package publisher

import (
	"backEnd/config"
	"backEnd/helper"
	"backEnd/syncsys"
	"github.com/rs/zerolog/log"
)

func DefaultChannelHandler() {
	for {
		msg := <-helper.DefaultChanel
		go func() {
			err := GetKafkaManager().Publish(msg.TopicName, msg.Payload)
			if err != nil {
				if log.Trace().Enabled() {
					log.Trace().
						Bool(config.CONFIDENTIAL, true).
						Str("topicName:", msg.TopicName).
						Interface("payload", msg.Payload).
						Msg("publish successful from defaultChannel")
				} else {
					log.Error().Err(err).Msg("unable to publish from defaultChannel")
				}
				log.Warn().Msg("backing up message to database channel")
				helper.InsertDbChannel <- msg
			} else {
				log.Info().Msg("publish successful from defaultChannel")

			}
		}()

	}

}

func DatabaseChannelHandler() {
	for {
		msg := <-helper.DatabaseChannel
		syncsys.Add()
		go func() {
			defer syncsys.Delete()
			err := GetKafkaManager().Publish(msg.TopicName, msg.Payload)
			if err != nil {
				//log.Println("[ERROR] error publishing to kafka from dbChannel", err.Error())
			} else {
				//log.Println("[INFO] kafka message publish successful from dbChannel, topic: ", msg.TopicName)
				//log.Println("[INFO] deleting message from sqlite after successful publish")
				helper.DeleteDbChannel <- msg.OID
			}
		}()

	}
}
