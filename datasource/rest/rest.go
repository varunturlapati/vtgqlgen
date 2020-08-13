package rest

type RTX interface {
	/*
	Get(context.Context, string, ...interface{}) (http.Response, error)
	Post(context.Context, string, ...interface{}) (http.Response, error)
	Put(context.Context, string, ...interface{}) (http.Response, error)
	Delete(context.Context, string, ...interface{}) (bool, error)
	*/

}

func New(r string) *RestRequests {
	c := NewClient(r)
	return &RestRequests{rreq: c}
}

type RestRequests struct {
	rreq RTX
}

type RestClient struct {
	Url string
}

func NewClient(r string) *RestClient {
	return &RestClient{Url: r}
}