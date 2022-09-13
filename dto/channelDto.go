package dto

type DefaultChannelStruct struct {
	TopicName string
	Payload   []byte
}

type DbChannelStruct struct {
	OID       int
	TopicName string
	Payload   []byte
}
