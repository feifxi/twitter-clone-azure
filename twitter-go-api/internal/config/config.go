package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Environment              string `mapstructure:"ENVIRONMENT"`
	DBSource                 string `mapstructure:"DATABASE_URL"`
	HTTPServerAddress        string `mapstructure:"HTTP_SERVER_ADDRESS"`
	FrontendURL              string `mapstructure:"FRONTEND_URL"`
	CookieDomain             string `mapstructure:"COOKIE_DOMAIN"`
	CookieSameSite           string `mapstructure:"COOKIE_SAME_SITE"`
	CookieSecure             bool   `mapstructure:"COOKIE_SECURE"`
	TokenSymmetricKey        string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDurationMinutes     int    `mapstructure:"TOKEN_DURATION_MINUTES"`
	RefreshTokenDurationDays int    `mapstructure:"REFRESH_TOKEN_DURATION_DAYS"`
	GoogleClientID           string `mapstructure:"GOOGLE_CLIENT_ID"`
	AzureStorageConnString   string `mapstructure:"AZURE_STORAGE_CONNECTION_STRING"`
	AzureStorageContainer    string `mapstructure:"AZURE_STORAGE_CONTAINER_NAME"`
	RedisAddress             string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword            string `mapstructure:"REDIS_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetDefault("FRONTEND_URL", "http://localhost:3000")
	viper.SetDefault("COOKIE_DOMAIN", "")
	viper.SetDefault("COOKIE_SAME_SITE", "Lax")
	viper.SetDefault("COOKIE_SECURE", false)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// It's ok if app.env doesn't exist, we will use auto env
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
