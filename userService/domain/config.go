package domain

type config struct {
	DatabaseUrl   string `env:"POSTGRES_DB_URL"`
	NatsUrl       string `env:"NATS_URL"`
	JwtKey        string `env:"JWT_KEY" `
	GmailPassword string `env:"GMAIL_PASSWORD"`
	GmailAddress  string `env:"GMAIL_ADDRESS"`
	Url           string `env:"URL"`
}

var Config *config

func InitConfig() error {
	if Config == nil {
		Config = &config{}
		//if err := env.Parse(Config); err != nil {
		//	log.Fatalf("something went wrong with environment, %e", err)
		//	return err
		//}

		Config.JwtKey = "874967EC3EA3490F8F2EF6478B72A756"
		Config.DatabaseUrl = "postgres://postgres:postgres@host.docker.internal:5432/userService?sslmode=disable"
		Config.NatsUrl = "nats://host.docker.internal:4222"
		Config.GmailPassword = "xpqaovslkvtfpefb"
		Config.GmailAddress = "InternetBankingOV@gmail.com"
		Config.Url = "http://localhost:12345"
		return nil
	}
	return nil
}
