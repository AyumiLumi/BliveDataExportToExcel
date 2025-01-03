package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateEventPage creates a reusable layout for event pages.
func CreateEventPage(eventChan chan string, title string) *fyne.Container {
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

	// 监听事件
	go func() {
		for msg := range eventChan {
			appendMessage(msg)
			messageList.Refresh()
		}
	}()

	// 布局
	content := container.NewVBox(
		widget.NewLabel(title), // 页面标题
		scrollContainer,
	)

	return content
}
