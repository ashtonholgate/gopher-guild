package config

import (
	"gopherService/customErrors"
	"os"
	"strconv"
)

type Config struct {
	DB_HOST                     string
	DB_PORT                     int
	DB_USER                     string
	DB_PASSWORD                 string
	DB_NAME                     string
	MAXIMUM_GOPHER_NAME_LENGTH  int
	MAXIMUM_GOPHER_COLOR_LENGTH int
}

func New() (Config, error) {
	var config = Config{}
	if value, exists := os.LookupEnv("DB_HOST"); exists {
		config.DB_HOST = value
	} else {
		return Config{}, &customErrors.EnvVarNotFoundError{EnvVar: "DB_HOST"}
	}

	if value, exists := os.LookupEnv("DB_PORT"); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, &customErrors.EnvVarWrongTypeError{EnvVar: "DB_PORT", Value: value, Type: "int"}
		} else {
			config.DB_PORT = intValue
		}
	} else {
		return Config{}, &customErrors.EnvVarNotFoundError{EnvVar: "DB_PORT"}
	}

	if value, exists := os.LookupEnv("DB_USER"); exists {
		config.DB_USER = value
	} else {
		return Config{}, &customErrors.EnvVarNotFoundError{EnvVar: "DB_USER"}
	}

	if value, exists := os.LookupEnv("DB_PASSWORD"); exists {
		config.DB_PASSWORD = value
	} else {
		return Config{}, &customErrors.EnvVarNotFoundError{EnvVar: "DB_PASSWORD"}
	}

	if value, exists := os.LookupEnv("DB_NAME"); exists {
		config.DB_NAME = value
	} else {
		return Config{}, &customErrors.EnvVarNotFoundError{EnvVar: "DB_NAME"}
	}

	if value, exists := os.LookupEnv("MAXIMUM_GOPHER_NAME_LENGTH"); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, &customErrors.EnvVarWrongTypeError{EnvVar: "MAXIMUM_GOPHER_NAME_LENGTH", Value: value, Type: "int"}
		} else {
			config.MAXIMUM_GOPHER_NAME_LENGTH = intValue
		}
	} else {
		config.MAXIMUM_GOPHER_NAME_LENGTH = 50
	}

	if value, exists := os.LookupEnv("MAXIMUM_GOPHER_COLOR_LENGTH"); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, &customErrors.EnvVarWrongTypeError{EnvVar: "MAXIMUM_GOPHER_COLOR_LENGTH", Value: value, Type: "int"}
		} else {
			config.MAXIMUM_GOPHER_COLOR_LENGTH = intValue
		}
	} else {
		config.MAXIMUM_GOPHER_COLOR_LENGTH = 50
	}

	return config, nil
}
