package lb_context

import (
	"dynamicpath/lib/loadbalancer_api"
	"net/http"
	"sync"
)

var LBChannel chan ChannelMessage
var mtx sync.Mutex

const (
	MaxChannel         int    = 100000
	EventPduSessionAdd string = "PduSessionAdd"
	EventPduSessionDel string = "PduSessionDelete"
)

func init() {
	// init Pool
	LBChannel = make(chan ChannelMessage, MaxChannel)
}

type Response struct {
	Header http.Header
	Status int
	Body   interface{}
}

func NewResponse(code int, h http.Header, body interface{}) (ret *Response) {
	ret = &Response{}
	ret.Status = code
	ret.Header = h
	ret.Body = body
	return
}

type HttpResponseMessage struct {
	HTTPResponse *Response
}

/* Send HTTP Response to HTTP handler thread through HTTP channel, args[0] is response payload and args[1:] is Additional Value*/
func SendHttpResponseMessage(channel chan HttpResponseMessage, header http.Header, status int, body interface{}) {
	responseMsg := HttpResponseMessage{}
	responseMsg.HTTPResponse = NewResponse(status, header, body)
	channel <- responseMsg
}

type ChannelMessage struct {
	Event       string
	HttpChannel chan HttpResponseMessage // return Http response
	Value       interface{}              // input/request value
}
type PduSessionRequest struct {
	Supi        string
	SessionInfo loadbalancer_api.SessionInfo
}

func NewHttpChannelMessage(event string, value interface{}) (msg ChannelMessage) {
	msg = ChannelMessage{}
	msg.Event = event
	msg.HttpChannel = make(chan HttpResponseMessage)
	msg.Value = value
	return msg
}

func SendMessage(msg ChannelMessage) {
	mtx.Lock()
	LBChannel <- msg
	mtx.Unlock()
}
