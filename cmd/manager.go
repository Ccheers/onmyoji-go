package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Ccheers/onmyoji-go/internal/biz"
	"github.com/Ccheers/onmyoji-go/internal/impl"
	"github.com/Ccheers/onmyoji-go/internal/pkg/imgx"
	"github.com/fatih/color"
)

type ManagerConfig struct {
	RefreshTime time.Duration //刷新时间，单位毫秒
	StopAt      time.Time     // 在某个时间停止运行
	ImgPath     string        //样本文件夹
}

func (x *ManagerConfig) Check() error {
	if x.RefreshTime == 0 {
		return fmt.Errorf("刷新时间不能为0")
	}

	if x.StopAt.Before(time.Now()) {
		return fmt.Errorf("停止时间不能在当前时间之前")
	}

	_, err := os.Stat(x.ImgPath)
	if err != nil {
		return fmt.Errorf("图片文件夹不存在")
	}
	return nil
}

type ImgInfo struct {
	path    string //图片路径
	ImgMaxX int
	ImgMaxY int
	Img     image.Image
}

type Manager struct {
	conf ManagerConfig

	workflow *biz.Workflow
}

func NewManager(cfg ManagerConfig) (*Manager, error) {
	err := cfg.Check()
	if err != nil {
		return nil, err
	}

	var root *biz.Scene
	sceneMap := make(map[string]*biz.Scene)
	dirs, err := os.ReadDir(cfg.ImgPath)
	if err != nil {
		return nil, err
	}
	for _, de := range dirs {
		if !de.IsDir() {
			continue
		}
		if de.Name() == ".DS_Store" {
			continue
		}
		modName := de.Name()
		path := filepath.Join(cfg.ImgPath, modName)

		color.Green("扫描模块: %s", modName)

		_, err = os.Stat(filepath.Join(path, KeyWordRoot))
		if err == nil {
			_, ok := sceneMap[modName]
			if !ok {
				sceneMap[modName] = biz.NewScene(modName)
			}
			root = sceneMap[modName]
		}

		err = formatScene(sceneMap, modName, filepath.Join(path, KeyWordMatch), KeyWordMatch)
		if err != nil {
			color.Red("格式化匹配场景失败 %s 文件地址: %s", err.Error(), filepath.Join(path, KeyWordMatch))
		}
		err = formatScene(sceneMap, modName, filepath.Join(path, KeyWordNext), KeyWordNext)
		if err != nil {
			color.Red("格式化下一步场景失败 %s 文件地址: %s", err.Error(), filepath.Join(path, KeyWordNext))
		}
		err = formatScene(sceneMap, modName, filepath.Join(path, KeyWordBtns), KeyWordBtns)
		if err != nil {
			color.Red("格式化按钮场景失败 %s 文件地址: %s", err.Error(), filepath.Join(path, KeyWordBtns))
		}
	}

	if root == nil {
		return nil, fmt.Errorf("根场景不存在")
	}

	allScenes := make([]*biz.Scene, 0, len(sceneMap))
	for _, scene := range sceneMap {
		allScenes = append(allScenes, scene)
	}

	flow, err := biz.NewWorkflow(root, allScenes)
	if err != nil {
		return nil, err
	}

	return &Manager{
		workflow: flow,
		conf:     cfg,
	}, nil
}

func formatScene(sceneMap map[string]*biz.Scene, modName string, imgPath string, keyword string) error {
	return filepath.Walk(imgPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Println("imgPath", imgPath)
		log.Println("path", path)
		if info.IsDir() {
			return nil
		}
		if info.Name() == ".DS_Store" {
			return nil
		}

		reader, err := os.Open(path)
		defer reader.Close()

		filename := info.Name()
		img := imgx.ReadPic(path)
		color.Green("已读取 %s", path)
		scene, ok := sceneMap[modName]
		if !ok {
			scene = biz.NewScene(modName)
			sceneMap[modName] = scene
		}

		switch keyword {
		case KeyWordMatch:
			scene.AddMatcher(impl.NewScreenMatcherImpl(filename, img))
		case KeyWordNext:
			fields := strings.Split(filename, ".")
			if len(fields) < 3 {
				return fmt.Errorf("文件名格式错误 %s 文件名格式需要是 {按钮名}.{场景名}.{文件后缀}", filename)
			}
			targetName := fields[1]
			target, ok := sceneMap[targetName]
			if !ok {
				target = biz.NewScene(targetName)
				sceneMap[targetName] = target
			}
			scene.AddNextBtn(impl.NewScreenBtn(filename, img, target))
		default:
			scene.AddBtn(impl.NewScreenBtn(filename, img, scene))
		}
		return nil
	})
}

func (m *Manager) work() {
	fmt.Println("============================================================================================================")
	color.Cyan(logo)
	fmt.Println("============================================================================================================")

	ctx, cancel := context.WithDeadline(context.Background(), m.conf.StopAt)
	defer cancel()

	bold := color.New(color.Bold).Add(color.FgGreen)
	bold.Println("开始运行脚本，请切换到游戏界面")
	for {
		select {
		case <-ctx.Done():
			color.Red("当前时间 %s 脚本运行结束", time.Now().Format("2006-01-02 15:04:05"))
			return
		default:
		}
		err := m.workflow.Run(ctx)
		if err != nil {
			color.Red("当前时间 %s 脚本运行出错：%s", time.Now().Format("2006-01-02 15:04:05"), err)
		}
		time.Sleep(m.conf.RefreshTime)
	}
}
