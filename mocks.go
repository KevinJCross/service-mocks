package server_mocks

import "net/http"

type Mocks struct{}

func New() *Mocks {
	return &Mocks{}
}

func (*Mocks) Handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"message":"hello"}`))
		defer func() { _ = request.Body.Close() }()
	}
}
