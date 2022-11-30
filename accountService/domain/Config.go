package domain

type config struct {
	DatabaseUrl string `env:"POSTGRES_DB_URL"`
	JwtKey      string `env:"JWT_KEY" `
}

var Config *config

func InitConfig() (string, error) {
	if Config == nil {

		Config = &config{}
		//if err := env.Parse(Config); err != nil {
		//	log.Fatalf("something went wrong with environment, %e", err)
		//	return "something went wrong", err
		//}

		Config.JwtKey = "874967EC3EA3490F8F2EF6478B72A756"
		Config.DatabaseUrl = "postgres://postgres:postgres@host.docker.internal:5433/accountService?sslmode=disable"
		return "", nil
	}
	return "", nil
}
