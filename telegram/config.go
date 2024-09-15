package telegram

type (
	Config struct {
		WebApp WebAppConfig `yaml:"webapp"`
		Bot    BotConfig    `yaml:"bot"`
	}
	BotConfig struct {
		Token string `validate:"required"`
	}
	WebAppConfig struct {
		Auth WebAppAuthConfig
	}
	WebAppAuthConfig struct {
		BotTokens        map[string]string `yaml:"bot_tokens" validate:"required,dive,keys,required,endkeys,required"` //nolint:lll
		NoSignatureCheck bool              `yaml:"no_signature_check"`
	}
)
