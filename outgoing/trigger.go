package outgoing

type TriggerHandleFunc func(req *Request, resp *Response) (err error)

type Trigger interface {
	Handle(*Request, *Response) error
}
