package utils

import (
	"bytes"
	"context"
	"hotsearch/log"
	"io"
	"net/http"
)

type Request struct {
	Header  map[string]string
	request *http.Request
	client  *http.Client
	Ctx     context.Context
}

type Response struct {
	Body   []byte
	Header http.Header
}

func NewRequest(method, url string, body io.Reader, client *http.Client) *Request {
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		log.LogOutErr("Cread New Request Err", err)
		return nil
	}

	return &Request{
		request: request,
		client:  client,
	}
}

func (r *Request) Do(buffer *bytes.Buffer) (http.Header, []byte) {
	for key, value := range r.Header {
		r.request.Header.Set(key, value)
	}

	r.request.WithContext(r.Ctx)

	response, err := r.client.Do(r.request)
	if err != nil {
		log.LogOutErr("Request "+r.request.RequestURI, err)
		return nil, nil
	}
	defer response.Body.Close()

	_, err = io.Copy(buffer, response.Body)

	if err != nil {
		log.LogOutErr("Reader "+r.request.RequestURI+" response body", err)
		return nil, nil
	}

	return response.Header, buffer.Bytes()
}
