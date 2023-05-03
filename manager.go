package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
)

type ManagerConfig struct {
	RefreshTime time.Duration //刷新时间，单位毫秒
	ImgPath     string        //样本文件夹
}

type ImgInfo struct {
	path    string //图片路径
	ImgMaxX int
	ImgMaxY int
	Img     image.Image
}

type Manager struct {
	ManagerConfig

	ImgInfos []ImgInfo
}

func NewManager(cfg ManagerConfig) (*Manager, error) {
	var imgInfos []ImgInfo
	err := filepath.Walk(cfg.ImgPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Base(info.Name()) == ".DS_Store" {
			return nil
		}
		reader, err := os.Open(path)
		defer reader.Close()

		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			return fmt.Errorf("图片错误: %s, 图片名: %s", err, path)
		}

		color.Green("已读取 %s", path)
		imgInfos = append(imgInfos, ImgInfo{
			path:    path,
			ImgMaxX: im.Width,
			ImgMaxY: im.Height,
			Img:     readPic(path),
		})
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &Manager{
		ImgInfos:      imgInfos,
		ManagerConfig: cfg,
	}, nil
}

// 进行图像识别，在img中找temp，并返回在img中找到的的temp左上角坐标
func findTempPos(temp, img image.Image) (int, int, float32) {
	//把image.Image统一转换成image.RGBA，不然会断言失败
	_, num, _, pos := gcv.FindImg(jpg2RGBA(temp), jpg2RGBA(img))
	return pos.X, pos.Y, num
}

// 根据img中找到的的temp左上角坐标，和temp的最大长宽，计算出一块可以点击的区域，并随机点击
func randomClick(tempPosX, tempPosY, tempMaxX, tempMaxY int) {
	//用qq截图软件截下来的图，分辨率是真实分辨率二倍，所以除以2以对应真实分辨率
	tempPosX, tempPosY, tempMaxX, tempMaxY = tempPosX/2, tempPosY/2, tempMaxX/2, tempMaxY/2

	//计算按钮的中心点
	centerX, centerY := tempPosX+tempMaxX/2, tempPosY+tempMaxY/2

	//在中心点加或减offset就是随机坐标的上限和下限
	offsetX := tempMaxX / 2
	offsetY := tempMaxY / 2
	_, randomX := RandomNormalInt64(int64(centerX-offsetX), int64(centerX+offsetX), int64(centerX), 10)
	_, randomY := RandomNormalInt64(int64(centerY-offsetY), int64(centerY+offsetY), int64(centerY), 10)

	robotgo.MouseSleep = 10
	robotgo.Move(int(randomX), int(randomY))
	robotgo.Click("left")
}

func (m *Manager) work() {
	fmt.Println("============================================================================================================")
	color.Cyan(logo)
	fmt.Println("============================================================================================================")

	bold := color.New(color.Bold).Add(color.FgGreen)
	bold.Println("开始运行脚本，请切换到游戏界面")
	for {
		//捕获当前屏幕
		start := time.Now()
		screenImg := robotgo.CaptureImg()
		cost := time.Since(start)

		fmt.Println("_______________________________")
		fmt.Println("成功捕获并保存当前屏幕，耗时：", cost)

		//逐一匹配样板图片
		for _, img := range m.ImgInfos {
			start := time.Now()
			tempPosX, tempPosY, num := findTempPos(img.Img, screenImg)
			cost := time.Since(start)

			fmt.Print(" 正在匹配：", img.path, " 相似度：", num, " 匹配耗时：", cost)

			if num > 0.9 {
				start := time.Now()
				randomClick(tempPosX, tempPosY, img.ImgMaxX, img.ImgMaxY)
				cost := time.Since(start)
				bold.Println(" 匹配成功, 耗时：", cost)
				break
			}
			fmt.Println(" 匹配不到相似的图片")
		}
		time.Sleep(m.RefreshTime)
	}
}
