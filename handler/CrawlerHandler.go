package handler

import "C"
import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"os"
	"strings"
	"thz-spider/Logger"
	"thz-spider/constant"
)

type Handler struct {
	C *colly.Collector
	Q *queue.Queue
}

func init() {

}
func (h Handler) ThemeHandler() {
	h.C.OnHTML("div[id='category_166'] a[href^='forum-']", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		link := e.Attr("href")
		Logger.Info.Printf("Link found: %q -> %s\n", e.Text, link)
		h.Q.AddURL(e.Request.AbsoluteURL(link))
	})
}
func (h Handler) TitleHandler() {
	h.C.OnHTML("div[class='bm_c'] a[class='s xst']", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		link := e.Attr("href")
		// Print link
		Logger.Info.Printf("Link found: %q -> %s\n", e.Text, link)
		//h.Q.AddURL(e.Request.AbsoluteURL(link))
		h.Q.AddURL(e.Request.AbsoluteURL(link) + "?name=" + e.Text)
	})
}
func (h Handler) NextPageHandler() {
	h.C.OnHTML("div[id='pgt'] div[class='pg'] a[class='nxt'][href^='forum-']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		Logger.Info.Printf("下一页: %q -> %s\n", e.Text, link)
		h.Q.AddURL(e.Request.AbsoluteURL(link))
	})
}
func (h Handler) ImageHandler() {
	h.C.OnHTML("img[file^='http://pic.thzpic.com/forum']", func(e *colly.HTMLElement) {
		src := e.Attr("file")

		// Print link
		Logger.Info.Printf("图片地址: %s\n", src)
		h.Q.AddURL(src + "?name=" + e.Request.URL.Query().Get("name"))
	})
}
func (h Handler) TorrentHandler() {
	h.C.OnHTML("p[class='attnm'] a", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		aid := strings.TrimPrefix(href, constant.TORRENT_PREFIX)
		Logger.Info.Printf("种子地址: %s\n", constant.TORRENT_DOWNLOAD_URI+aid)
		h.Q.AddURL(constant.TORRENT_DOWNLOAD_URI + aid + "&name=" + e.Request.URL.Query().Get("name"))
	})
}
func (h Handler) VisitingLogHandler() {
	h.C.OnRequest(func(r *colly.Request) {

		Logger.Info.Printf("Visiting:%s\n", r.URL.String())
	})
}
func (h Handler) SaveHandler() {
	h.C.OnResponse(func(r *colly.Response) {
		path := constant.DIR_PATH + r.Request.URL.Query().Get("name") + "/"
		if strings.Index(r.Headers.Get("Content-Type"), "image") > -1 {
			os.MkdirAll(path, os.ModePerm)
			r.Save(path + r.FileName()[7:43])
			Logger.Info.Printf("图片保存: %s\n", r.FileName()[7:43])
			return
		}
		if strings.Index(r.Headers.Get("Content-Type"), "application/octet-stream") > -1 {
			os.MkdirAll(path, os.ModePerm)
			r.Save(path + r.FileName())
			Logger.Info.Printf("种子保存: %s\n", r.FileName())
			return
		}
	})
}
