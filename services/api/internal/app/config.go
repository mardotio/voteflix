package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const minSecretLength = 32

var envVars = struct {
	ApiPort          string
	ApiJwtSecret     string
	ApiInDocker      string
	ApiBotKey        string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
}{
	ApiPort:          "API_PORT",
	ApiJwtSecret:     "API_JWT_SECRET",
	ApiInDocker:      "API_IN_DOCKER",
	ApiBotKey:        "API_BOT_KEY",
	PostgresHost:     "POSTGRES_HOST",
	PostgresPort:     "POSTGRES_PORT",
	PostgresUser:     "POSTGRES_USER",
	PostgresPassword: "POSTGRES_PASSWORD",
	PostgresDb:       "POSTGRES_DB",
}

type AppConfig struct {
	Port             int
	jwtSecret        []byte
	ApiBotKey        string
	ApiInDocker      bool
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
}

func validateRequiredString(envVarName string, errSlice *[]string) (string, error) {
	value := os.Getenv(envVarName)

	if value == "" {
		err := fmt.Sprintf("- %s cannot be empty", envVarName)
		*errSlice = append(*errSlice, err)
		return "", errors.New(err)
	}

	return value, nil
}

func validateRequiredPort(envVarName string, errSlice *[]string) (int, error) {
	value, err := strconv.Atoi(os.Getenv(envVarName))

	if err != nil {
		err := fmt.Sprintf("- %s must be a valid port", envVarName)
		*errSlice = append(*errSlice, err)
		return -1, errors.New(err)
	}

	return value, nil
}

func (config *AppConfig) setPostgresPort(errSlice *[]string) {
	if port, err := validateRequiredPort(envVars.PostgresPort, errSlice); err == nil {
		config.PostgresPort = port
	}
}

func (config *AppConfig) setPostgresHost(errSlice *[]string) {
	if host, err := validateRequiredString(envVars.PostgresHost, errSlice); err == nil {
		config.PostgresHost = host
	}
}

func (config *AppConfig) setPostgresPassword(errSlice *[]string) {
	if password, err := validateRequiredString(envVars.PostgresPassword, errSlice); err == nil {
		config.PostgresPassword = password
	}
}

func (config *AppConfig) setPostgresDb(errSlice *[]string) {
	if db, err := validateRequiredString(envVars.PostgresDb, errSlice); err == nil {
		config.PostgresDb = db
	}
}

func (config *AppConfig) setPostgresUser(errSlice *[]string) {
	if user, err := validateRequiredString(envVars.PostgresUser, errSlice); err == nil {
		config.PostgresUser = user
	}
}

func (config *AppConfig) setApiPort(errSlice *[]string) {
	if port, err := validateRequiredPort(envVars.ApiPort, errSlice); err == nil {
		config.Port = port
	}
}

func (config *AppConfig) setApiJwtSecret(errSlice *[]string) {
	secret, secretErr := validateRequiredString(envVars.ApiJwtSecret, errSlice)

	if secretErr != nil {
		return
	}

	if minSecretLength > len(secret) {
		*errSlice = append(*errSlice, fmt.Sprintf("- %s must be at least %d characters", envVars.ApiJwtSecret, minSecretLength))
		return
	}

	config.jwtSecret = []byte(secret)
}

func (config *AppConfig) setApiBotKey(errSlice *[]string) {
	secret, secretErr := validateRequiredString(envVars.ApiBotKey, errSlice)

	if secretErr != nil {
		return
	}

	if minSecretLength > len(secret) {
		*errSlice = append(*errSlice, fmt.Sprintf("- %s must be at least %d characters", envVars.ApiBotKey, minSecretLength))
		return
	}

	config.ApiBotKey = secret
}

func parseEnv() (*AppConfig, error) {
	log.Println("Loading application config")
	var parseErr []string

	config := AppConfig{}

	config.setApiPort(&parseErr)
	config.setApiJwtSecret(&parseErr)
	config.setApiBotKey(&parseErr)
	config.setPostgresDb(&parseErr)
	config.setPostgresUser(&parseErr)
	config.setPostgresHost(&parseErr)
	config.setPostgresPassword(&parseErr)
	config.setPostgresPort(&parseErr)

	if nil == parseErr {
		return &config, nil
	}

	return nil, fmt.Errorf("\n%s", strings.Join(parseErr[:], "\n"))
}

func (app *App) loadAppConfig() {
	config, err := parseEnv()

	if nil != err {
		log.Fatal(err)
	}

	app.config = config
}
