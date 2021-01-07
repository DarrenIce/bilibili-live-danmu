package danmu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/issue9/term/v2/colors"
)


func (d *DanmuClient) process() {
	for {
		select {
		case m := <- d.unzlibChannel:
			uz := m[16:]
			js := new(receivedInfo)
			json.Unmarshal(uz, js)
			switch js.Cmd {
			case "ACTIVITY_BANNER_UPDATE_V2":
				continue
			case "COMBO_SEND":
				gs := colors.New(colors.Normal, 203, colors.Black)
				gs.Printf("[%s] %s: [%s] 送给 [%s] %d 个 %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.Data["uname"].(string), js.Data["r_uname"].(string), int(js.Data["combo_num"].(float64)), js.Data["gift_name"].(string))
			case "DANMU_MSG":
				dm := colors.New(colors.Normal, 158, colors.Black)
				dm.Printf("[%s] %s: [%s] %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.Info[2].([]interface{})[1].(string), js.Info[1].(string))
			case "ENTRY_EFFECT":
				ee := colors.New(colors.Normal, 75, colors.Black)
				ee.Printf("[%s] %s: %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.Data["copy_writing_v2"].(string))
			case "GUARD_BUY":
				fmt.Printf("[%s] %s: [%s]购买了[%s]. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.Data["username"].(string), js.Data["gift_name"].(string))
			case "INTERACT_WORD":
				iw := colors.New(colors.Normal, 29, colors.Black)
				iw.Printf("[%s] %s: [%s]进入了房间[%d]. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.Data["uname"].(string), uint32(js.Data["roomid"].(float64)))
			case "LIVE_INTERACTIVE_GAME":
				continue
			case "LIVE":
				continue
			case "NOTICE_MSG":
				fmt.Printf("[%s] %s: %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.MsgSelf)
			case "ONLINE_RANK_COUNT":
				continue
			case "ONLINE_RANK_TOP3":
				continue
			case "ONLINE_RANK_V2":
				continue
			case "PANEL":
				continue
			case "PREPARING":
				continue
			case "SEND_GIFT":
				gs := colors.New(colors.Normal, 203, colors.Black)
				gs.Printf("[%s] %s: [%s] 在[%d] 投喂了 %d 个 %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.Data["uname"].(string), d.roomID, int(js.Data["num"].(float64)), js.Data["giftName"].(string))
			case "USER_TOAST_MSG":
				fmt.Printf("[%s] %s: %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.Data["toast_msg"].(string))
			case "WIDGET_BANNER":
				continue
			default:
				fmt.Println(BytesToStringFast(uz))
			}
		case msg := <- d.heartBeatChannel:
			hb := colors.New(colors.Normal, 190, colors.Black)
			hb.Printf("[%s] HeartBeat...实时人气: %d. \n", time.Now().Format("2006-01-02 15:04:05"), ByteArrToDecimal(msg[16:]))
		case msg := <- d.serverNoticeChannel:
			if msg[7] == 0 {
				uz := msg[16:]
				js := new(receivedInfo)
				json.Unmarshal(uz, js)
				sn := colors.New(colors.Normal, 223, colors.Black)
				switch js.Cmd {
				case "NOTICE_MSG":
					sn.Printf("[%s] From Server %s: %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, js.MsgSelf)
				case "ROOM_RANK":
					sn.Printf("[%s] From Server %s: %d %s. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, uint32(js.Data["roomid"].(float64)), js.Data["rank_desc"].(string))
				case "ROOM_REAL_TIME_MESSAGE_UPDATE":
					sn.Printf("[%s] From Server %s: [%d] 关注: %d, 粉丝团: %d. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, uint32(js.Data["roomid"].(float64)), int(js.Data["fans"].(float64)), int(js.Data["fans_club"].(float64)))
				case "HOT_RANK_CHANGED":
					sn.Printf("[%s] From Server %s: [%d] Rank: %d, Trend: %d, Area name: %d, Countdown: %d. \n", time.Now().Format("2006-01-02 15:04:05"), js.Cmd, d.roomID, int(js.Data["rank"].(float64)), int(js.Data["trend"].(float64)), js.Data["area_name"].(string), int(js.Data["countdown"].(float64)))
				default:
					sn.Println(BytesToStringFast(uz))
				}
			} else {
				fmt.Println(msg)
			}
		}
	}
}