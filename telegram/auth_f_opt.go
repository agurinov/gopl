package telegram

import "go.uber.org/zap"

func WithAuthDummy(dummyEnabled bool) AuthOption {
	return func(s *Auth) error {
		s.dummyEnabled = dummyEnabled

		return nil
	}
}

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

func WithAuthNoBot() AuthOption {
	return func(s *Auth) error {
		s.noBot = true

		return nil
	}
}
