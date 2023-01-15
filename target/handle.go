package target

import (
	"bytes"
	"context"
	"hotsearch/log"
	"hotsearch/pool"
	"hotsearch/utils"
	"net/http"
	"time"
)

type Target interface {
	New() any
}

type TargetData interface {
	Decode() []string
}

type TargetHandle struct {
	targets []Target
	header  map[string]string
	data    map[string][]string
	urls    map[TargetData]string
	name    string
	b       *pool.BufferPool
	c       *pool.ClientPool
}

func NewTargets(b *pool.BufferPool, c *pool.ClientPool) *TargetHandle {

	var targets []Target
	targets = append(targets, new(Blibili), new(WeiBo), new(DouYin), new(WeiBu))

	return &TargetHandle{
		targets: targets,
		b:       b,
		c:       c,
	}
}

func (t *TargetHandle) Do() map[string][]string {
	t.data = make(map[string][]string)
	for i := 0; i < len(t.targets); i++ {
		switch v := t.targets[i].New().(type) {
		case *Blibili:
			t.name = v.Name
			t.header = v.Header
			t.urls = v.Urls
		case *WeiBo:
			t.name = v.Name
			t.header = v.Header
			t.urls = v.Urls
		case *DouYin:
			t.name = v.Name
			t.header = v.Header
			t.urls = v.Urls
		case *WeiBu:
			t.name = v.Name
			t.header = v.Header
			t.urls = v.Urls
		}

		for body, url := range t.urls {
			log.LogPut("[INFO] Start Request %s %s\n", t.name, url)
			t.c.Signal <- struct{}{}
			c := <-t.c.Client
			client := c.Get().(*http.Client)
			request := utils.NewRequest("GET", url, "", client)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			request.Ctx = ctx
			request.Header = t.header

			t.b.Signal <- struct{}{}
			p := <-t.b.Buffer
			buffer := p.Get().(*bytes.Buffer)
			buffer.Reset()

			_, byteBody := request.Do(buffer)

			json := new(utils.JsonDate)

			json.Date = byteBody
			json.Decode = body
			json.Decoder()
			t.data[t.name] = append(t.data[t.name], body.Decode()...)

			c.Put(client)
			client = nil
			p.Put(buffer)
			buffer = nil
		}
	}

	return t.data
}
