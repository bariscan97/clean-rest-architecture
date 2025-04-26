package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Database struct {
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        User     string `mapstructure:"user"`
        Password string `mapstructure:"password"`
        DBName   string `mapstructure:"dbname"`
        SSLMode  string `mapstructure:"sslmode"`
    } `mapstructure:"database"`
    App struct {
        Port int  
    } `mapstructure:"app"`
}

func LoadConfig(path string) (*Config, error) {
    var c Config

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(path)

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    if err := viper.Unmarshal(&c); err != nil {
        return nil, err
    }

    return &c, nil
}
