# Onmyoji-GO
适用于 mac 的 阴阳师自动点击脚本

## 依赖
- opencv

### macos 安装 opencv
```shell
brew install pkg-config
brew install opencv
export PATH="/opt/homebrew/opt/opencv/bin:$PATH"
export PKG_CONFIG_PATH="/usr/local/opt/opencv/lib/pkgconfig"
```

## 使用
### 在文件夹下创建一个文件夹，放入需要检测的按钮截图，如`example`

![image-20230503213712677](https://raw.githubusercontent.com/Ccheers/pic/main/img/image-20230503213712677.png)

### 使用 QQ、微信、飞书、的截图工具截取你想要检测并点击的按钮，放置于你创建的文件夹下

注意：检测优先级为图片的字典序

### 打开一个终端，使用 ` onmyoji-go` 启动脚本 ，按照macos的提示赋予相关权限（不要用直接点击程序的方法打开，可能会让macos无法正常赋予权限）

```shell
# 查看帮助
./onmyoji-go -h
Usage of ./onmyoji-go:
  -d string
        图片文件夹 (default "img")
  -s uint
        停止时间，单位分钟 (default 90)
  -t uint
        刷新时间，单位毫秒 (default 500)

```

### 输入截图文件夹的名称，切换到游戏界面，开刷

![image-20230503213647166](https://raw.githubusercontent.com/Ccheers/pic/main/img/image-20230503213647166.png)




## 关于防封
首先，opencv的图片识别函数在运行过程中存在执行时间的不确定性，使每次点击的间隔时间存在随机性。
<img width="793" alt="image" src="https://user-images.githubusercontent.com/39732766/213117469-11227f93-b2e3-4306-bbd2-63b81ef28cbb.png">


其次，使用了随机函数，每次随机点击按钮内的一个点。

但是这种处理方式会造成一个问题：真人在操作的时候，每次点击的像素点并不会是真随机分布的，而使用随机函数点击，散点图则是真随机分布，非常不自然，很容易被检测到。

<img width="224" alt="iShot_2023-01-18_15 54 03" src="https://user-images.githubusercontent.com/39732766/213117628-a8d2364f-13e5-4321-8046-65e5fa849de4.png">

所以这里使用了正态分布生成随机数，接近于真人点击效果，如图是模拟250次结果：
<img width="621" alt="image" src="https://user-images.githubusercontent.com/39732766/213117804-85c26491-c889-43d7-b2e5-b923eda79d73.png">



## 稳定性
稳定操作阴阳师刷魂土一个月。

但是仍然存在被检测的可能性，使用本程序所造成的一切后果需由自己承担。

推荐挂2-3小时左右停一段时间，不要作死挂一晚上。

## thanks
- [auto-click](https://github.com/WinterBokeh/auto-click)