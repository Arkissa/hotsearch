package target

type FreeBufData struct {
	Data []struct {
		BugTitle string `json:"bug_title"`
	} `json:"data"`
}

type FreeBuf struct {
	Name   string
	Header map[string]string
	Urls   map[TargetData]string
}

func (b *FreeBuf) New() any {
	urls := make(map[TargetData]string)
	header := make(map[string]string)
	header["referer"] = "https://www.freebuf.com"
	header["user-agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

    urls[new(FreeBufData)] = "https://www.freebuf.com/fapi/frontend/home/clipped?page=1"
	b.Name = "freebuf"
	b.Header = header
	b.Urls = urls

	return b
}

func (b *FreeBufData) Decode() []string {
	var code []string

	for _, list := range b.Data {
		code = append(code, list.BugTitle)
	}

	return code
}
