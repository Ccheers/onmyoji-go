package impl

import (
	"context"
	"image"

	"github.com/Ccheers/onmyoji-go/internal/biz"
	"github.com/Ccheers/onmyoji-go/internal/pkg/imgx"
	"github.com/Ccheers/onmyoji-go/internal/pkg/mathx"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
)

type ScreenBtn struct {
	name   string
	target image.Image
	X, Y   int
	next   *biz.Scene
}

func NewScreenBtn(name string, target image.Image, next *biz.Scene) *ScreenBtn {
	return &ScreenBtn{name: name, target: target, next: next}
}

func (x *ScreenBtn) Name() string {
	return x.name
}

func (x *ScreenBtn) Match(ctx context.Context, screen image.Image) bool {
	// 把image.Image统一转换成image.RGBA，不然会断言失败
	// 进行图像识别
	_, num, _, pos := gcv.FindImg(imgx.Jpg2RGBA(x.target), imgx.Jpg2RGBA(screen))
	if num < 0.9 {
		return false
	}
	x.X = pos.X
	x.Y = pos.Y
	return true
}

func (x *ScreenBtn) Click() error {
	x.randomClick(x.X, x.Y, x.target.Bounds().Max.X, x.target.Bounds().Max.Y)
	return nil
}

// 根据img中找到的的temp左上角坐标，和temp的最大长宽，计算出一块可以点击的区域，并随机点击
func (x *ScreenBtn) randomClick(tempPosX, tempPosY, tempMaxX, tempMaxY int) {
	//用qq截图软件截下来的图，分辨率是真实分辨率二倍，所以除以2以对应真实分辨率
	tempPosX, tempPosY, tempMaxX, tempMaxY = tempPosX/2, tempPosY/2, tempMaxX/2, tempMaxY/2

	//计算按钮的中心点
	centerX, centerY := tempPosX+tempMaxX/2, tempPosY+tempMaxY/2

	//在中心点加或减offset就是随机坐标的上限和下限
	offsetX := tempMaxX / 2
	offsetY := tempMaxY / 2
	_, randomX := mathx.RandomNormalInt64(int64(centerX-offsetX), int64(centerX+offsetX), int64(centerX), 10)
	_, randomY := mathx.RandomNormalInt64(int64(centerY-offsetY), int64(centerY+offsetY), int64(centerY), 10)
	robotgo.MouseSleep = 10
	robotgo.Move(int(randomX), int(randomY))
	robotgo.Click("left")
}
func (x *ScreenBtn) Next() *biz.Scene {
	return x.next
}
