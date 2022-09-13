package main

import (
	"backEnd/config"
	"backEnd/database"
	"backEnd/helper"
	"backEnd/logging"
	"backEnd/publisher"
	"backEnd/subscriber"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	//PreInit our logging package
	logging.PreInit()
	//initialize environment variables
	err := config.InitEnvironmentVariables()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("environment variables init error")
	} else {
		log.Debug().Msg("environment variables loaded")
	}

	//initializing our logging package
	log.Debug().Msg("initializing logging")
	logging.Init()
	//init topicMap
	log.Debug().Msg("initializing TopicMap")
	helper.InitTopicMap(config.TopicRoute)
	log.Debug().Msg("initializing WildCardMap")
	helper.InitWildCard(config.WildCardTopicRoute)

	//init channels
	log.Info().Msg("initializing channels")
	helper.InitChannels()
	//init kafka connection
	log.Info().Msg("initializing kafka")
	err = publisher.GetKafkaManager().InitProducer()
	if err != nil {
		log.Error().Err(err).Msg("kafka error")
	}

	//database operations
	err = database.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("database init error")
	}
	//init kafka publishers
	go publisher.DatabaseChannelHandler()
	go publisher.DefaultChannelHandler()

	go database.DeleteDbChannelHandler()
	go database.InsertChannelHandler()
	//init mqtt client and test connection
	err = subscriber.InitMqttClient()
	if err != nil {
		log.Fatal().Err(err).Msg("MQTT client connection error")
	}
	//init mqtt subscribers
	err = subscriber.InitMqttSubscribers()
	if err != nil {
		log.Error().Err(err).Msg("MQTT subscriber connection error")
	}

	helper.FlashLogo()
	//invoke DbRestore function
	go database.DbInvoke()

	for {
		time.Sleep(time.Hour * 1000)
	}

}
