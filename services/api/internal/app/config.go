package app

import (
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
	ApiBotKey        string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
}{
	ApiPort:          "API_PORT",
	ApiJwtSecret:     "API_JWT_SECRET",
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
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string

	errors []string
}

type validatorFn[T comparable] func(T, string) (T, error)

func (c *AppConfig) appendError(e error) {
	c.errors = append(c.errors, e.Error())
}

func (c *AppConfig) validateInt(envVarName string, validators []validatorFn[int]) int {
	value, err := strconv.Atoi(os.Getenv(envVarName))

	if err != nil {
		c.appendError(fmt.Errorf("- %s must be an integer", envVarName))
		return -1
	}

	for _, validator := range validators {
		value, err = validator(value, envVarName)

		if err != nil {
			c.appendError(err)
			break
		}
	}

	return value
}

func (c *AppConfig) validateString(envVarName string, validators []validatorFn[string]) string {
	value := os.Getenv(envVarName)
	var err error

	for _, validator := range validators {
		value, err = validator(value, envVarName)

		if err != nil {
			c.appendError(err)
			break
		}
	}

	return value
}

func validateMinLength(minLength int) validatorFn[string] {
	return func(value string, envVarName string) (string, error) {

		if len(value) < minLength {
			return "", fmt.Errorf("- %s must be at least %d characters", envVarName, minLength)
		}

		return value, nil
	}
}

func validateRequiredPort(value int, envVarName string) (int, error) {
	minPort, maxPort := 1024, 49151

	if value < minPort || value > maxPort {
		return -1, fmt.Errorf("- %s must be a valid port (%d - %d)", envVarName, minPort, maxPort)
	}

	return value, nil
}

func newAppConfig(dbOnly bool) (*AppConfig, error) {
	log.Println("Loading application config")

	c := AppConfig{}

	if !dbOnly {
		c.Port = c.validateInt(envVars.ApiPort, []validatorFn[int]{validateRequiredPort})
		c.jwtSecret = []byte(c.validateString(envVars.ApiJwtSecret, []validatorFn[string]{validateMinLength(minSecretLength)}))
		c.ApiBotKey = c.validateString(envVars.ApiBotKey, []validatorFn[string]{validateMinLength(minSecretLength)})
	}

	c.PostgresPort = c.validateInt(envVars.PostgresPort, []validatorFn[int]{validateRequiredPort})
	c.PostgresDb = c.validateString(envVars.PostgresDb, []validatorFn[string]{validateMinLength(1)})
	c.PostgresHost = c.validateString(envVars.PostgresHost, []validatorFn[string]{validateMinLength(1)})
	c.PostgresUser = c.validateString(envVars.PostgresUser, []validatorFn[string]{validateMinLength(1)})
	c.PostgresPassword = c.validateString(envVars.PostgresPassword, []validatorFn[string]{validateMinLength(1)})

	if nil == c.errors {
		return &c, nil
	}

	return nil, fmt.Errorf("\n%s", strings.Join(c.errors[:], "\n"))
}

func (app *App) loadAppConfig() {
	config, err := newAppConfig(app.dbOnly)

	if nil != err {
		log.Fatal(err)
	}

	app.config = config
}
