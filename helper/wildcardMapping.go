package helper

import (
	"container/list"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

// wildCardTtoRegExp
// Generates regular expression from a wildCard topic
// WhenEver we get a message from a mqtt broker, we need to find against which wildcard that topic was
func wildCardTtoRegExp(wildcard string) string {

	convertedRegExp := strings.ReplaceAll(wildcard, "/", "\\/")
	convertedRegExp = strings.ReplaceAll(convertedRegExp, "+", "[a-zA-Z0-9-_.]+")
	convertedRegExp = strings.ReplaceAll(convertedRegExp, "#", "[a-zA-Z0-9-_./]+")
	convertedRegExp = "^" + convertedRegExp + "$"
	log.Trace().
		Str("wildcard", wildcard).
		Str("regExp", convertedRegExp).
		Msg("mqtt wildcard to regExp converted")
	return convertedRegExp
}

// FindKafkaTopicFromMqttTopic
// converts MQTT topicName name to respective kafka topic
//
// It tries to match all regular expression. RegExp are generated from wildcards
// In ideal condition there will be exactly one match. If there's
// a collision it may be mapped to multiple regExp and respective wildcards.
// for example, default/+/app and default/v1/# all maps to default/v1/app.
// in a collision we log warn msg and return the last match
func FindKafkaTopicFromMqttTopic(topicName string) (string, error) {

	totalMatchCount := 0
	mappedName := ""
	lst := list.New()
	for regExp, kafKaTopic := range RegExpToKafkaTopic {
		match, err := regexp.MatchString(regExp, topicName)
		if err != nil {
			log.Error().Err(err).
				Str("regExp", regExp).
				Msg("regexp match error")
		} else {
			if match == true {
				totalMatchCount++
				mappedName = kafKaTopic
				wildCard := RegExpToWildCard[regExp]
				lst.PushBack(&wildCard)
			}
		}
	}
	if mappedName == "" {
		return mappedName, errors.New("topic map failed")
	}
	if lst.Len() > 1 {
		temp := make([]string, lst.Len())
		i := 0
		for e := lst.Front(); e != nil; e = e.Next() {
			temp[i] = *e.Value.(*string)
			i++
		}
		log.Warn().
			Str("MqttTopic", topicName).
			Strs("regExp", temp).
			Msg("mqttTopicName mapped to multiple wild cards")

	}
	return mappedName, nil
}
