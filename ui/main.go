package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/AyumiLumi/BliveDataExportToExcel/ui/pages"
)

func main() {
	myApp := app.NewWithID("看看你的")
	myWindow := myApp.NewWindow("看看你的")

	// 定义事件 channels
	eventChans := map[string]chan string{
		"home":      make(chan string, 300),
		"danmaku":   make(chan string, 150),
		"gift":      make(chan string, 200),
		"superchat": make(chan string, 300),
		"redpocket": make(chan string, 300),
		"guard":     make(chan string, 100),
	}

	cancelChan := make(chan struct{})

	// 创建 Home 页面
	homePage := pages.CreateHomePage(myWindow, eventChans, cancelChan)

	// 其他事件页面
	danmakuPage := pages.CreateEventPage(eventChans["danmaku"], "弹幕事件")
	giftPage := pages.CreateEventPage(eventChans["gift"], "礼物事件")
	superChatPage := pages.CreateEventPage(eventChans["superchat"], "SC 事件")
	guardPage := pages.CreateEventPage(eventChans["guard"], "大航海事件")
	redPocketPage := pages.CreateEventPage(eventChans["redpocket"], "红包事件")

	// 使用 Tabs 创建侧边栏
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), homePage),
		container.NewTabItem("弹幕", danmakuPage),
		container.NewTabItem("礼物", giftPage),
		container.NewTabItem("SC", superChatPage),
		container.NewTabItem("红包", redPocketPage),
		container.NewTabItem("大航海", guardPage),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	// 设置窗口内容
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
