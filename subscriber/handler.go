package subscriber

import (
	"backEnd/config"
	"backEnd/dto"
	"backEnd/helper"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

func sub(client *mqtt.Client, topicName string) error {
	token := (*client).Subscribe(topicName, 1, messageSubHandler)
	token.Wait()
	log.Info().Msg("subscribed to topic:" + topicName)
	return nil
}

func InitMqttSubscribers() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.MqttBrokerUrl, config.MqttBrokerPort))
	opts.SetClientID(config.MqttClientId)
	opts.SetUsername(config.MqttUserName)
	opts.SetPassword(config.MqttPassword)
	opts.SetDefaultPublishHandler(messageSubHandler)
	opts.SetOrderMatters(false)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	for key, _ := range helper.TopicMap {
		err := sub(&client, key)
		if err != nil {
			log.Error().Err(err).
				Str("topicName", key).
				Msg("topic subscription error")
		}
	}

	for _, key := range helper.RegExpToWildCard {
		err := sub(&client, key)
		if err != nil {
			log.Error().Err(err).
				Str("topicName", key).
				Msg("topic subscription error")
		}
	}
	return nil
}

var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if log.Trace().Enabled() {
		log.Trace().
			Str("topicName", msg.Topic()).
			Interface("payload", msg.Payload()).
			Msg("received message from mqtt")
	} else {
		log.Debug().Msg("received message from mqtt")
	}

	kafkaTopicName := helper.TopicMap[msg.Topic()]
	if kafkaTopicName == "" {
		var err error
		kafkaTopicName, err = helper.FindKafkaTopicFromMqttTopic(msg.Topic())
		if err != nil {
			log.Error().Err(err).
				Str("topicName", kafkaTopicName).
				Msg("map failed")
			return
		}
	}
	helper.DefaultChanel <- dto.DefaultChannelStruct{
		TopicName: kafkaTopicName,
		Payload:   msg.Payload(),
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info().Msg("Connected with MQTT broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Error().Err(err).Msg("connection lost with MQTT broker")
}
