package subscriber

import (
	"backEnd/config"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//for MQTT v5.0: "github.com/eclipse/paho.golang/paho"
//for MQTT v3.1: "github.com/eclipse/paho.subscriber.golang""

func InitMqttClient() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.MqttBrokerUrl, config.MqttBrokerPort))
	opts.SetClientID(config.MqttClientId)
	opts.SetUsername(config.MqttUserName)
	opts.SetPassword(config.MqttPassword)
	opts.SetDefaultPublishHandler(messageSubHandler)
	opts.SetOrderMatters(false)
	opts.SetCleanSession(false)
	opts.SetAutoReconnect(true)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	client.Disconnect(100)
	return nil
}
