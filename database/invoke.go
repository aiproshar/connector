package database

import (
	"backEnd/config"
	"backEnd/syncsys"
	"github.com/rs/zerolog/log"
	"time"
)

func DbInvoke() {
	for {
		err := database.Ping()
		if err != nil {
			log.Error().Err(err).Msg("database ping error")
		}
		cont := RetrieveMessages()
		if cont == false {

			log.Debug().Msg("Database is empty, DbInvoke going to sleep")
			time.Sleep(config.INVOKE_LATENCY_SEC * time.Second)
			log.Trace().Msg("DbInvoke woken from sleep")
		} else {
			for syncsys.ReadCount() > config.InvokeMinGoRoutines {
				log.Debug().Msg("DbInvoke ")
				time.Sleep(time.Millisecond * 10)
			}

		}

	}
}
