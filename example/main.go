package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/AyumiLumi/BliveDataExportToExcel/api"
	"github.com/AyumiLumi/BliveDataExportToExcel/client"
	"github.com/AyumiLumi/BliveDataExportToExcel/message"
	_ "github.com/AyumiLumi/BliveDataExportToExcel/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	c := client.NewClient(21738566) //25977291
	c.SetCookie("buvid3=6A93F404-5CBE-5689-ABAF-C86651B9795833943infoc; b_nut=1713776833; _uuid=BC42ED39-4AA6-A562-F4FB-68EAB9158ADD34983infoc; buvid4=2F3D7726-8DE1-ECFC-37A8-37E747BBE57A34742-024042209-j9RQN9%2FKfsrv2a7%2BSKiPKg%3D%3D; CURRENT_FNVAL=4048; rpdid=|(J|~uklm|mm0J'u~uRukJkRm; DedeUserID=189194391; DedeUserID__ckMd5=f9f76ed12099abcc; buvid_fp_plain=undefined; LIVE_BUVID=AUTO7817165511075577; CURRENT_QUALITY=120; bsource=search_google; hit-dyn-v2=1; header_theme_version=CLOSE; enable_web_push=DISABLE; home_feed_column=5; browser_resolution=1707-932; bp_t_offset_189194391=997367083454955520; PVID=1; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzIxNjk2ODAsImlhdCI6MTczMTkxMDQyMCwicGx0IjotMX0.7rvQlsKMgz8Nu6dSivowWqDJAO-jZq0_NrkBp0qd3K4; bili_ticket_expires=1732169620; SESSDATA=5e187120%2C1747474462%2C0f91c%2Ab2CjAZ4SDajWGzZkIw7z0cfFt-TMcLwrdHNuOakrXLg_HfeZodN_Y0CT8oE5dZGMIR7NgSVjU4YnJkcEhZSGItRWNpcGV0UWFZdlZvVUhKMThjb0cyZ0EzLVV2azR0RUNLYk5rQmsta1RxQVh1aTlSZXkyMnNGY3pxT0lZMVNfVDJnMHRLbmtCODZ3IIEC; bili_jct=ef850374e44502973f6f59f0e3fb0688; sid=5catmsa3; b_lsid=984EC8D10_1934769A734; fingerprint=5c4287d93140d832c73b98439445e0b7; buvid_fp=5c4287d93140d832c73b98439445e0b7")
	//guardBuyChan := make(chan struct{})

	// 创建一个新的 Excel 文件
	file := xlsx.NewFile()

	// 创建第一个工作表
	sheet1, err1 := file.AddSheet("礼物")
	if err1 != nil {
		fmt.Printf(err1.Error())
		//return
	}

	sheet2, err2 := file.AddSheet("SC")
	if err2 != nil {
		fmt.Printf(err2.Error())
		//return
	}

	sheet3, err3 := file.AddSheet("大航海")
	if err3 != nil {
		fmt.Printf(err3.Error())
		//return
	}

	sheet4, err4 := file.AddSheet("弹幕")
	if err4 != nil {
		fmt.Printf(err4.Error())
		//return
	}

	sheet6, err6 := file.AddSheet("直播营收总计")
	if err4 != nil {
		fmt.Printf(err6.Error())
		//return
	}

	var UpName string
	UpName = "lumi"

	// 向醒目留言工作表写入数据
	row0_super_chat := sheet2.AddRow()
	cell0_super_chat := row0_super_chat.AddCell()
	cell0_super_chat.Value = "Price"
	cell1_super_chat := row0_super_chat.AddCell()
	cell1_super_chat.Value = "Uname"
	cell3_super_chat := row0_super_chat.AddCell()
	cell3_super_chat.Value = "Uid"
	cell2_super_chat := row0_super_chat.AddCell()
	cell2_super_chat.Value = "Message"

	// 向礼物工作表写入数据
	row0_gift := sheet1.AddRow()
	cell0_gift := row0_gift.AddCell()
	cell0_gift.Value = "Uname"
	cell4_gift := row0_gift.AddCell()
	cell4_gift.Value = "Uid"
	cell1_gift := row0_gift.AddCell()
	cell1_gift.Value = "GiftName"
	cell3_gift := row0_gift.AddCell()
	cell3_gift.Value = "Number"
	cell2_gift := row0_gift.AddCell()
	cell2_gift.Value = "gift.Num*gift.Price"

	// 向大航海工作表写入数据
	row0_guard := sheet3.AddRow()
	cell0_guard := row0_guard.AddCell()
	cell0_guard.Value = "Uname"
	cell1_guard := row0_guard.AddCell()
	cell1_guard.Value = "Uid"
	cell3_guard := row0_guard.AddCell()
	cell3_guard.Value = "GuardLevel"
	cell4_guard := row0_guard.AddCell()
	cell4_guard.Value = "Number"
	cell2_guard := row0_guard.AddCell()
	cell2_guard.Value = "Price"

	superChatAllCount := 0.00
	giftAllCount := 0.00
	guardAllCount := 0.00

	go func() {
		//弹幕事件
		c.OnDanmaku(func(danmaku *message.Danmaku) {
			row := sheet4.AddRow()
			cell := row.AddCell()
			// 格式化为可读的时间字符串
			formattedTime := format(danmaku.Timestamp)
			if len(danmaku.Sender.Medal.UpName) > 0 {
				UpName = danmaku.Sender.Medal.UpName
			}
			if danmaku.Type == message.EmoticonDanmaku {
				cell.Value = fmt.Sprintf("%s [弹幕表情] %s %s Lv%d：%s 表情URL： %s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content, danmaku.Emoticon.Url)
				fmt.Printf("%s [弹幕表情] %s %s Lv%d：%s 表情URL： %s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content, danmaku.Emoticon.Url)
			} else {
				cell.Value = fmt.Sprintf("%s [弹幕] %s %s Lv%d：%s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content)
				fmt.Printf("%s [弹幕] %s %s Lv%d：%s\n", formattedTime, danmaku.Sender.Uname, danmaku.Sender.Medal.Name, danmaku.Sender.Medal.Level, danmaku.Content)
			}
		})

		// 醒目留言事件
		c.OnSuperChat(func(superChat *message.SuperChat) {
			fmt.Printf("[SC|%d元] %s: %s\n", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
			// 向醒目留言工作表写入数据
			row1_super_chat := sheet2.AddRow()
			row1_cell0_super_chat := row1_super_chat.AddCell()
			row1_cell0_super_chat.Value = strconv.Itoa(superChat.Price)
			row1_cell1_super_chat := row1_super_chat.AddCell()
			row1_cell1_super_chat.Value = superChat.UserInfo.Uname
			row1_cell3_super_chat := row1_super_chat.AddCell()
			row1_cell3_super_chat.Value = strconv.Itoa(superChat.Uid)
			row1_cell2_super_chat := row1_super_chat.AddCell()
			row1_cell2_super_chat.Value = superChat.Message + "  " + format(int64(superChat.StartTime))
			superChatAllCount += float64(superChat.Price)
		})

		// 礼物事件
		c.OnGift(func(gift *message.Gift) {
			if gift.CoinType == "gold" {
				fmt.Printf("[礼物] %s 的 %s %d 个 共%.2f元\n", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
				// 向礼物工作表写入数据
				row1_gift := sheet1.AddRow()
				row1_cell0_gift := row1_gift.AddCell()
				row1_cell0_gift.Value = gift.Uname
				row1_cell4_gift := row1_gift.AddCell()
				row1_cell4_gift.Value = strconv.Itoa(gift.Uid)
				row1_cell1_gift := row1_gift.AddCell()
				row1_cell1_gift.Value = gift.GiftName + "  " + format(int64(gift.Timestamp))
				row1_cell3_gift := row1_gift.AddCell()
				row1_cell3_gift.Value = strconv.Itoa(gift.Num)
				row1_cell2_gift := row1_gift.AddCell()
				row1_cell2_gift.Value = strconv.FormatFloat(float64(gift.Num*gift.Price)/1000, 'f', 2, 64)
				giftAllCount += float64(gift.Num*gift.Price) / 1000
			}
		})

		// 上舰事件
		c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
			fmt.Printf("[大航海] %s 开通了 %d 等级的大航海 * %d，金额 %d 元\n", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Num, float64(guardBuy.Price*guardBuy.Num)/1000)
			//close(guardBuyChan)
			// 向大航海工作表写入数据
			row1_guard := sheet3.AddRow()
			row1_cell0_guard := row1_guard.AddCell()
			row1_cell0_guard.Value = guardBuy.Username + "  " + format(int64(guardBuy.StartTime))
			row1_cell1_guard := row1_guard.AddCell()
			row1_cell1_guard.Value = strconv.Itoa(guardBuy.Uid)
			row1_cell3_guard := row1_guard.AddCell()
			row1_cell3_guard.Value = strconv.Itoa(guardBuy.GuardLevel)
			row1_cell4_guard := row1_guard.AddCell()
			row1_cell4_guard.Value = strconv.Itoa(guardBuy.Num)
			row1_cell2_guard := row1_guard.AddCell()
			row1_cell2_guard.Value = strconv.FormatFloat(float64(guardBuy.Price*guardBuy.Num)/1000, 'f', 2, 64)
			guardAllCount += float64(guardBuy.Price / 1000)
		})

		// 监听自定义事件
		c.RegisterCustomEventHandler("STOP_LIVE_ROOM_LIST", func(s string) {
			data := gjson.Get(s, "data").String()
			fmt.Printf("STOP_LIVE_ROOM_LIST: %s\n", data)
		})
		err := c.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("started")
	// 需要自行阻塞什么方法都可以

	// 捕获退出信号
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, os.Kill)

	// 等待退出信号
	//<-stopChan
	//fmt.Println("接收到退出信号，正在保存文件...")

	// 使用 WaitGroup 确保保存文件完成后退出
	var wg sync.WaitGroup
	wg.Add(1)
	// 等待退出信号或手动停止
	select {
	case <-stopChan: // 退出信号到来
		fmt.Println("接收到退出信号，正在保存文件...")
		// 在退出时保存文件
		go func() {
			defer wg.Done()
			giftRow := sheet6.AddRow()
			giftCell := giftRow.AddCell()
			giftCell.Value = "礼物总计：     " + strconv.FormatFloat(giftAllCount, 'f', 2, 64) + "元"
			superChatRow := sheet6.AddRow()
			superChatCell := superChatRow.AddCell()
			superChatCell.Value = "SC总计：      " + strconv.FormatFloat(superChatAllCount, 'f', 2, 64) + "元"
			guardRow := sheet6.AddRow()
			guardCell := guardRow.AddCell()
			guardCell.Value = "大航海总计：     " + strconv.FormatFloat(guardAllCount, 'f', 2, 64) + "元"
			totalRow := sheet6.AddRow()
			totalCell := totalRow.AddCell()
			totalCell.Value = "总计：     " + strconv.FormatFloat((guardAllCount+giftAllCount+superChatAllCount), 'f', 2, 64) + "元"
			saveExcelFile(file, UpName)
		}()
	}

	// 等待文件保存完成
	wg.Wait()
	fmt.Println("文件保存完成，程序退出")
}

func saveExcelFile(file *xlsx.File, upName string) {
	currentTime := time.Now()
	fileName := currentTime.Format("2006-01-02_15-04-05") + "_" + upName + ".xlsx"
	savePath := "D:\\excel\\" + fileName

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
	// 加载时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("加载时区失败:", err)
		return err.Error()
	}
	// 将时间戳转换为 Time 类型
	t := time.Unix(timestamp, 0).In(location)

	// 格式化为可读的时间字符串
	return t.Format("2006-01-02 15:04:05")
}
