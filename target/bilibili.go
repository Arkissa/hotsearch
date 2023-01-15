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
	Name   string
	Header map[string]string
	Urls   map[TargetData]string
}

func (b *Blibili) New() any {
	urls := make(map[TargetData]string)
	header := make(map[string]string)
	header["referer"] = "https://www.bilibili.com"
	header["user-agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

	urls[new(BlibiliData)] = "https://api.bilibili.com/x/web-interface/wbi/search/square?limit=50"
	b.Name = "bilibili"
	b.Header = header
	b.Urls = urls

	return b
}

func (b *BlibiliData) Decode() []string {
	var code []string

	for _, list := range b.Data.Trending.List {
		code = append(code, list.Keyword)
	}

	return code
}
