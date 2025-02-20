package export

import (
	"context"
	"fmt"
	"github.com/AyumiLumi/BliveDataExportToExcel/api"
	"github.com/AyumiLumi/BliveDataExportToExcel/client"
	"github.com/AyumiLumi/BliveDataExportToExcel/message"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
	"time"
)

var savePath string

func ExportExcel(ctx context.Context, roomId int, cookie string, eventChans map[string]chan string, cancelChan chan struct{}) {
	log.SetLevel(log.DebugLevel)
	c := client.NewClient(roomId) //25977291
	c.SetCookie(cookie)

	// 创建 Excel 文件
	file := xlsx.NewFile()
	sheet1, _ := file.AddSheet("礼物")
	sheet2, _ := file.AddSheet("SC")
	sheet3, _ := file.AddSheet("大航海")
	sheet4, _ := file.AddSheet("弹幕")
	sheet5, _ := file.AddSheet("红包")
	sheet6, _ := file.AddSheet("直播营收总计")

	// 向醒目留言工作表写入数据
	addHeaders(sheet2, []string{"Price", "Uname", "Uid", "Message"})

	// 向礼物工作表写入数据
	addHeaders(sheet1, []string{"Uname", "Uid", "GiftName", "Number", "gift.Num*gift.Price"})

	// 向大航海工作表写入数据
	addHeaders(sheet3, []string{"Uname", "Time", "Uid", "GuardLevel", "Number", "Price", "Message"})

	// 向红包工作表写入数据
	addHeaders(sheet5, []string{"Uname", "Time", "Uid", "Number", "Price", "Message"})

	var UpName string
	UpName = "lumi"

	superChatAllCount := 0.00
	giftAllCount := 0.00
	guardAllCount := 0.00
	redPocketAllCount := 0.00

	// 启动事件处理协程
	go func() {
		// 弹幕事件
		c.OnDanmaku(func(danmaku *message.Danmaku) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			row := sheet4.AddRow()
			cell := row.AddCell()
			// 格式化为可读的时间字符串
			formattedTime := format(danmaku.Timestamp)
			if len(danmaku.Sender.Medal.UpName) > 0 {
				UpName = danmaku.Sender.Medal.UpName
			}
			var content string
			if danmaku.Type == message.EmoticonDanmaku {
				content = fmt.Sprintf("%s [弹幕表情] <%d> %s %s Lv%d：%s 表情URL： %s  Uid:%d\n", formattedTime, danmaku.Sender.Level, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content, danmaku.Emoticon.Url, danmaku.Sender.Uid)
				cell.Value = fmt.Sprintf("%s Uid:%d [弹幕表情] <%d> %s %s Lv%d：%s 表情URL： %s\n", formattedTime, danmaku.Sender.Uid, danmaku.Sender.Level, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content, danmaku.Emoticon.Url)

				fmt.Printf("%s [弹幕表情] %s %s Lv%d：%s 表情URL： %s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content, danmaku.Emoticon.Url)
			} else {
				content = fmt.Sprintf("%s [弹幕] <%d> %s %s Lv%d：%s Uid:%d\n", formattedTime, danmaku.Sender.Level, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content, danmaku.Sender.Uid)
				cell.Value = fmt.Sprintf("%s Uid:%d [弹幕] <%d> %s %s Lv%d：%s\n", formattedTime, danmaku.Sender.Uid, danmaku.Sender.Level, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content)

				fmt.Printf("%s [弹幕] %s %s Lv%d：%s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content)
			}
			eventChans["danmaku"] <- content // 发送到 danmaku
			eventChans["home"] <- content    // 发送到 home
		})

		// SC事件
		c.OnSuperChat(func(superChat *message.SuperChat) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			formattedTime := "  " + format(superChat.StartTime)
			MadelInfo := " " + superChat.MedalInfo.MedalName + " Lv" + strconv.Itoa(superChat.MedalInfo.MedalLevel)
			content := fmt.Sprintf("[SC|%d元] %s %s: %s，%s", superChat.Price, superChat.UserInfo.Uname, MadelInfo, superChat.Message, formattedTime)
			eventChans["home"] <- content      // 发送到 UI
			eventChans["superchat"] <- content // 发送到 UI
			row := sheet2.AddRow()
			row.AddCell().Value = strconv.Itoa(superChat.Price)
			row.AddCell().Value = superChat.UserInfo.Uname + MadelInfo
			row.AddCell().Value = strconv.Itoa(superChat.Uid)
			row.AddCell().Value = superChat.Message + formattedTime
			superChatAllCount += float64(superChat.Price)
		})

		// 礼物事件
		c.OnGift(func(gift *message.Gift) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if gift.CoinType == "gold" {
				MedalInfo := " " + gift.MedalInfo.MedalName + " Lv" + strconv.Itoa(gift.MedalInfo.MedalLevel)
				content := fmt.Sprintf("[礼物] %s %s 的 %s x%d，共 %.2f 元，%s", gift.Uname, MedalInfo, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000, format(gift.Timestamp))
				eventChans["gift"] <- content // 发送到 UI
				eventChans["home"] <- content // 发送到 UI
				row := sheet1.AddRow()
				row.AddCell().Value = gift.Uname + MedalInfo
				row.AddCell().Value = strconv.Itoa(gift.Uid)
				row.AddCell().Value = gift.GiftName + "  " + format(gift.Timestamp)
				row.AddCell().Value = strconv.Itoa(gift.Num)
				row.AddCell().Value = fmt.Sprintf("%.2f", float64(gift.Num*gift.Price)/1000)
				giftAllCount += float64(gift.Num*gift.Price) / 1000
			}
		})

		c.OnGuard(func(guard *message.Guard) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			content := fmt.Sprintf("[大航海] %s  %s", guard.ToastMsg, format(guard.GuardInfo.StartTime))
			eventChans["home"] <- content  // 发送到 UI
			eventChans["guard"] <- content // 发送到 UI
			row := sheet3.AddRow()
			row.AddCell().Value = guard.SenderUinfo.Base.Name
			row.AddCell().Value = format(guard.GuardInfo.StartTime)
			row.AddCell().Value = strconv.Itoa(guard.SenderUinfo.Uid)
			row.AddCell().Value = guard.GuardInfo.RoleName
			row.AddCell().Value = strconv.Itoa(guard.PayInfo.Num) + "个" + guard.PayInfo.Unit
			row.AddCell().Value = fmt.Sprintf("%.2f", float64(guard.PayInfo.Price)/1000)
			row.AddCell().Value = fmt.Sprintf("%s", guard.ToastMsg)
			guardAllCount += float64(guard.PayInfo.Price / 1000)
		})

		//红包事件
		c.OnRedPocket(func(redPocket *message.RedPocket) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			content := fmt.Sprintf("[红包] %s  [Uid:%d]  %s  Lv%d  %s%s， 共%d个，%.2f元 %s", redPocket.Uname, redPocket.Uid, redPocket.SenderInfo.Medal.Name, redPocket.SenderInfo.Medal.Level,
				redPocket.Action, redPocket.GiftName, redPocket.Num, float64(redPocket.Price/10), format(redPocket.StartTime))
			eventChans["home"] <- content      // 发送到 UI
			eventChans["redpocket"] <- content // 发送到 UI
			row := sheet5.AddRow()
			row.AddCell().Value = redPocket.Uname
			row.AddCell().Value = format(redPocket.StartTime)
			row.AddCell().Value = strconv.Itoa(redPocket.Uid)
			//row.AddCell().Value = guard.GuardInfo.RoleName
			row.AddCell().Value = strconv.Itoa(redPocket.Num) + "个"
			row.AddCell().Value = fmt.Sprintf("%.2f", float64(redPocket.Price)/10)
			row.AddCell().Value = fmt.Sprintf("%s", content)
			redPocketAllCount += float64(redPocket.Price / 10)
		})

		c.Start()
	}()

	// 保存数据
	go func() {
		select {
		case <-ctx.Done():
			// 汇总数据
			addHeaders(sheet6, []string{"数据仅仅是礼物原价信息，由于API没有折扣信息，仅供参考！"})
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("礼物总计：%.2f元", giftAllCount)
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("SC总计：%.2f元", superChatAllCount)
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("红包总计：%.2f元", redPocketAllCount)
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("大航海总计：%.2f元", guardAllCount)
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("总计：%.2f元", giftAllCount+superChatAllCount+guardAllCount+redPocketAllCount)
			saveExcelFile(file, UpName)
			eventChans["home"] <- "任务完成，文件保存成功！路径：" + savePath
			cancelChan <- struct{}{}
		}
	}()
}

