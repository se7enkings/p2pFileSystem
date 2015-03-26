package transfer

type Message interface {
	Type() string
	Destination() string
	Payload() []byte
}
