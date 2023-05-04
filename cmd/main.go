package main

import (
	"flag"
	"time"

	"github.com/fatih/color"
)

var (
	dirName     = flag.String("d", "img", "图片文件夹")
	refreshTime = flag.Uint("t", 500, "刷新时间，单位毫秒")
	stopTime    = flag.Uint("s", 90, "停止时间，单位分钟")
)

const (
	// KeyWordRoot 根场景关键字 文件名就是 .root
	KeyWordRoot = ".root"
	// KeyWordMatch 匹配场景关键字 {图片名}.{文件后缀}
	KeyWordMatch = "match"
	// KeyWordNext 下一步关键字 {按钮名}.{场景名}.{文件后缀}
	KeyWordNext = "next"
	// KeyWordBtns 按钮关键字 {按钮名}.{文件后缀}
	KeyWordBtns = "btns"
)

// 从用户端获取脚本的启动参数
func getManagerConfig() ManagerConfig {
	dirName := *dirName
	refreshTime := *refreshTime
	stopTime := *stopTime

	color.Green("图片文件夹: %s", dirName)
	color.Green("刷新时间: %d", refreshTime)
	color.Green("停止时间: %d", stopTime)

	if refreshTime == 0 {
		refreshTime = 500
	}

	return ManagerConfig{
		RefreshTime: time.Duration(refreshTime) * time.Millisecond,
		StopAt:      time.Now().Add(time.Duration(stopTime) * time.Minute),
		ImgPath:     dirName,
	}
}

func main() {
	flag.Parse()

	cfg := getManagerConfig()
	color.Green("当前配置参数:%+v", cfg)
	manager, err := NewManager(cfg)
	if err != nil {
		panic(err)
	}

	manager.work()
}
