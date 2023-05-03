package main

import (
	"flag"
	"time"

	"github.com/fatih/color"
)

var (
	dirName     = flag.String("-d", "img", "图片文件夹")
	refreshTime = flag.Uint("-t", 500, "刷新时间，单位毫秒")
)

// 从用户端获取脚本的启动参数
func getManagerConfig() ManagerConfig {
	flag.Parse()
	flag.Args()

	dirName := *dirName

	refreshTime := *refreshTime

	color.Green("图片文件夹: %s", dirName)
	color.Green("刷新时间: %d", refreshTime)

	if refreshTime == 0 {
		refreshTime = 500
	}

	return ManagerConfig{RefreshTime: time.Duration(refreshTime) * time.Millisecond, ImgPath: dirName}
}

func main() {
	cfg := getManagerConfig()
	manager, err := NewManager(cfg)
	if err != nil {
		panic(err)
	}

	manager.work()
}