// 定义一个通用函数用于写入表头
func addHeaders(sheet *xlsx.Sheet, headers []string) {
	row := sheet.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}
}

func saveExcelFile(file *xlsx.File, upName string) {
	currentTime := time.Now()
	fileName := currentTime.Format("2006-01-02_15-04-05") + "_" + upName + ".xlsx"
	savePath = "D:\\excel\\" + fileName

	fmt.Printf("准备保存文件到: %s\n", savePath)

	// 确保目录存在
	if err := os.MkdirAll("D:\\excel", os.ModePerm); err != nil {
		fmt.Printf("创建目录失败: %s\n", err)
		return
	}

	// 保存文件
	err := file.Save(savePath)
	if err != nil {
		fmt.Printf("保存 Excel 文件失败: %s\n", err)
	} else {
		fmt.Printf("Excel 文件已保存到: %s\n", savePath)
	}
}

func sendDanmaku() error {
	dmReq := &api.DanmakuRequest{
		Msg:      "official_13",
		RoomID:   "732",
		Bubble:   "0",
		Color:    "16777215",
		FontSize: "25",
		Mode:     "1",
		DmType:   "1",
	}
	d, err := api.SendDanmaku(dmReq, &api.BiliVerify{
		Csrf:     "",
		SessData: "",
	})
	if err != nil {
		return err
	}
	fmt.Println(d)
	return nil
}

func format(timestamp int64) string {
	// 如果时间戳长度大于13位，可能是毫秒
	if timestamp >= 1e12 && timestamp < 1e15 {
		timestamp /= 1000
	}

	// 转换为 Time 类型
	t := time.Unix(timestamp, 0)

	// 使用本地时间（系统默认时区）
	return t.Local().Format("2006-01-02 15:04:05")
}

//import (
//"github.com/evansb/go-timezone"
//)
//
//func format(timestamp int64) string {
//	if timestamp > 1e12 {
//		timestamp /= 1000 // 转换为秒
//	}
//
//	// 使用 timezonedb 来加载时区
//	location, err := timezone.LoadLocation("Asia/Shanghai")
//	if err != nil {
//		fmt.Println("加载时区失败:", err)
//		return err.Error()
//	}
//
//	// 将时间戳转换为 Time 类型
//	t := time.Unix(timestamp, 0).In(location)
//
//	// 格式化为可读的时间字符串
//	return t.Format("2006-01-02 15:04:05")
//}
