# BliveDataExportToExcel

## 环境准备

系统版本  Windows 11 24H2

Go version go1.23.0 windows/amd64

安装Fyne UI，Fyne 需要 3 个基本元素，Go 工具（至少为 1.12 版）、C 编译器（用于连接系统图形驱动程序）和系统图形驱动程序。

```
$ go get fyne.io/fyne/v2@latest
$ go install fyne.io/fyne/v2/cmd/fyne@latest
```

 C 编译器推荐安装MinGW，右侧[Releases](https://github.com/AyumiLumi/BliveDataExportToExcel/releases)提供了压缩包，解压后，将/bin的路径添加到环境变量运行终端查看是否安装成功

```shell
g++ -v

gcc -v
```

## bilibili 直播弹幕 golang 库

基于github.com/Akegarasu/blivedm-go库修改，增加了简单的UI和导出数据到本地的功能，保证要有d盘，默认导出到d:/excel/下，感兴趣的可以修改代码使用Fyne提供的组件自定义保存路径。

## 安装

```shell
go get https://github.com/AyumiLumi/BliveDataExportToExcel
```

## 功能

### Home页面

输入主播房间号，输入你的Cookie，点击start按钮即可开始监听记录直播事件，点击取消或者关闭程序，会自动把程序运行期间的直播事件写入以时间+主播名为文件名的excel文件中，点击UI上的数据行，即可复制到系统粘贴板，直播事件做了如下图的分类

https://github.com/AyumiLumi/BliveDataExportToExcel/blob/main/image-20241126155655012.png

## 代码示例

```
package main

import (
    "context"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
    "github.com/AyumiLumi/BliveDataExportToExcel/export"
    "strconv"
    "time"
)

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("看看你的")

    // 输入框
    roomIdEntry := widget.NewEntry()
    roomIdEntry.SetPlaceHolder("主播直播房间号")
    cookieEntry := widget.NewMultiLineEntry()
    cookieEntry.SetPlaceHolder("你的登录 Cookie")

    // 状态标签
    statusLabel := widget.NewLabel("准备中...")

    // 消息存储与显示
    messages := make([]string, 0, 10000) // 最大存储 10000 条消息
    messageList := widget.NewList(
       func() int {
          return len(messages) // 返回消息数量
       },
       func() fyne.CanvasObject {
          return widget.NewLabel("") // 创建一个 Label
       },
       func(id widget.ListItemID, item fyne.CanvasObject) {
          item.(*widget.Label).SetText(messages[id])

       },
    )

    // 自动滚动
    scrollContainer := container.NewVScroll(messageList)
    scrollContainer.SetMinSize(fyne.NewSize(400, 600)) // 设置滚动容器大小

    // 消息追加函数
    appendMessage := func(msg string) {
       if len(messages) >= 10000 { // 超过 1000 条时删除最旧的一条
          messages = messages[:len(messages)-1]
       }
       messages = append([]string{msg}, messages...) // 将新消息加到顶部
       // 延迟调用 ScrollToTop 以确保刷新完成后滚动
       go func() {
          // 延迟触发滚动到底部
          time.Sleep(60 * time.Millisecond)
          scrollContainer.ScrollToBottom() // 自动滚动到底部
       }()
    }

    // 控制按钮
    var cancel context.CancelFunc

    // 先声明 startButton
    startButton := widget.NewButton("Start", nil)

    // 设置 startButton 的功能
    startButton.OnTapped = func() {
       roomId := roomIdEntry.Text
       cookie := cookieEntry.Text
       id, _ := strconv.Atoi(roomId)

       if roomId == "" || cookie == "" {
          statusLabel.SetText("请填写房间号和登录 Cookie,没有Cookie将无法完整显示信息！")
          //return
       }

       // 开始任务
       statusLabel.SetText("任务进行中...")
       startButton.Disable()

       // 创建上下文
       ctx, cancelFunc := context.WithCancel(context.Background())
       cancel = cancelFunc

       eventChan := make(chan string)
       cancelChan := make(chan struct{})

       go func() {
          export.ExportExcel(ctx, id, cookie, eventChan, cancelChan)
       }()

       go func() {
          for {
             select {
             case msg := <-eventChan:
                appendMessage(msg)    // 新消息加入滚动列表
                messageList.Refresh() // 刷新列表

             case <-cancelChan:
                statusLabel.SetText("任务已完成") // 更新按钮状态
                startButton.Enable()
                return
             }
          }
       }()
    }

    // Cancel按钮
    cancelButton := widget.NewButton("Cancel", func() {
       if cancel != nil {
          cancel()                           // 停止监听
          statusLabel.SetText("取消中，保存文件...") // 更新状态
       }
    })

    // 版权声明
    copyrightLabel := widget.NewLabel("© 2024 By 滨崎步的水钻麦克风. All rights reserved.")
    copyrightContainer := container.NewCenter(copyrightLabel) // 使用 Center 来居中显示
    // 设置窗口关闭时保存文件
    myWindow.SetCloseIntercept(func() {
       if cancel != nil {
          cancel()                        // 取消监听
          statusLabel.SetText("保存文件中...") // 更新状态
          // 执行保存文件的操作
          // You might need to replace this with your actual save logic
          // Example: export.ExportExcel(ctx, id, cookie, eventChan, cancelChan)
       }
       myWindow.Close() // 关闭窗口
    })

    // 居中显示按钮并加入间隔
    buttonsContainer := container.NewCenter(
       container.NewHBox(
          startButton,
          cancelButton,
       ),
    )

    // 布局
    content := container.NewVBox(
       widget.NewLabel("房间号"),
       roomIdEntry,
       widget.NewLabel("登录 Cookie"),
       cookieEntry,
       statusLabel,     // 按钮状态
       scrollContainer, // 消息滚动列表
       buttonsContainer,
       copyrightContainer, // 版权声明
    )

    myWindow.SetContent(content)
    myWindow.Resize(fyne.NewSize(500, 700))
    myWindow.ShowAndRun()
}
```
