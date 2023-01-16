package pool

import (
	"bytes"
	"sync"
	"time"
)

type BufferPool struct {
	Signal chan struct{}
	Buffer chan *sync.Pool
}

func NewBufferPool() *BufferPool {
	pool := &BufferPool{
		Signal: make(chan struct{}, 1),
		Buffer: make(chan *sync.Pool, 10),
	}

	go pool.buffers()

	return pool
}

func (p *BufferPool) buffers() {
	bufferPool := sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 512))
		},
	}

	for range p.Signal {
		p.Buffer <- &bufferPool
	}
}

func (p *BufferPool) Close() {
	for len(p.Buffer) != 0 || len(p.Signal) != 0 {
		time.Sleep(1e6)
	}

    close(p.Buffer)
    close(p.Signal)
}
