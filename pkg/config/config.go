package config

type (
	App struct {
		Name       string `mapstructure:"name"`
		Version    string `mapstructure:"version"`
		Production bool   `mapstructure:"production"`
	}

	HTTP struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	GRPC struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Password string `mapstructure:"password"`
	}
)
