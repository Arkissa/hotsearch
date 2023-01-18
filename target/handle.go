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
	Urls() map[TargetData][2]string
	Header() map[string]string
	Name() string
	New()
}

type TargetData interface {
	Decode() []string
}

type TargetHandle struct {
	targets []Target
	b       *pool.BufferPool
	c       *pool.ClientPool
}

func NewTargets(b *pool.BufferPool, c *pool.ClientPool) *TargetHandle {
	var targets []Target
	targets = append(targets, new(Blibili), new(WeiBo), new(DouYin), new(WeiBu), new(FreeBuf))

	return &TargetHandle{
		targets: targets,
		b:       b,
		c:       c,
	}
}

func (t *TargetHandle) Do() map[string][]string {
    datas := make(map[string][]string)
	for i := 0; i < len(t.targets); i++ {
		t.targets[i].New()
		urls := t.targets[i].Urls()
		name := t.targets[i].Name()
		header := t.targets[i].Header()

		for targetData, url := range urls {
			log.LogPut("[INFO] Start Request %s %s\n", name, url[1])
			t.c.Signal <- struct{}{}
			c := <-t.c.Client
			client := c.Get().(*http.Client)
			request := utils.NewRequest(url[0], url[1], "", client)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

			request.Ctx = ctx
			request.Header = header

			t.b.Signal <- struct{}{}
			p := <-t.b.Buffer
			buffer := p.Get().(*bytes.Buffer)
			buffer.Reset()

			_, byteBody := request.Do(buffer)
			cancel()
			c.Put(client)
			client = nil

			json := new(utils.JsonDate)
			json.Date = byteBody
			json.Decode = targetData
			json.Decoder()
			datas[name] = append(datas[name], targetData.Decode()...)

			p.Put(buffer)
			buffer = nil
		}
	}

	return datas
}
