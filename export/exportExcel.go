package export

import (
	"context"
	"fmt"
	"github.com/Akegarasu/blivedm-go/api"
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
	"time"
)

var savePath string

func ExportExcel(ctx context.Context, roomId int, cookie string, eventChan chan string, cancelChan chan struct{}) {
	log.SetLevel(log.DebugLevel)
	c := client.NewClient(roomId) //25977291
	c.SetCookie(cookie)

	// 创建 Excel 文件
	file := xlsx.NewFile()
	sheet1, _ := file.AddSheet("礼物")
	sheet2, _ := file.AddSheet("SC")
	sheet3, _ := file.AddSheet("大航海")
	sheet4, _ := file.AddSheet("弹幕")
	sheet6, _ := file.AddSheet("直播营收总计")

	// 向醒目留言工作表写入数据
	addHeaders(sheet2, []string{"Price", "Uname", "Uid", "Message"})

	// 向礼物工作表写入数据
	addHeaders(sheet1, []string{"Uname", "Uid", "GiftName", "Number", "gift.Num*gift.Price"})

	// 向大航海工作表写入数据
	addHeaders(sheet3, []string{"Uname", "FansMedalName&Level", "Time", "Uid", "GuardLevel", "Number", "Price"})

	var UpName string
	UpName = "lumi"

	superChatAllCount := 0.00
	giftAllCount := 0.00
	guardAllCount := 0.00

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
				content = fmt.Sprintf("%s [弹幕表情] %s %s Lv%d：%s 表情URL： %s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content, danmaku.Emoticon.Url)
				cell.Value = content
				fmt.Printf("%s [弹幕表情] %s %s Lv%d：%s 表情URL： %s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content, danmaku.Emoticon.Url)
			} else {
				content = fmt.Sprintf("%s [弹幕] %s %s Lv%d：%s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content)
				cell.Value = content
				fmt.Printf("%s [弹幕] %s %s Lv%d：%s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content)
			}
			eventChan <- content // 发送到 UI
		})

		// SC事件
		c.OnSuperChat(func(superChat *message.SuperChat) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			content := fmt.Sprintf("[SC|%d元] %s: %s", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
			eventChan <- content // 发送到 UI
			row := sheet2.AddRow()
			row.AddCell().Value = strconv.Itoa(superChat.Price)
			row.AddCell().Value = superChat.UserInfo.Uname + superChat.FansMedalName + " Lv" + strconv.Itoa(superChat.FansMedalLevel)
			row.AddCell().Value = strconv.Itoa(superChat.Uid)
			row.AddCell().Value = superChat.Message + "  " + format(superChat.Timestamp)
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
				content := fmt.Sprintf("[礼物] %s 的 %s x%d，共 %.2f 元", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
				eventChan <- content // 发送到 UI
				row := sheet1.AddRow()
				row.AddCell().Value = gift.Uname + gift.FansMedalName + " Lv" + strconv.Itoa(gift.FansMedalLevel)
				row.AddCell().Value = strconv.Itoa(gift.Uid)
				row.AddCell().Value = gift.GiftName + "  " + format(gift.Timestamp)
				row.AddCell().Value = strconv.Itoa(gift.Num)
				row.AddCell().Value = fmt.Sprintf("%.2f", float64(gift.Num*gift.Price)/1000)
				giftAllCount += float64(gift.Num*gift.Price) / 1000
			}
		})

		// 大航海事件
		c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			content := fmt.Sprintf("[大航海] %s 开通了 %d 等级的大航海，共 %d 个，金额 %.2f 元", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Num, float64(guardBuy.Price*guardBuy.Num)/1000)
			eventChan <- content // 发送到 UI
			row := sheet3.AddRow()
			row.AddCell().Value = guardBuy.Username
			row.AddCell().Value = guardBuy.FansMedalName + " Lv" + strconv.Itoa(guardBuy.FansMedalLevel)
			row.AddCell().Value = format(guardBuy.Timestamp)
			row.AddCell().Value = strconv.Itoa(guardBuy.Uid)
			row.AddCell().Value = strconv.Itoa(guardBuy.GuardLevel)
			row.AddCell().Value = strconv.Itoa(guardBuy.Num) + guardBuy.GuardUnit
			row.AddCell().Value = fmt.Sprintf("%.2f", float64(guardBuy.Price*guardBuy.Num)/1000)
			guardAllCount += float64(guardBuy.Price * guardBuy.Num / 1000)
		})

		c.Start()
	}()

	// 保存数据
	go func() {
		select {
		case <-ctx.Done():
			// 汇总数据
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("礼物总计：%.2f元", giftAllCount)
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("SC总计：%.2f元", superChatAllCount)
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("大航海总计：%.2f元", guardAllCount)
			sheet6.AddRow().AddCell().Value = fmt.Sprintf("总计：%.2f元", giftAllCount+superChatAllCount+guardAllCount)
			saveExcelFile(file, UpName)
			eventChan <- "任务完成，文件保存成功！路径：" + savePath
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
	if timestamp > 1e12 {
		timestamp /= 1000 // 转换为秒
	}

	// 将时间戳转换为 Time 类型
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