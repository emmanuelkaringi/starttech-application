package config

import (
        "strings"

        "github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
        ServerPort         string   `mapstructure:"PORT"`
        MongoURI           string   `mapstructure:"MONGO_URI"`
        DBName             string   `mapstructure:"DB_NAME"`
        JWTSecretKey       string   `mapstructure:"JWT_SECRET_KEY"`
        JWTExpirationHours int      `mapstructure:"JWT_EXPIRATION_HOURS"`
        EnableCache        bool     `mapstructure:"ENABLE_CACHE"`
        RedisAddr          string   `mapstructure:"REDIS_ADDR"`
        RedisPassword      string   `mapstructure:"REDIS_PASSWORD"`
        LogLevel           string   `mapstructure:"LOG_LEVEL"`
        LogFormat          string   `mapstructure:"LOG_FORMAT"`
        CookieDomains      []string `mapstructure:"COOKIE_DOMAINS"`
        SecureCookie       bool     `mapstructure:"SECURE_COOKIE"`
        AllowedOrigins     []string `mapstructure:"ALLOWED_ORIGINS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
        viper.AddConfigPath(path)
        viper.SetConfigName(".env")
        viper.SetConfigType("env")

        viper.AutomaticEnv()

        // Set default values
        viper.SetDefault("PORT", "8080")
        viper.SetDefault("ENABLE_CACHE", false)
        viper.SetDefault("JWT_EXPIRATION_HOURS", 72)
        viper.SetDefault("COOKIE_DOMAINS", "")
        viper.SetDefault("SECURE_COOKIE", false)
        viper.SetDefault("ALLOWED_ORIGINS", "http://localhost:5173")

        err = viper.ReadInConfig()
        if err != nil {
                if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
                        return
                }
        }

        err = viper.Unmarshal(&config)
        if err != nil {
                return
        }

        // Manually handle comma-separated strings for slices
        if allowedOrigins := viper.GetString("ALLOWED_ORIGINS"); allowedOrigins != "" {
                parts := strings.Split(allowedOrigins, ",")
                var cleaned []string
                for _, p := range parts {
                        trimmed := strings.TrimSpace(p)
                        trimmed = strings.Trim(trimmed, "\"'")
                        if trimmed != "" {
                                cleaned = append(cleaned, trimmed)
                        }
                }
                if len(cleaned) > 0 {
                        config.AllowedOrigins = cleaned
                }
        }

        // Handle cookie domains - only override if explicitly set
        cookieDomains := viper.GetString("COOKIE_DOMAINS")
        if cookieDomains != "" {
                parts := strings.Split(cookieDomains, ",")
                var cleaned []string
                for _, p := range parts {
                        trimmed := strings.TrimSpace(p)
                        trimmed = strings.Trim(trimmed, "\"'")
                        if trimmed != "" {
                                cleaned = append(cleaned, trimmed)
                        }
                }
                if len(cleaned) > 0 {
                        config.CookieDomains = cleaned
                }
        } else {
                // Empty string means no domain restriction
                config.CookieDomains = []string{}
        }

        return
}
