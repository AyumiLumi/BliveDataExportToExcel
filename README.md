# BliveDataExportToExcel

## 环境准备

系统版本  Windows 11 24H2

Go version go1.23.0 windows/amd64

安装Fyne UI，快速上手
https://docs.fyne.io/started/
Fyne 需要 3 个基本元素，Go 工具（至少为 1.12 版）、C 编译器（用于连接系统图形驱动程序）和系统图形驱动程序。

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

### b站的弹幕消息数据样例

```
{
  cmd: "DANMU_MSG",
  dm_v2: "",
  info: [
    [
      0,
      1,
      25,
      5816798,
      1735119251294,
      1244493633,
      0,
      "d1b9df36",
      0,
      0,
      0,
      "",
      1,
      {
        bulge_display: 1,
        emoticon_unique: "room_23945626_3367",
        height: 162,
        in_player_area: 1,
        is_dynamic: 0,
        url: "http://i0.hdslb.com/bfs/live/953617cc7ff461f52c33a8daf4b35f956e8834e7.png",
        width: 162,
      },
      "{}",
      {
        extra:
          '{"send_from_me":false,"mode":0,"color":5816798,"dm_type":1,"font_size":25,"player_mode":1,"show_player_type":0,"content":"狂","user_hash":"3518619446","emoticon_unique":"room_23945626_3367","bulge_display":1,"recommend_score":0,"main_state_dm_color":"","objective_state_dm_color":"","direction":0,"pk_direction":0,"quartet_direction":0,"anniversary_crowd":0,"yeah_space_type":"","yeah_space_url":"","jump_to_url":"","space_type":"","space_url":"","animation":{},"emots":null,"is_audited":false,"id_str":"4a59e96b8375d8f36def9dab60676bd19334","icon":null,"show_reply":true,"reply_mid":0,"reply_uname":"","reply_uname_color":"","reply_is_mystery":false,"reply_type_enum":0,"hit_combo":0,"esports_jump_url":""}',
        mode: 0,
        show_player_type: 0,
        user: {
          base: {
            face: "https://i0.hdslb.com/bfs/face/ce02d62eb1fa99f625ab7c41ae5b66a3c9743430.webp",
            is_mystery: false,
            name: "滨崎步的水钻麦克风",
            name_color: 0,
            name_color_str: "",
            official_info: { desc: "", role: 0, title: "", type: -1 },
            origin_info: {
              face: "https://i0.hdslb.com/bfs/face/ce02d62eb1fa99f625ab7c41ae5b66a3c9743430.webp",
              name: "滨崎步的水钻麦克风",
            },
            risk_ctrl_info: null,
          },
          guard: null,
          guard_leader: { is_guard_leader: false },
          medal: {
            color: 12478086,
            color_border: 12478086,
            color_end: 12478086,
            color_start: 12478086,
            guard_icon: "",
            guard_level: 0,
            honor_icon: "",
            id: 541356,
            is_light: 1,
            level: 15,
            name: "黯灭头",
            ruid: 1855519979,
            score: 76610,
            typ: 0,
            user_receive_count: 0,
            v2_medal_color_border: "#C85DC499",
            v2_medal_color_end: "#C85DC499",
            v2_medal_color_level: "#59005699",
            v2_medal_color_start: "#C85DC499",
            v2_medal_color_text: "#FFFFFFFF",
          },
          title: { old_title_css_id: "", title_css_id: "" },
          uhead_frame: null,
          uid: 189194391,
          wealth: null,
        },
      },
      { activity_identity: "", activity_source: 0, not_show: 0 },
      0,
    ],
    "狂",
    [189194391, "滨崎步的水钻麦克风", 0, 0, 0, 10000, 1, ""],
    [
      15,
      "黯灭头",
      "dodo_Official",
      23945626,
      12478086,
      "",
      0,
      12478086,
      12478086,
      12478086,
      0,
      1,
      1855519979,
    ],
    [21, 0, 5805790, "\u003e50000", 0],
    ["", ""],
    0,
    0,
    null,
    { ct: "96156738", ts: 1735119251 },
    0,
    0,
    null,
    null,
    0,
    484,
    [34],
    null,
  ],
};

```

### B站大航海消息



