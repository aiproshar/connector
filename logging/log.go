package logging

import (
	"backEnd/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func getWarnLvl() zerolog.Level {
	switch config.LoggingLvl {
	case config.PanicLvl:
		return zerolog.PanicLevel
	case config.FatalLvl:
		return zerolog.FatalLevel
	case config.ErrorLvl:
		return zerolog.ErrorLevel
	case config.WarnLvl:
		return zerolog.WarnLevel
	case config.InfoLvl:
		return zerolog.InfoLevel
	case config.DebugLvl:
		return zerolog.DebugLevel
	case config.TraceLvl:
		return zerolog.TraceLevel
	}
	log.Warn().Msg(string("LOG_LVL:" + config.LoggingLvl + " not understood. Setting to infoLevel as default"))
	config.LoggingLvl = config.InfoLvl
	return zerolog.InfoLevel
}

func Init() {

	//if you want to customize the default logger
	//log.logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(getWarnLvl())
	log.Log().Str("Configured Logging level:", string(config.LoggingLvl)).Msg("")
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}
func PreInit() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Log().Msg("Pre Initialing Logging Package, Temporary logLvl: trace")

}
