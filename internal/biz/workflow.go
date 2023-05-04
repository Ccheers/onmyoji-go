package biz

import (
	"context"
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/Ccheers/onmyoji-go/internal/pkg/fnx"
	"github.com/fatih/color"
	"github.com/go-vgo/robotgo"
)

// Button 按钮
type Button interface {
	Name() string
	Match(ctx context.Context, screen image.Image) bool
	Click() error
	Next() *Scene
}

// Matcher 匹配器
type Matcher interface {
	Name() string
	Match(ctx context.Context, screen image.Image) bool
}

// Scene 场景
type Scene struct {
	name          string
	sceneMatchers []Matcher

	mu       sync.Mutex
	btns     []Button
	nextBtns []Button
}

func NewScene(name string) *Scene {
	return &Scene{
		name: name,
	}
}

func (x *Scene) AddMatcher(matchers ...Matcher) {
	x.sceneMatchers = append(x.sceneMatchers, matchers...)
}

func (x *Scene) Check() error {
	if len(x.sceneMatchers) == 0 {
		return fmt.Errorf("场景匹配器未装载")
	}
	return nil
}

func (x *Scene) AddBtn(btns ...Button) {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.btns = append(x.btns, btns...)
}

func (x *Scene) AddNextBtn(btns ...Button) {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.nextBtns = append(x.nextBtns, btns...)
}

func (x *Scene) match(ctx context.Context, screenImage image.Image) bool {
	for _, matcher := range x.sceneMatchers {
		if !matcher.Match(ctx, screenImage) {
			return false
		}
	}
	return true
}

type Workflow struct {
	allScenes []*Scene
	scene     *Scene
}

func NewWorkflow(scene *Scene, allScenes []*Scene) (*Workflow, error) {
	err := scene.Check()
	if err != nil {
		return nil, err
	}
	for _, aScene := range allScenes {
		err = aScene.Check()
		if err != nil {
			return nil, err
		}
	}
	return &Workflow{scene: scene, allScenes: allScenes}, nil
}

func (x *Workflow) Run(ctx context.Context) error {
	// 捕获当前屏幕
	var screenImg image.Image
	fnx.FnCost(ctx, fmt.Sprintf("%s:捕获并保存当前屏幕", x.scene.name), func(ctx context.Context) {
		screenImg = robotgo.CaptureImg()
	})

	// 错误场景修复
	if !x.scene.match(ctx, screenImg) {
		color.Yellow("当前场景(%s)不匹配，尝试修复", x.scene.name)
		x.autoSelectScene(ctx, screenImg)
	}

	for _, btn := range x.scene.btns {
		err := x.opBtn(ctx, screenImg, btn)
		if err != nil {
			return err
		}
	}

	for _, btn := range x.scene.nextBtns {
		match, err := x.opNextBtn(ctx, screenImg, btn)
		if err != nil {
			return err
		}
		if match {
			return nil
		}
	}

	return nil
}

func (x *Workflow) opBtn(ctx context.Context, screenImg image.Image, btn Button) error {
	var match bool
	fnx.FnCost(ctx, btn.Name()+":匹配", func(ctx context.Context) {
		match = btn.Match(ctx, screenImg)
	})
	if match {
		var err error
		fnx.FnCost(ctx, btn.Name()+":点击", func(ctx context.Context) {
			err = btn.Click()
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func (x *Workflow) opNextBtn(ctx context.Context, screenImg image.Image, btn Button) (bool, error) {
	var match bool
	fnx.FnCost(ctx, btn.Name()+":匹配", func(ctx context.Context) {
		match = btn.Match(ctx, screenImg)
	})
	if match {
		var err error
		fnx.FnCost(ctx, btn.Name()+":点击", func(ctx context.Context) {
			err = btn.Click()
		})
		if err != nil {
			return false, err
		}
		time.Sleep(time.Second * 2)
		fnx.FnCost(ctx, fmt.Sprintf("%s:下一步捕获并保存当前屏幕", x.scene.name), func(ctx context.Context) {
			screenImg = robotgo.CaptureImg()
		})

		next := btn.Next()
		if next.match(ctx, screenImg) {
			x.selectScene(next)
			return match, nil
		}
		color.Yellow("场景(%s)不匹配", next.name)
	}
	return match, nil
}

func (x *Workflow) autoSelectScene(ctx context.Context, screenImg image.Image) {
	for _, scene := range x.allScenes {
		if scene.match(ctx, screenImg) {
			x.selectScene(scene)
			return
		}
	}
}

func (x *Workflow) selectScene(scene *Scene) {
	color.Yellow("进入场景: %s", scene.name)
	x.scene = scene
}
