package config

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
)

// General EnVars
var RunMode string
var TopicRoute []string
var WildCardTopicRoute []string

// MQTT envars
var MqttBrokerUrl string
var MqttBrokerPort int
var MqttClientId string
var MqttUserName string
var MqttPassword string

// KAFKA envars
var KafkaBroker string
var SaslSsl bool
var SaslOauthbearerClientId string
var SaslOauthbearerClientSecret string
var SaslOauthbearerTokenEndpointUri string

// sqlite envars
var DatabasePath string
var InvokeLatencySec int
var InvokeMaxGoRoutines int
var InvokeMinGoRoutines int

// Logger Envars
var LoggingLvl logLvl

// temp variables
var boolVal bool
var err error

func InitEnvironmentVariables() error {
	log.Info().Msg("initializing environment variables")
	//checking runMode
	RunMode = os.Getenv("RUN_MODE")
	if RunMode == "" {
		log.Warn().Msg("RUN_MODE empty, switching to develop mode")
		RunMode = DEVELOP
	}
	var err error

	log.Info().
		Str("RUN_MODE", RunMode).
		Msg("")

	//loading envArs from .env file is runMode != PRODUCTION
	if RunMode != PRODUCTION {
		err = godotenv.Load()
		if err != nil {
			return err
		}
	}

	err = initMqttEnvironmentVariables()
	if err != nil {
		return err
	}

	err = initKafkaEnvironmentVariables()
	if err != nil {
		return err
	}
	err = initDatabaseEnvironmentVariables()
	if err != nil {
		return err
	}

	TopicRoute = strings.Split(os.Getenv("TOPIC_LIST"), ",")
	WildCardTopicRoute = strings.Split(os.Getenv("WILDCARD_TOPIC_LIST"), ",")
	return nil
}

func initMqttEnvironmentVariables() error {

	//BrokerURL
	MqttBrokerUrl, boolVal = os.LookupEnv("MQTT_BROKER_URL")
	if boolVal == false {
		return errors.New("MQTT_BROKER_URL not found in envVars")
	}

	//BrokerPort
	BrokerPortString, boolVal := os.LookupEnv("MQTT_BROKER_PORT")
	if boolVal == false {
		return errors.New("MQTT_BROKER_PORT not found in envVars")
	}
	MqttBrokerPort, err = strconv.Atoi(BrokerPortString)
	if err != nil {
		return errors.New("unable to convert MQTT_BROKER_PORT to valid int")
	}

	//clientId
	MqttClientId, boolVal = os.LookupEnv("MQTT_CLIENT_ID")
	if boolVal == false {
		return errors.New("MQTT_CLIENT_ID not found in envVars")
	}

	//username
	MqttUserName, boolVal = os.LookupEnv("MQTT_USER_NAME")
	if boolVal == false {
		return errors.New("MQTT_USER_NAME not found in envVars")
	}
	//user password
	MqttPassword, boolVal = os.LookupEnv("MQTT_USER_PASSWORD")
	if boolVal == false {
		return errors.New("MQTT_USER_PASSWORD not found in envVars")
	}

	initLoggingEnvironmentVariables()

	return nil
}

func initKafkaEnvironmentVariables() error {
	KafkaBroker, boolVal = os.LookupEnv("KAFKA_BROKER")
	if boolVal == false {
		return errors.New("KAFKA_BROKER not found in envVars")
	}

	tempSaslSsl, boolVal := os.LookupEnv("SASL_SSL")
	if boolVal == false {
		SaslSsl = false
	} else {
		if tempSaslSsl != TRUE && tempSaslSsl != FALSE {
			return errors.New("SASL_SSL should be true or false")
		}
		if tempSaslSsl == TRUE {
			SaslOauthbearerClientId, boolVal = os.LookupEnv("SASL_OAUTHBEARER_CLIENT_ID")
			if boolVal == false {
				return errors.New("SASL_OAUTHBEARER_CLIENT_ID not found in envVars")
			}

			SaslOauthbearerClientSecret, boolVal = os.LookupEnv("SASL_OAUTHBEARER_CLIENT_SECRET")
			if boolVal == false {
				return errors.New("SASL_OAUTHBEARER_CLIENT_SECRET not found in envVars")
			}

			SaslOauthbearerTokenEndpointUri, boolVal = os.LookupEnv("SASL_OAUTHBEARER_TOKEN_ENDPOINT_URI")
			if boolVal == false {
				return errors.New("SASL_OAUTHBEARER_TOKEN_ENDPOINT_URI not found in envVars")
			}
			SaslSsl = true
		} else {
			SaslSsl = false
		}
	}

	return nil
}

func initDatabaseEnvironmentVariables() error {
	DatabasePath, boolVal = os.LookupEnv("DB_PERSIST_PATH")
	if boolVal == false {
		log.Warn().Msg("DB_PERSIST_PATH not found, switching to in-memory mode")
	}

	invokeLatencySec, boolVal := os.LookupEnv("INVOKE_LATENCY_SEC")
	if boolVal == false {
		return errors.New("INVOKE_LATENCY_SEC not found in envVars")
	}
	InvokeLatencySec, err = strconv.Atoi(invokeLatencySec)
	if err != nil {
		return errors.New("INVOKE_LATENCY_SEC is not a valid integer")
	}

	invokeMaxGoRoutines, boolVal := os.LookupEnv("INVOKE_MAX_GOROUTINE")
	if boolVal == false {
		return errors.New("INVOKE_MAX_GOROUTINE not found in envVars")
	}
	InvokeMaxGoRoutines, err = strconv.Atoi(invokeMaxGoRoutines)
	if err != nil {
		return errors.New("INVOKE_MAX_GOROUTINE is not a valid integer")
	}

	invokeMinGoRoutines, boolVal := os.LookupEnv("INVOKE_MIN_GOROUTINE")
	if boolVal == false {
		return errors.New("INVOKE_MIN_GOROUTINE not found in envVars")
	}
	InvokeMinGoRoutines, err = strconv.Atoi(invokeMinGoRoutines)
	if err != nil {
		return errors.New("INVOKE_MIN_GOROUTINE is not a valid integer")
	}

	return nil
}
func initLoggingEnvironmentVariables() {
	tempLevel, boolVal := os.LookupEnv("LOG_LVL")
	if boolVal == false {
		LoggingLvl = WarnLvl

	} else {
		LoggingLvl = logLvl(tempLevel)
	}
}
