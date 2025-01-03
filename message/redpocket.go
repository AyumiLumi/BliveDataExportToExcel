package message

import (
	"github.com/AyumiLumi/BliveDataExportToExcel/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type RedPocket struct {
	LotID       int    `json:"lot_id"`
	StartTime   int64  `json:"start_time"`
	CurrentTime int64  `json:"current_time"`
	WaitNum     int    `json:"wait_num"`
	WaitNumV2   int    `json:"wait_num_v2"`
	Uname       string `json:"uname"`
	Uid         int    `json:"uid"`
	Action      string `json:"action"`
	Num         int    `json:"num"`
	GiftName    string `json:"gift_name"`
	GiftID      int    `json:"gift_id"`
	Price       int    `json:"price"`
	NameColor   string `json:"name_color"`
	MedalInfo   struct {
		TargetID         int64  `json:"target_id"`
		Special          string `json:"special"`
		IconID           int    `json:"icon_id"`
		AnchorUname      string `json:"anchor_uname"`
		AnchorRoomID     int    `json:"anchor_roomid"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		MedalColor       int    `json:"medal_color"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorBorder int    `json:"medal_color_border"`
		IsLighted        int    `json:"is_lighted"`
		GuardLevel       int    `json:"guard_level"`
	} `json:"medal_info"`
	WealthLevel int         `json:"wealth_level"`
	GroupMedal  interface{} `json:"group_medal"`
	IsMystery   bool        `json:"is_mystery"`
	SenderInfo  struct {
		Uid  int `json:"uid"`
		Base struct {
			Name       string `json:"name"`
			Face       string `json:"face"`
			NameColor  int    `json:"name_color"`
			IsMystery  bool   `json:"is_mystery"`
			OriginInfo struct {
				Name string `json:"name"`
				Face string `json:"face"`
			} `json:"origin_info"`
			OfficialInfo struct {
				Role  int    `json:"role"`
				Title string `json:"title"`
				Desc  string `json:"desc"`
				Type  int    `json:"type"`
			} `json:"official_info"`
			NameColorStr string `json:"name_color_str"`
		} `json:"base"`
		Medal struct {
			Name               string `json:"name"`
			Level              int    `json:"level"`
			ColorStart         int    `json:"color_start"`
			ColorEnd           int    `json:"color_end"`
			ColorBorder        int    `json:"color_border"`
			Color              int    `json:"color"`
			ID                 int    `json:"id"`
			Typ                int    `json:"typ"`
			IsLight            int    `json:"is_light"`
			Ruid               int64  `json:"ruid"`
			GuardLevel         int    `json:"guard_level"`
			Score              int    `json:"score"`
			GuardIcon          string `json:"guard_icon"`
			HonorIcon          string `json:"honor_icon"`
			V2MedalColorStart  string `json:"v2_medal_color_start"`
			V2MedalColorEnd    string `json:"v2_medal_color_end"`
			V2MedalColorBorder string `json:"v2_medal_color_border"`
			V2MedalColorText   string `json:"v2_medal_color_text"`
			V2MedalColorLevel  string `json:"v2_medal_color_level"`
			UserReceiveCount   int    `json:"user_receive_count"`
		} `json:"medal"`
		Wealth struct {
			Level     int    `json:"level"`
			DmIconKey string `json:"dm_icon_key"`
		} `json:"wealth"`
		Title interface{} `json:"title"`
		Guard struct {
			Level      int    `json:"level"`
			ExpiredStr string `json:"expired_str"`
		} `json:"guard"`
		UheadFrame  interface{} `json:"uhead_frame"`
		GuardLeader interface{} `json:"guard_leader"`
	} `json:"sender_info"`
	GiftIcon string `json:"gift_icon"`
	RpType   int    `json:"rp_type"`
}

func (g *RedPocket) Parse(data []byte) {
	sb := utils.BytesToString(data)
	sd := gjson.Get(sb, "data").String()
	err := utils.UnmarshalStr(sd, g)
	if err != nil {
		log.Error("parse RedPocket failed")
	}
}
