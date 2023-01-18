package target

type WeiBo struct {
	urls   map[TargetData][2]string
	header map[string]string
	name   string
}

type WeiBoData struct {
	Ok   int `json:"ok"`
	Data struct {
		Cards []struct {
			CardGroup []struct {
				Desc string `json:"desc"`
			} `json:"card_group"`
		} `json:"cards"`
	} `json:"data"`
}

func (w *WeiBo) New() {

	urls := make(map[TargetData][2]string)
	header := make(map[string]string)
	header["referer"] = "https://m.weibo.cn/search?containerid=231583"
	header["user-agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

	urls[new(WeiBoData)] = [2]string{"GET", "https://m.weibo.cn/api/container/getIndex?containerid=106003%26filter_type%3Drealtimehot&title=%E5%BE%AE%E5%8D%9A%E7%83%AD%E6%90%9C"}
	w.name = "weibo"
	w.header = header
	w.urls = urls
}

func (b *WeiBo) Urls() map[TargetData][2]string {
	return b.urls
}

func (b *WeiBo) Name() string {
	return b.name
}

func (b *WeiBo) Header() map[string]string {
	return b.header
}

func (w *WeiBoData) Decode() []string {
	var code []string

	for _, cards := range w.Data.Cards {
		for _, card := range cards.CardGroup {
			code = append(code, card.Desc)
		}
	}

	return code
}
