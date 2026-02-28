package main

import (
    "fmt"
    "os"
    "github.com/spf13/viper"
)

type Config struct {
    DBSource string `mapstructure:"DATABASE_URL"`
}

func main() {
    os.Setenv("DATABASE_URL", "my_db_url")

    // Simulate what config.go does
    viper.SetConfigName("nonexistent")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()
    
    _ = viper.ReadInConfig() // Will fail, ConfigFileNotFoundError

    var config Config
    viper.Unmarshal(&config)

    fmt.Printf("DBSource parsed: '%s'\n", config.DBSource)
}
