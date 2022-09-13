package config

import "github.com/rs/zerolog"

const TRUE = "true"
const FALSE = "false"

const DEVELOP = "DEVELOP"
const PRODUCTION = "PRODUCTION"
const DATABASE_NAME = "default.db"
const INVOKE_LATENCY_SEC = 10
const APP_VERSION = "v0.3.5"

const SASL_SSL = "SASL_SSL"
const SASL_MECHANISM_AOUTHBEARER = "OAUTHBEARER"
const SASL_OAUTHBEARER_METHOD_OIDC = "oidc"

const CONFIDENTIAL = "confidential"

type logLvl string

const PreLoggingLevel = zerolog.TraceLevel
const DefaultLoggingLvl = zerolog.InfoLevel

const (
	PanicLvl logLvl = "panic"
	FatalLvl logLvl = "fatal"
	ErrorLvl logLvl = "error"
	WarnLvl  logLvl = "warn"
	InfoLvl  logLvl = "info"
	DebugLvl logLvl = "debug"
	TraceLvl logLvl = "trace"
)
