package target

type DouYinData struct {
	Data struct {
		WordList []struct {
			Word string `json:"word"`
		} `json:"word_list"`
		TrendingList []struct {
			Word string `json:"word"`
		} `json:"trending_list"`
	} `json:"data"`
}

type DouYin struct {
	name   string
	header map[string]string
	urls   map[TargetData]string
}

func (d *DouYin) New() {
	urls := make(map[TargetData]string)
	header := make(map[string]string)
	header["referer"] = "https://www.douyin.com/"
	header["user-agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

	urls[new(DouYinData)] = "https://www.douyin.com/aweme/v1/web/hot/search/list/?device_platform=webapp&aid=6383&channel=channel_pc_web&detail_list=1&source=6&pc_client_type=1&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=2240&screen_height=1400&browser_language=zh-CN&browser_platform=Linux+x86_64&browser_name=Chrome&browser_version=108.0.0.0&browser_online=true&engine_name=Blink&engine_version=108.0.0.0&os_name=Linux&os_version=x86_64&cpu_core_num=16&device_memory=8&platform=PC&downlink=3.45&effective_type=4g&round_trip_time=0&webid=7185133404112143927&msToken=kh7EMmHldNYKviiTCUnp3a-Pk9GrGm_OKitQ5lw0T2c2BEiTgobhLA91zoQPMz-bJYEUGETzKi-6C4neZUbZIvHC4ez6Sbi3Y8dukcb49nePCi9gOELV&X-Bogus=DFSzswVLKTTANnFXSDsSo1XzGAXQ"
	d.name = "douyin"
	d.header = header
	d.urls = urls
}

func (b *DouYin) Urls() map[TargetData]string {
	return b.urls
}

func (b *DouYin) Name() string {
	return b.name
}

func (b *DouYin) Header() map[string]string {
	return b.header
}

func (d *DouYinData) Decode() []string {
	var code []string

	for _, wordList := range d.Data.WordList {
		code = append(code, wordList.Word)
	}

	for _, trendingList := range d.Data.TrendingList {
		code = append(code, trendingList.Word)
	}

	return code
}
