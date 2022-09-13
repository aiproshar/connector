package helper

import (
	"backEnd/dto"
	"github.com/rs/zerolog/log"
)

var DefaultChanel chan dto.DefaultChannelStruct
var DatabaseChannel chan dto.DbChannelStruct
var DeleteDbChannel chan int
var InsertDbChannel chan dto.DefaultChannelStruct

func InitChannels() {
	log.Debug().Msg("creating channels")
	DefaultChanel = make(chan dto.DefaultChannelStruct)
	DatabaseChannel = make(chan dto.DbChannelStruct)
	DeleteDbChannel = make(chan int)
	InsertDbChannel = make(chan dto.DefaultChannelStruct)
	log.Debug().Msg("channels created")
}
