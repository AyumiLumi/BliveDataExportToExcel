package message

import (
	"github.com/AyumiLumi/BliveDataExportToExcel/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Guard struct {
	SenderUinfo    SenderUinfo   `json:"sender_uinfo"`
	ReceiverUinfo  ReceiverUinfo `json:"receiver_uinfo"`
	GuardInfo      GuardInfo     `json:"guard_info"`
	GroupGuardInfo interface{}   `json:"group_guard_info"`
	PayInfo        PayInfo       `json:"pay_info"`
	GiftInfo       GiftInfo      `json:"gift_info"`
	EffectInfo     EffectInfo    `json:"effect_info"`
	ToastMsg       string        `json:"toast_msg"`
	Option         Option        `json:"option"`
}

// 其他嵌套的结构体定义
type SenderUinfo struct {
	Uid  int  `json:"uid"`
	Base Base `json:"base"`
}

type Base struct {
	Name string `json:"name"`
	Face string `json:"face"`
}

type ReceiverUinfo struct {
	Uid  int  `json:"uid"`
	Base Base `json:"base"`
}

type GuardInfo struct {
	GuardLevel     int    `json:"guard_level"`
	RoleName       string `json:"role_name"`
	RoomGuardCount int    `json:"room_guard_count"`
	OpType         int    `json:"op_type"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
}

type PayInfo struct {
	PayflowID string `json:"payflow_id"`
	Price     int    `json:"price"`
	Num       int    `json:"num"`
	Unit      string `json:"unit"`
}

type GiftInfo struct {
	GiftID int `json:"gift_id"`
}

type EffectInfo struct {
	EffectID          int `json:"effect_id"`
	RoomEffectID      int `json:"room_effect_id"`
	FaceEffectID      int `json:"face_effect_id"`
	RoomGiftEffectID  int `json:"room_gift_effect_id"`
	RoomGroupEffectID int `json:"room_group_effect_ID"`
}

type Option struct {
	AnchorShow bool   `json:"anchor_show"`
	UserShow   bool   `json:"user_show"`
	IsGroup    bool   `json:"is_group"`
	IsShow     bool   `json:"is_show"`
	Source     int    `json:"source"`
	SvgaBlock  bool   `json:"svga_block"`
	Color      string `json:"color"`
}

func (g *Guard) Parse(data []byte) {
	sb := utils.BytesToString(data)
	sd := gjson.Get(sb, "data").String()
	err := utils.UnmarshalStr(sd, g)
	if err != nil {
		log.Error("parse Guard failed")
	}
}
