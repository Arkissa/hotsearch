package pool

import (
	"net/http"
	"sync"
	"time"
)

type ClientPool struct {
	Signal chan struct{}
	Client chan *sync.Pool
}

func NewClient() *ClientPool {
	pool := &ClientPool{
		Signal: make(chan struct{}, 1),
		Client: make(chan *sync.Pool, 1),
	}

	go pool.client()

	return pool
}

func (p *ClientPool) client() {
	clientPool := sync.Pool{
		New: func() any {
			return new(http.Client)
		},
	}

	for range p.Signal {
		p.Client <- &clientPool
	}
}

func (p *ClientPool) Close() {
	for len(p.Client) != 0 || len(p.Signal) != 0 {
		time.Sleep(1e6)
	}

	close(p.Client)
	close(p.Signal)
}
