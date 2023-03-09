package configs

type Config struct{}

func Setup() {
	config := Config{}
	config.GormDatabase()
	config.RedisConfig()
}
