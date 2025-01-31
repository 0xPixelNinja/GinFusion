package config

import (
    "log"

    "github.com/spf13/viper"
)

// Config holds the configuration values for the application.
type Config struct {
    Server struct {
        Address string `mapstructure:"address"`
    } `mapstructure:"server"`
    Admin struct {
        Address string `mapstructure:"address"`
    } `mapstructure:"admin"`
    Redis struct {
        Addr     string `mapstructure:"addr"`
        Password string `mapstructure:"password"`
        DB       int    `mapstructure:"db"`
    } `mapstructure:"redis"`
}

// LoadConfig reads configuration from config.yaml and environment variables.
func LoadConfig() *Config {

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    // Read the config file
    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Warning: Error reading config file, %s", err)
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        log.Fatalf("Error decoding configuration into struct: %v", err)
    }

    return &cfg
}
