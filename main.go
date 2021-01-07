package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"os"
	"reflect"
	"runtime"
	"thz-spider/Logger"
	"thz-spider/handler"
	"time"
)

func main() {

	start := time.Now().UnixNano()
	threadNum := runtime.NumCPU() * 2
	fmt.Printf("CPU核心数:%d \n", runtime.NumCPU())
	c := colly.NewCollector(
		colly.AllowedDomains("taohuazu7.com", "thz.cc", "pic.thzpic.com"),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: threadNum})
	storage := &redisstorage.Storage{
		Address:  os.Args[1],
		Password: os.Args[2],
		DB:       0,
		Prefix:   "thz",
	}
	err := c.SetStorage(storage)
	if err != nil {
		panic(err)
	}
	if os.Args[3] == "master" {
		storage.Clear()
	}
	q, _ := queue.New(8, storage)
	v := reflect.ValueOf(handler.Handler{c, q})
	t := reflect.TypeOf(handler.Handler{c, q})
	for i := 0; i < v.NumMethod(); i++ {
		Logger.Info.Print(t.Method(i).Name)
		v.Method(i).Call([]reflect.Value{})
	}
	q.AddURL("http://taohuazu7.com")

	for !q.IsEmpty() {
		q.Run(c)
		c.Wait()
	}
	duration := (time.Now().UnixNano() - start) / 1e6
	fmt.Printf("总耗时:%d毫秒", duration)
}
