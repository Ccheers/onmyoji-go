package mathx

import (
	"math"
	"math/rand"
	"time"
)

func RandFromRangeInt64(min int64, max int64) int64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Int63n(max-min) + min
}

// NormalFloat64 正态分布公式
func NormalFloat64(x int64, miu int64, sigma int64) float64 {
	randomNormal := 1 / (math.Sqrt(2*math.Pi) * float64(sigma)) * math.Pow(math.E, -math.Pow(float64(x-miu), 2)/(2*math.Pow(float64(sigma), 2)))
	//注意下是x-miu，我看网上好多写的是miu-miu是不对的
	return randomNormal
}

// RandomNormalInt64 正态分布随机数生产器：min:最小值，max:最大值，miu:期望值（均值），sigma:方差
func RandomNormalInt64(min int64, max int64, miu int64, sigma int64) (bool, int64) {
	if min >= max {
		return false, 0
	}
	if miu < min {
		miu = min
	}
	if miu > max {
		miu = max
	}
	var x int64
	var y, dScope float64
	for {
		x = RandFromRangeInt64(min, max)
		y = NormalFloat64(x, miu, sigma) * 100000
		dScope = float64(RandFromRangeInt64(0, int64(NormalFloat64(miu, miu, sigma)*100000)))
		//注意下传的是两个miu
		if dScope <= y {
			break
		}
	}
	return true, x
}
