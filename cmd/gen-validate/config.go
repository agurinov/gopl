package main

import (
	"bytes"
	"go/types"
	"io"
	"os"
	"strings"

	"github.com/agurinov/gopl/ast"
	"github.com/agurinov/gopl/env/envvars"
	"github.com/dave/jennifer/jen"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	output         io.ReadWriter
	sourceFile     string
	sourcePkg      string
	pkg            string
	outputFilepath string
	types          []string
	sourceLine     int
}

func ParseConfig() (cfg Config, err error) {
	cfg.sourceFile, err = envvars.GoFile.Value()
	if err != nil {
		return cfg, err
	}

	cfg.sourceLine, err = envvars.GoLine.Value()
	if err != nil {
		return cfg, err
	}

	cfg.sourcePkg, err = envvars.GoPackage.Value()
	if err != nil {
		return cfg, err
	}

	cfg.types = strings.Split(*typesPtr, ",")
	cfg.pkg = *packagePtr

	switch o := *outputPtr; o {
	case "-", "":
		cfg.output = os.Stdout
	default:
		cfg.outputFilepath = o
		cfg.output = new(bytes.Buffer)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	s := struct {
		Output     io.ReadWriter `validate:"required"`
		SourceFile string        `validate:"required"`
		SourcePkg  string        `validate:"required"`
		Pkg        string        `validate:"required"`
		Types      []string      `validate:"required"`
		SourceLine int           `validate:"required"`
	}{
		SourceFile: c.sourceFile,
		SourceLine: c.sourceLine,
		SourcePkg:  c.sourcePkg,
		Types:      c.types,
		Pkg:        c.pkg,
		Output:     c.output,
	}

	if err := validator.New().Struct(s); err != nil {
		return err
	}

	return nil
}

func (c Config) Generate() error {
	targets, err := c.analyze()
	if err != nil {
		return err
	}

	f := c.generate(targets...)

	return c.render(f)
}

func (c Config) analyze() ([]Target, error) {
	asts, err := c.getStructTypes()
	if err != nil {
		return nil, err
	}

	targets := make([]Target, 0, len(asts))

	for structName, structData := range asts {
		target := Target{
			name:     structName,
			generics: make([]string, 0),
			fields:   make([]TargetField, 0, len(structData.structType.Fields.List)),
			imports:  ast.ParseImports(structData.imports),
		}

		// Possible generics
		if structData.typeSpec.TypeParams != nil {
			for _, genericType := range structData.typeSpec.TypeParams.List {
				target.generics = append(target.generics,
					genericType.Names[0].Name,
				)
			}
		}

		for _, field := range structData.structType.Fields.List {
			target.fields = append(target.fields, TargetField{
				name: field.Names[0].Name,
				typ:  types.ExprString(field.Type),
				tags: ast.ParseStructTags(field.Tag),
			})
		}

		targets = append(targets, target)
	}

	return targets, nil
}

func (c Config) generate(targets ...Target) *jen.File {
	const (
		codegenComment        = ""
		validator             = "github.com/go-playground/validator/v10"
		nonStandartValidators = validator + "/non-standard/validators"
	)

	f := jen.NewFile(c.pkg)
	f.HeaderComment(codegenComment)

	f.ImportName(validator, "validator")
	f.ImportName(nonStandartValidators, "validators")

	var (
		validateFuncBody = []jen.Code{
			jen.Id("v").Op(":=").Qual(validator, "New").Call(),
			jen.If(
				jen.Err().Op(":=").Id("v").Dot("RegisterValidation").Call(
					jen.Lit("notblank"),
					jen.Qual(nonStandartValidators, "NotBlank"),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return().Err(),
			),
			jen.Line(),

			jen.If(
				jen.Err().Op(":=").Id("v").Dot("Struct").Call(jen.Id("s")),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return().Err(),
			),
			jen.Line(),

			jen.Return().Nil(),
		}
		validateEmptyFuncBody = []jen.Code{
			jen.Return().Nil(),
		}
	)

	for _, target := range targets {
		var (
			structDefinition     = make([]jen.Code, 0, len(target.fields))
			structInitialization = make([]jen.Code, 0, len(target.fields))
		)

		for _, field := range target.fields {
			if field.tags == nil {
				continue
			}

			validateStructTag := field.tags["validate"]
			if validateStructTag == "" {
				continue
			}

			fieldTypeDefinitions := make([]string, 2)

			switch splittedType := strings.Split(field.typ, "."); len(splittedType) {
			case 2:
				importAlias := splittedType[0]
				importPath, exists := target.imports[importAlias]

				if exists {
					fieldTypeDefinitions[0] = importPath
				} else {
					fieldTypeDefinitions[0] = splittedType[0]
				}

				fieldTypeDefinitions[1] = splittedType[1]
			default:
				fieldTypeDefinitions[0] = ""
				fieldTypeDefinitions[1] = field.typ
			}

			structDefinition = append(structDefinition,
				jen.Id(
					strings.Title(field.name),
				).Qual(
					fieldTypeDefinitions[0], fieldTypeDefinitions[1],
				).Tag(
					map[string]string{"validate": validateStructTag},
				),
			)
			structInitialization = append(structInitialization,
				jen.Id(
					strings.Title(field.name),
				).Op(":").Id("obj").Dot(field.name).Id(","),
			)
		}

		var body []jen.Code

		switch target.InvolvedInGenerate() {
		case true:
			body = []jen.Code{
				jen.Id("s").Op(":=").Struct(structDefinition...).Block(structInitialization...),
				jen.Line(),
			}
			body = append(body, validateFuncBody...)
		default:
			body = validateEmptyFuncBody
		}

		for importAlias, importPath := range target.imports {
			f.ImportAlias(importPath, importAlias)
		}

		f.Func().
			Params(
				jen.Id("obj").Id(target.name).TypesFunc(func(g *jen.Group) {
					for _, genericType := range target.generics {
						g.Add(jen.Id(genericType))
					}
				}),
			).
			Id("Validate").
			Params().Error().
			Block(body...).Line()
	}

	return f
}

func (c Config) render(f *jen.File) error {
	if err := f.Render(c.output); err != nil {
		return err
	}

	if c.outputFilepath == "" {
		return nil
	}

	file, err := os.Create(c.outputFilepath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(file, c.output); err != nil {
		return err
	}

	return nil
}