```
{
  "cmd": "USER_TOAST_MSG_V2",
  "data": {
    "sender_uinfo": {
      "uid": 9123548,
      "base": {
        "name": "黑听的小耳朵",
        "face": ""
      }
    },
    "receiver_uinfo": {
      "uid": 922573,
      "base": {
        "name": "黑泽诺亚NOIR",
        "face": "https://i1.hdslb.com/bfs/face/bfb3f9c7cfdf347c9588c5f302cdd2bc096dd8fe.webp"
      }
    },
    "guard_info": {
      "guard_level": 3,
      "role_name": "舰长",
      "room_guard_count": 147,
      "op_type": 1,
      "start_time": 1735882838,
      "end_time": 1735882838
    },
    "group_guard_info": null,
    "pay_info": {
      "payflow_id": "2501031340270412435486336",
      "price": 138000,
      "num": 1,
      "unit": "月"
    },
    "gift_info": {
      "gift_id": 10003
    },
    "effect_info": {
      "effect_id": 397,
      "room_effect_id": 590,
      "face_effect_id": 44,
      "room_gift_effect_id": 0,
      "room_group_effect_id": 1337
    },
    "toast_msg": "\\u003c%黑听的小耳朵%\\u003e 在主播黑泽诺亚NOIR的直播间开通了舰长，今天是TA陪伴主播的第1天",
    "option": {
      "anchor_show": true,
      "user_show": true,
      "is_group": 0,
      "is_show": 0,
      "source": 0,
      "svga_block": 0,
      "color": "#00D1F1"
    }
  }
};
```

### B站红包数据格式

```
{
  cmd: "POPULARITY_RED_POCKET_V2_NEW",
  data: {
    lot_id: 25239368,
    start_time: 1735893924,
    current_time: 1735893924,
    wait_num: 0,
    wait_num_v2: 0,
    uname: "滨崎步的水钻麦克风",
    uid: 189194391,
    action: "送出",
    num: 1,
    gift_name: "红包",
    gift_id: 13000,
    price: 20,
    name_color: "#00D1F1",
    medal_info: {
      target_id: 3493126803032322,
      special: "",
      icon_id: 0,
      anchor_uname: "",
      anchor_roomid: 0,
      medal_level: 22,
      medal_name: "确实卡",
      medal_color: 1725515,
      medal_color_start: 1725515,
      medal_color_end: 5414290,
      medal_color_border: 6809855,
      is_lighted: 1,
      guard_level: 3,
    },
    wealth_level: 35,
    group_medal: null,
    is_mystery: false,
    sender_info: {
      uid: 189194391,
      base: {
        name: "滨崎步的水钻麦克风",
        face: "https://i0.hdslb.com/bfs/face/ce02d62eb1fa99f625ab7c41ae5b66a3c9743430.webp",
        name_color: 0,
        is_mystery: false,
        origin_info: {
          name: "滨崎步的水钻麦克风",
          face: "https://i0.hdslb.com/bfs/face/ce02d62eb1fa99f625ab7c41ae5b66a3c9743430.webp",
        },
        official_info: { role: 0, title: "", desc: "", type: -1 },
        name_color_str: "#00D1F1",
      },
      medal: {
        name: "确实卡",
        level: 22,
        color_start: 1725515,
        color_end: 5414290,
        color_border: 6809855,
        color: 1725515,
        id: 0,
        typ: 0,
        is_light: 1,
        ruid: 3493126803032322,
        guard_level: 3,
        score: 50002767,
        guard_icon:
          "https://i0.hdslb.com/bfs/live/143f5ec3003b4080d1b5f817a9efdca46d631945.png",
        honor_icon: "",
        v2_medal_color_start: "#43B3E3CC",
        v2_medal_color_end: "#43B3E3CC",
        v2_medal_color_border: "#5FC7F4FF",
        v2_medal_color_text: "#FFFFFFFF",
        v2_medal_color_level: "#00308C99",
        user_receive_count: 0,
      },
      wealth: { level: 35, dm_icon_key: "" },
      title: null,
      guard: { level: 3, expired_str: "2025-01-03 23:59:59" },
      uhead_frame: null,
      guard_leader: null,
    },
    gift_icon: "",
    rp_type: 0,
  },
};
```

