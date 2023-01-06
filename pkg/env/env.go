package env

const (
	DEV     Env = "development"
	TEST    Env = "test"
	STAGING Env = "staging"
	PROD    Env = "production"
)

type Env string

func (e Env) String() string {
	return string(e)
}

func Parse(s string) Env {
	switch s {
	case "development", "dev":
		return DEV
	case "test", "testing":
		return TEST
	case "staging":
		return STAGING
	case "production", "prod":
		return PROD
	default:
		return DEV
	}
}
