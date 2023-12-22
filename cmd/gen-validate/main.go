package main

import (
	"flag"
	"log"
)

var (
	typesPtr   = flag.String("types", "", "Comma separated list of types")
	packagePtr = flag.String("package", "", "Name of package for codegen")
	outputPtr  = flag.String("output", "-", "Output for codegen")
)

func main() {
	if flag.Parse(); !flag.Parsed() {
		log.Fatal("can't parse cli flags")
	}

	cfg, err := ParseConfig()
	if err != nil {
		log.Fatal("can't parse config: ", err)
	}

	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	if err := cfg.Generate(); err != nil {
		panic(err)
	}
}
