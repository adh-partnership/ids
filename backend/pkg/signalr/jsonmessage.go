package signalr

const JSON_PACKET_SEPARATOR = "\u001e" // ASCII Character 0x1E (Record Separator)

type MessageType int

const (
	TYPE_INVOCATION MessageType = iota + 1 // Starts at 1
	TYPE_STREAM_ITEM
	TYPE_COMPLETION
	TYPE_STREAM_INVOCATION
	TYPE_CANCELLATION
	TYPE_PING
	TYPE_CLOSE
)
