package fnx

import (
	"context"
	"time"

	"github.com/fatih/color"
)

func FnCost(ctx context.Context, fnName string, fn func(ctx context.Context)) {
	start := time.Now()
	fn(ctx)
	cost := time.Since(start)
	color.White("函数 %s 耗时：%v", fnName, cost)
}
