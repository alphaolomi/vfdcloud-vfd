package env

const (
	DevEnv     Env = "development"
	TestEnv    Env = "test"
	StagingEnv Env = "staging"
	ProdEnv    Env = "production"
)

type Env string

func (e Env) String() string {
	return string(e)
}

func ParseEnv(s string) Env {
	switch s {
	case "development", "dev":
		return DevEnv
	case "test", "testing":
		return TestEnv
	case "staging":
		return StagingEnv
	case "production", "prod":
		return ProdEnv
	default:
		return DevEnv
	}
}
