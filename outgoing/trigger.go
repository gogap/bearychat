package outgoing

type Trigger interface {
	Handle(*Request, *Response) error
}
