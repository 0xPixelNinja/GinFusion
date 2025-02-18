package config

import (
    "log"
    "strings"

    "github.com/spf13/viper"
)

// Config holds all the configuration for the application.
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
    JWT struct {
        Secret     string `mapstructure:"secret"`
        Expiration int    `mapstructure:"expiration"`
    } `mapstructure:"jwt"`
}

// LoadConfig reads configuration from config.yaml and environment variables.
func LoadConfig() *Config {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Warning: no config file found: %v", err)
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        log.Fatalf("Error decoding config: %v", err)
    }
    return &cfg
}
