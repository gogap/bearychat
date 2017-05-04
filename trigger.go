package bearychat

type TriggerHandleFunc func(req *OutgoingRequest, msg *Message) (err error)

type Trigger interface {
	Handle(*OutgoingRequest, *Message) error
}
