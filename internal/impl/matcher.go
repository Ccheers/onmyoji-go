package impl

import (
	"context"
	"image"

	"github.com/Ccheers/onmyoji-go/internal/pkg/imgx"
	"github.com/vcaesar/gcv"
)

type ScreenMatcherImpl struct {
	name   string
	target image.Image
}

func NewScreenMatcherImpl(name string, target image.Image) *ScreenMatcherImpl {
	return &ScreenMatcherImpl{name: name, target: target}
}

func (x *ScreenMatcherImpl) Name() string {
	return x.name
}

func (x *ScreenMatcherImpl) Match(ctx context.Context, screen image.Image) bool {
	// 把image.Image统一转换成image.RGBA，不然会断言失败
	// 进行图像识别
	_, num, _, _ := gcv.FindImg(imgx.Jpg2RGBA(x.target), imgx.Jpg2RGBA(screen))
	return num > 0.9
}
