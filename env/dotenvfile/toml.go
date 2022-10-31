package dotenvfile

type DotenvFileOption byte

const (
	DOTENV_NO_OVERRIDE DotenvFileOption = 1 << iota
)
