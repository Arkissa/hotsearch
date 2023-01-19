package target

import (
	"bytes"
	"context"
	"fmt"
	"hotsearch/log"
	"hotsearch/pool"
	"hotsearch/utils"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/dop251/goja"
)

type Cnvd struct {
	rpx     map[string]string
	header  map[string]string
	replace []string
	xpaths  []string
	js      *goja.Runtime
	b       *pool.BufferPool
	c       *pool.ClientPool
	format  string
	url     string
}

func NewCnvd(b *pool.BufferPool, c *pool.ClientPool) *Cnvd {
	rpx := make(map[string]string)
	rpx[`{var _0x\w+=window\[\S+\s\S+{return!!\[\];}`] = "{return false;"
	rpx[`var _0x\w+=window.{1,}return!!\[\];}}}`] = "return false;}"
	rpx[`setTimeout\(function\(\){document\[_0x\w+\(.{1,20}\]=`] = "return "
	rpx[`_0x\w+\[_0x\w+\(.{1,20}\]\(setTimeout,function\(\){document\[_0x\w+.{1,30}=`] = `return `
	rpx[`\);location\[.{1,},_0x\w+\);`] = ");"
	rpx[`location\[.{1,},_0x\w+\);`] = ""
	rpx[`go\({`] = `return go({`

	header := make(map[string]string)
	header["referer"] = "https://www.cnvd.org.cn/"
	header["user-agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

	var (
		replace = make([]string, 5)
		xpaths  = make([]string, 6)
	)

	replace[0] = "path=/; HttpOnly; SameSite=None; secure"
	replace[1] = "<script>document.cookie="
	replace[2] = "location.href=location.pathname+location.search</script>"
	replace[3] = "<script>"
	replace[4] = "</script>"

	format := "var cookie = ()=>{%s};cookie()"

	xpaths[0] = "/html/body/div[3]/div[3]/div/div[1]/ul/li/a"
	xpaths[1] = "/html/body/div[3]/div[3]/div/div[2]/ul/li/a"
	xpaths[2] = "/html/body/div[3]/div[4]/div[2]/ul/li/a"
	xpaths[3] = "/html/body/div[3]/div[5]/div/div[1]/ul/li/a"
	xpaths[4] = "/html/body/div[3]/div[7]/div/div[1]/ul/li/a"
	xpaths[5] = "/html/body/div[3]/div[7]/div/div[2]/ul/li/a"

	return &Cnvd{
		js:      goja.New(),
		rpx:     rpx,
		header:  header,
		replace: replace,
		format:  format,
		url:     "https://www.cnvd.org.cn/",
		xpaths:  xpaths,
		b:       b,
		c:       c,
	}
}

func (c *Cnvd) runJs(jsCode *string) string {
	jsReturn, err := c.js.RunString(*jsCode)
	if err != nil {
		log.LogOutErr("Run js err", err)
		return ""
	}
	return jsReturn.String()
}

func (c *Cnvd) strReplace(str *string) {
	for i := 0; i < len(c.replace); i++ {
		*str = strings.ReplaceAll(*str, c.replace[i], "")
	}
}

func (c *Cnvd) rgx(str *string) {
	for rpx, value := range c.rpx {
		r, err := regexp.Compile(rpx)
		if err != nil {
			log.LogOutErr("regexp set err", err)
			continue
		}
		*str = r.ReplaceAllString(*str, value)
	}
}

func (c *Cnvd) xpath(body io.Reader) []string {
	doc, err := htmlquery.Parse(body)
	if err != nil {
		log.LogOutErr("generate html err", err)
		return nil
	}

	var data []string
	for i := 0; i < len(c.xpaths); i++ {
		lables, err := htmlquery.QueryAll(doc, c.xpaths[i])
		if err != nil {
			log.LogOutErr("xpath err"+c.xpaths[i], err)
			continue
		}

		for _, lable := range lables {
			data = append(data, htmlquery.SelectAttr(lable, "title"))
		}

	}

	return data
}

func (c *Cnvd) Do() map[string][]string {
	var (
		cookie, body string
		num          int
		done         = false
		result       = make(map[string][]string)
	)

	c.c.Signal <- struct{}{}
	cl := <-c.c.Client
	client := cl.Get().(*http.Client)
	request := utils.NewRequest("GET", c.url, nil, client)
	request.Header = c.header
	log.LogPut("[INFO] Start Request cnvd")
request:
	c.b.Signal <- struct{}{}
	p := <-c.b.Buffer
	buffer := p.Get().(*bytes.Buffer)
	buffer.Reset()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	request.Ctx = ctx
	responseHeader, byteBody := request.Do(buffer)
	cancel()
	body = string(byteBody)

	if num == 0 {
		if responseHeader["Set-Cookie"] != nil {
			cookie = responseHeader["Set-Cookie"][0]
			c.strReplace(&cookie)
		}
	}

	if !strings.Contains(body, "<title>") {
		c.strReplace(&body)
		c.rgx(&body)
		if strings.Contains(body, "var _0x") {
			body = fmt.Sprintf(c.format, body)
		}
		body = c.runJs(&body)
		request.Header["Cookie"] = fmt.Sprintf("%s %s", cookie, body)
	} else {
		done = true
	}

	p.Put(buffer)
	buffer = nil

	if num >= 5 {
		done = true
	}

	if !done {
		num++
		goto request
	}

	cl.Put(client)
	client = nil
	result["cnvd"] = c.xpath(strings.NewReader(body))

	return result
}
