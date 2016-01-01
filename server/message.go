package server

// Contains a message
type Message struct {
	Content string
	Host    string
	Type    MessageType
}

type MessageType uint8

const CONNECTION_ESTABLISHED MessageType = 1
const CONNECTION_LOST MessageType = 2
const INBOUND MessageType = 3
const OUTBOUND MessageType = 4


func NewMessage(host string, content string, messageType MessageType) *Message {
	return &Message{
		Content: content,
		Host: host,
		Type: messageType,
	}
}

func (message *Message) String() string {
	return "[" + message.Host + "] " + message.Content
}
