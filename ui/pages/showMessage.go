package pages

import (
	"fyne.io/fyne/v2"
)

func showNotification(message string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "复制成功",
		Content: message,
	})
}
