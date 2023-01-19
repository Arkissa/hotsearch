package target

type BlibiliData struct {
	Data struct {
		Trending struct {
			List []struct {
				Keyword string `json:"keyword"`
			} `json:"list"`
		} `json:"trending"`
	} `json:"data"`
}

type Blibili struct {
	urls   map[TargetData][3]any
	header map[string]string
	name   string
}

func (b *Blibili) New() {
	urls := make(map[TargetData][3]any)
	header := make(map[string]string)
	header["referer"] = "https://www.bilibili.com"
	header["user-agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

	urls[new(BlibiliData)] = [3]any{"GET", "https://api.bilibili.com/x/web-interface/wbi/search/square?limit=50", nil}
	b.name = "bilibili"
	b.header = header
	b.urls = urls
}

func (b *Blibili) Urls() map[TargetData][3]any {
	return b.urls
}

func (b *Blibili) Name() string {
	return b.name
}

func (b *Blibili) Header() map[string]string {
	return b.header
}

func (b *BlibiliData) Decode() []string {
	var code []string

	for _, list := range b.Data.Trending.List {
		code = append(code, list.Keyword)
	}

	return code
}
