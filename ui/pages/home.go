package pages

import (
	"context"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/AyumiLumi/BliveDataExportToExcel/export"
)

// CreateHomePage 构造 Home 页面
func CreateHomePage(myWindow fyne.Window, eventChans map[string]chan string, cancelChan chan struct{}) *fyne.Container {
	// 输入框
	roomIdEntry := widget.NewEntry()
	roomIdEntry.SetPlaceHolder("主播直播房间号")
	cookieEntry := widget.NewMultiLineEntry()
	cookieEntry.SetPlaceHolder("你的登录 Cookie")

	// 状态标签
	statusLabel := widget.NewLabel("准备中...")

	// 消息存储与显示
	messages := make([]string, 0, 10000)
	messageList := widget.NewList(
		func() int { return len(messages) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(messages[id])
		},
	)

	messageList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(messages) {
			msg := messages[id]
			clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
			clipboard.SetContent(msg)

			// 显示提示框
			showNotification(msg)
		}
	}

	// 自动滚动
	scrollContainer := container.NewVScroll(messageList)
	scrollContainer.SetMinSize(fyne.NewSize(400, 600))

	// 消息追加函数
	appendMessage := func(msg string) {
		if len(messages) >= 10000 {
			messages = messages[:len(messages)-1]
		}
		messages = append([]string{msg}, messages...)
		//go func() {
		//	//time.Sleep(60 * time.Millisecond)
		//	scrollContainer.ScrollToBottom()
		//}()
	}

	// 控制按钮
	var cancel context.CancelFunc
	startButton := widget.NewButton("Start", nil)

	// 设置 Start 按钮功能
	startButton.OnTapped = func() {
		roomId := roomIdEntry.Text
		cookie := cookieEntry.Text
		id, _ := strconv.Atoi(roomId)

		if roomId == "" || cookie == "" {
			statusLabel.SetText("请填写房间号和登录 Cookie,没有 Cookie 将无法完整显示信息！")
		}

		// 开始任务
		statusLabel.SetText("任务进行中...")
		startButton.Disable()

		// 创建上下文
		ctx, cancelFunc := context.WithCancel(context.Background())
		cancel = cancelFunc

		go func() {
			export.ExportExcel(ctx, id, cookie, eventChans, cancelChan)
		}()

		go func() {
			for {
				select {
				case msg := <-eventChans["home"]:
					appendMessage(msg)
					messageList.Refresh()
				case <-cancelChan:
					statusLabel.SetText("任务已完成")
					startButton.Enable()
					return
				}
			}
		}()
	}

	// Cancel 按钮
	cancelButton := widget.NewButton("Cancel", func() {
		if cancel != nil {
			cancel()
			statusLabel.SetText("取消中，保存文件...")
		}
	})

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

	// 版权声明
	copyrightLabel := widget.NewLabel("© 2024 By 滨崎步的水钻麦克风. All rights reserved.")
	copyrightContainer := container.NewCenter(copyrightLabel)

	// 按钮容器
	buttonsContainer := container.NewCenter(
		container.NewHBox(
			startButton,
			layout.NewSpacer(),
			cancelButton,
		),
	)

	// 布局
	content := container.NewVBox(
		widget.NewLabel("房间号"),
		roomIdEntry,
		widget.NewLabel("登录 Cookie"),
		cookieEntry,
		statusLabel,
		scrollContainer,
		buttonsContainer,
		copyrightContainer,
	)

	return content
}
