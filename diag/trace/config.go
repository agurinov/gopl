package trace

type Config struct {
	AppName string  `yaml:"app_name"`
	Ratio   float64 `yaml:"ratio" validate:"min=0,max=1"`
}
