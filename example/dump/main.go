package main

import (
	"fmt"
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	_ "github.com/Akegarasu/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)

	c := client.NewClient(25977291)
	//guardBuyChan := make(chan struct{})
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, os.Kill)

	// Excel 文件创建
	file := xlsx.NewFile()
	_, _ = file.AddSheet("礼物")
	_, _ = file.AddSheet("SC")
	_, _ = file.AddSheet("大航海")
	_, _ = file.AddSheet("弹幕")

	upName := "lumi"

	// 事件注册
	c.OnDanmaku(func(danmaku *message.Danmaku) {
		fmt.Println("[弹幕] 收到弹幕消息...")
	})
	err := c.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("started")

	// 捕获退出信号并同步保存文件
	<-stopChan
	fmt.Println("接收到退出信号，正在保存文件...")
	saveExcelFile(file, upName)
	fmt.Println("文件保存完成，程序退出")
}
func saveExcelFile(file *xlsx.File, upName string) {
	currentTime := time.Now()
	fileName := currentTime.Format("20060102_150405") + "_" + upName + ".xlsx"
	savePath := "D:\\excel\\" + fileName

	fmt.Printf("开始保存文件到: %s\n", savePath)

	// 确保目录存在
	if err := os.MkdirAll("D:\\excel", os.ModePerm); err != nil {
		fmt.Printf("创建目录失败: %s\n", err)
		return
	}

	// 保存文件
	start := time.Now()
	err := file.Save(savePath)
	if err != nil {
		fmt.Printf("保存 Excel 文件失败: %s\n", err)
	} else {
		fmt.Printf("Excel 文件已成功保存到: %s\n", savePath)
	}
	fmt.Printf("文件保存耗时: %v\n", time.Since(start))
}
