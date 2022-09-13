package database

import (
	"backEnd/config"
	"backEnd/dto"
	"backEnd/helper"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"strconv"
	"sync"
)

var database *sql.DB
var writeStatement *sql.Stmt

var mtx = &sync.Mutex{}

func Init() error {
	var tempDatabase *sql.DB
	var err error
	if config.DatabasePath == "" {
		config.DatabasePath = ":memory:"
		tempDatabase, err = sql.Open("sqlite3", config.DatabasePath)
	} else {
		tempDatabase, err = sql.Open("sqlite3", config.DatabasePath+config.DATABASE_NAME)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("sqlite open error")
		return err
	}
	statement, _ := tempDatabase.Prepare("CREATE TABLE IF NOT EXISTS messages(" +
		"topic STRING NOT NULL," +
		"payload blob" +
		")")
	_, err = statement.Exec()
	if err != nil {
		log.Fatal().Err(err).Msg("Sqlite DDL error")
		return err
	}
	database = tempDatabase
	writeStatement, err = database.Prepare("INSERT INTO messages (topic, payload) VALUES (?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("sqlite write statement prepare error")
	}

	return nil
}

func InsertChannelHandler() {
	for {
		msg := <-helper.InsertDbChannel
		mtx.Lock()
		log.Info().Msg("initializing sqlite insert")
		_, err := writeStatement.Exec(msg.TopicName, msg.Payload)
		if err != nil {
			log.Error().Err(err).Msg("sqlite write error")
		}
		mtx.Unlock()
		if log.Trace().Enabled() {
			log.Trace().
				Interface("data", msg).
				Msg("sqlite insert successful")
		} else {
			log.Debug().Msg("sqlite insert successful")
		}

	}
}

func RetrieveMessages() bool {
	mtx.Lock()
	log.Info().Msg("initializing retrieve messages handler")
	rows, err := database.Query("SELECT OID, topic, payload FROM messages")
	if err != nil {
		log.Error().Err(err).Msg("sqlite query statement prepare error")
	}
	tempCount := 0
	for rows.Next() {
		var topicName string
		var payload []byte
		var id int
		err := rows.Scan(&id, &topicName, &payload)
		if err != nil {
			log.Error().Err(err).Msg("sqlite row read error")
			continue
		} else {
			if log.Trace().Enabled() {
				log.Trace().
					Interface("Data:", dto.DbChannelStruct{
						OID:       id,
						TopicName: topicName,
						Payload:   payload,
					}).Msg("Data retrieved From Sqlite")
			} else {
				log.Debug().Msg("sqlite retrieve successfully")
			}

			helper.DatabaseChannel <- dto.DbChannelStruct{
				OID:       id,
				TopicName: topicName,
				Payload:   payload,
			}
		}
		tempCount++
		if tempCount > config.InvokeMaxGoRoutines {
			_ = rows.Close()
			log.Debug().Msg("partial retrieve message done")
			mtx.Unlock()
			return true
		}
	}

	_ = rows.Close()
	log.Debug().Msg(" full retrieve message done")
	mtx.Unlock()

	return false
}

func DeleteDbChannelHandler() {
	del, err := database.Prepare("DELETE FROM messages where OID == ?")
	if err != nil {
		log.Error().Err(err).Msg("sqlite delete statement prepare error")
	}

	for {
		oid := <-helper.DeleteDbChannel
		mtx.Lock()
		_, err = del.Exec(strconv.Itoa(oid))
		if err != nil {
			log.Error().Err(err).
				Int("Oid", oid).
				Msg("sqlite database delete error")
			mtx.Unlock()
			return
		}
		mtx.Unlock()
		log.Debug().
			Int("Oid", oid).
			Msg("sqlite delete done")
	}

}
