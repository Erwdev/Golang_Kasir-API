package config

import (
    "os"
    "strings"

    "github.com/spf13/viper"
)

// Config hanya pegang PORT dan DB_CONN.
// Railway set env di OS level, bukan di file .env,
// jadi kita baca os.Getenv() langsung sebagai prioritas utama.

type Config struct {
    Port   string
    DBConn string
}

func Load() Config {
    if _, err := os.Stat(".env"); err == nil {
        viper.SetConfigFile(".env")
        viper.SetConfigType("env")
        viper.ReadInConfig()
    }

    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    cfg := Config{}

    if port := os.Getenv("PORT"); port != "" {
        cfg.Port = port
    } else if port := viper.GetString("PORT"); port != "" {
        cfg.Port = port
    } else {
        cfg.Port = "8080"
    }

    if dbConn := os.Getenv("DB_CONN"); dbConn != "" {
        cfg.DBConn = dbConn
    } else if dbConn := os.Getenv("DATABASE_URL"); dbConn != "" {
        cfg.DBConn = dbConn
    } else {
        cfg.DBConn = viper.GetString("DB_CONN")
    }

    return cfg
}