package target

type FreeBufData struct {
	Data []struct {
		BugTitle string `json:"bug_title"`
	} `json:"data"`
}

type FreeBuf struct {
	urls   map[TargetData][2]string
	header map[string]string
	name   string
}

func (b *FreeBuf) New() {
	urls := make(map[TargetData][2]string)
	header := make(map[string]string)
	header["referer"] = "https://www.freebuf.com"
	header["user-agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

	urls[new(FreeBufData)] = [2]string{"GET", "https://www.freebuf.com/fapi/frontend/home/clipped?page=1"}
	b.name = "freebuf"
	b.header = header
	b.urls = urls
}

func (b *FreeBuf) Urls() map[TargetData][2]string {
	return b.urls
}

func (b *FreeBuf) Name() string {
	return b.name
}

func (b *FreeBuf) Header() map[string]string {
	return b.header
}

func (b *FreeBufData) Decode() []string {
	var code []string

	for _, list := range b.Data {
		code = append(code, list.BugTitle)
	}

	return code
}
