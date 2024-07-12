package telegram

type (
	Config struct {
		WebApp WebAppConfig `yaml:"webapp"`
	}
	WebAppConfig struct {
		Auth AuthConfig
	}
	AuthConfig struct {
		BotTokens map[string]string `yaml:"bot_tokens" validate:"required,dive,keys,required,endkeys,required"` //nolint:lll
		Dummy     bool
	}
)
