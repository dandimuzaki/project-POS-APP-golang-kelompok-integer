package utils

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Configuration struct {
	AppName     string
	Port        string
	Debug       bool
	Limit       int
	PathLogging string
	DB          DatabaseConfig
	SMTP SMTPConfig
	BaseURL string
	BusinessRules BusinessRules
}

type DatabaseConfig struct {
	Name     string
	Username string
	Password string
	Host     string
	Port     string
	MaxConn  int32
}

type SMTPConfig struct {
	Host string
	Port int
	Email string
	Password string
}

type BusinessRules struct {
	TaxRate int
	ProfitMargin int
	DefaultShiftStart string
	DefaultShiftEnd string
}

func ReadConfiguration() (Configuration, error) {
	// get config from env file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}

	// get config from os variable
	viper.AutomaticEnv()

	// get config from flag
	pflag.Int("port-app", 0, "port for app golang")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	return Configuration{
		AppName:     viper.GetString("APP_NAME"),
		Port:        viper.GetString("PORT"),
		Debug:       viper.GetBool("DEBUG"),
		Limit:       viper.GetInt("LIMIT"),
		PathLogging: viper.GetString("PATH_LOGGING"),
		DB: DatabaseConfig{
			Name:     viper.GetString("DB_NAME"),
			Username: viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			MaxConn:  viper.GetInt32("DB_MAX_CONN"),
		},
		SMTP: SMTPConfig{
			Host: viper.GetString("SMTP_HOST"),
			Port: viper.GetInt("SMTP_PORT"),
			Email: viper.GetString("SMTP_USER"),
			Password: viper.GetString("SMTP_PASSWORD"),
		},
		BaseURL: viper.GetString("APP_URL"),
		BusinessRules: BusinessRules{
			TaxRate: viper.GetInt("TAX_RATE"),
			ProfitMargin: viper.GetInt("PROFIT_MARGIN"),
			DefaultShiftStart: viper.GetString("DEFAULT_SHIFT_START"),
			DefaultShiftEnd: viper.GetString("DEFAULT_SHIFT_END"),
		},
	}, nil

}
