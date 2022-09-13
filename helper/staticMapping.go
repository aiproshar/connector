package helper

import (
	log "github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

// TopicMap
// static mapping from mqtt-topics to kafka-topics. No wildcard is allowed in static mapping
var TopicMap = make(map[string]string)

// RegExpToKafkaTopic
// Regular expression to actual mqtt wildCard topic Mapping
// Key: regExp , Value: mqtt-wildcard topic ("+" , "#")
var RegExpToKafkaTopic = make(map[string]string)
var RegExpToWildCard = make(map[string]string)

func InitTopicMap(list []string) {
	log.Debug().Msg("initializing initTopiMap map")
	for i := 0; i < len(list); i++ {
		if list[i] != "" {
			temp := strings.Split(list[i], ":")
			if len(temp) < 2 {
				log.Warn().Msg("unable to split substring using \":\", substring:" + list[i])
				continue
			}
			if temp[0] == "" || temp[1] == "" {
				log.Warn().Msg("topicName is empty, substring:" + list[i])
				continue
			}
			TopicMap[temp[0]] = temp[1]
			mapping := struct {
				Id         int    `json:"id"`
				MqttTopic  string `json:"mqttTopic"`
				KafkaTopic string `json:"kafkaTopic"`
			}{
				Id:         i + 1,
				MqttTopic:  temp[0],
				KafkaTopic: temp[1],
			}
			log.Trace().
				Interface("", mapping).
				Msg("static topic mapping")
		} else {
			log.Warn().Msg("empty split string found in TOPIC_LIST, Id:" + strconv.Itoa(i))
			continue
		}
	}
}

func InitWildCard(list []string) {
	log.Debug().Msg("initializing initTopiMap map")
	for i := 0; i < len(list); i++ {
		if list[i] != "" {
			temp := strings.Split(list[i], ":")
			if len(temp) < 2 {
				log.Warn().Msg("unable to split substring using \":\", substring:" + list[i])
				continue
			}
			if temp[0] == "" || temp[1] == "" {
				log.Warn().Msg("topicName is empty, substring:" + list[i])
				continue
			}
			mqttWildCardTopic := temp[0]
			kafKaTopicName := temp[1]
			log.Debug().Msg("converting mqtt wildcard topic:" + mqttWildCardTopic)
			regExp := wildCardTtoRegExp(mqttWildCardTopic)
			//saving actual wildCard topicName to map initiate mqtt listeners
			RegExpToWildCard[regExp] = mqttWildCardTopic
			RegExpToKafkaTopic[regExp] = kafKaTopicName
			mapping := struct {
				Id         int    `json:"id"`
				Wildcard   string `json:"wildcard"`
				RegExp     string `json:"regExp"`
				KafkaTopic string `json:"kafkaTopic"`
			}{
				Id:         i + 1,
				Wildcard:   mqttWildCardTopic,
				RegExp:     regExp,
				KafkaTopic: kafKaTopicName,
			}
			log.Trace().
				Interface("", mapping).
				Msg("wildcard topic mapping")

		} else {
			log.Warn().Msg("empty split string found in TOPIC_LIST, Id:" + strconv.Itoa(i))
			continue
		}
	}
}
