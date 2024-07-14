package telegram

import "go.uber.org/zap"

func WithAuthBotTokens(botTokens map[string]string) AuthOption {
	return func(s *Auth) error {
		s.botTokens = botTokens

		return nil
	}
}

func WithAuthLogger(logger *zap.Logger) AuthOption {
	return func(s *Auth) error {
		if logger == nil {
			return nil
		}

		s.logger = logger.Named("telegram.auth")

		return nil
	}
}

func WithAuthNoBotAllowed() AuthOption {
	return func(s *Auth) error {
		s.noBotAllowed = true

		return nil
	}
}

func WithAuthNoSignatureCheck(noSignatureCheck bool) AuthOption {
	return func(s *Auth) error {
		s.noSignatureCheck = noSignatureCheck

		return nil
	}
}
