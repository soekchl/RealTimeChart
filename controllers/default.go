package controllers

import (
	"math/rand"
	"sync"
	"time"

	"github.com/astaxie/beego"
	. "github.com/soekchl/myUtils"
)

type MainController struct {
	beego.Controller
}

var (
	dataList      []int32 // 目前保存110个
	dataListMutex sync.RWMutex
	max_length          = 100
	deviation     int32 = 5
)

func init() {
	var (
		r   int32 = 0
		min int32 = 1
		max int32 = 100
	)
	rand.Seed(time.Now().UnixNano())
	dataListMutex.Lock()
	defer dataListMutex.Unlock()

	for i := 0; i < max_length; i++ {
		r = rand.Int31n(max-min+1) + min
		dataList = append(dataList, r)
		max = r + deviation
		min = r - deviation
		Debug(max, " ", min)
	}

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				getOne()
			}
		}
	}()
}

func (c *MainController) Get() {
	c.Data["initData"] = getNewData() // 初始化数据
	c.TplName = "index.html"
}

func getOne() {
	dataListMutex.Lock()
	defer dataListMutex.Unlock()

	max := dataList[max_length-1] + deviation
	min := dataList[max_length-1] - deviation
	dataList = append(dataList, rand.Int31n(max-min+1)+min)
	Debug(max, " ", min)
	dataList = dataList[1:]
	sendWebSocket(EncodeInt32(dataList[max_length-1]))
}

func getNewData() (arr []int32) {
	dataListMutex.RLock()
	defer dataListMutex.RUnlock()
	return dataList
}

func EncodeInt32(n int32) (b []byte) {
	b = make([]byte, 4)
	b[0] = byte(n & 0xFF)
	b[1] = byte((n >> 8) & 0xFF)
	b[2] = byte((n >> 16) & 0xFF)
	b[3] = byte((n >> 24) & 0xFF)
	return
}
