package danmu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/asmcos/requests"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

var (
	getDanmuInfo = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%d&type=0"
)

type handShakeInfo struct {
	UID       uint8  `json:"uid"`
	Roomid    uint32 `json:"roomid"`
	Protover  uint8  `json:"protover"`
	Platform  string `json:"platform"`
	Clientver string `json:"clientver"`
	Type      uint8  `json:"type"`
	Key       string `json:"key"`
}

func (d *DanmuClient) connect() {
	r, err := requests.Get(fmt.Sprintf(getDanmuInfo, d.roomID))
	if err != nil {
		fmt.Println("request.Get DanmuInfo: ", err)
	}
	fmt.Println("获取弹幕服务器")
	token := gjson.Get(r.Text(), "data.token").String()
	hostList := []string{}
	gjson.Get(r.Text(), "data.host_list").ForEach(func(key, value gjson.Result) bool {
		hostList = append(hostList, value.Get("host").String())
		return true
	})
	hsInfo := handShakeInfo{
		UID:       0,
		Roomid:    d.roomID,
		Protover:  2,
		Platform:  "web",
		Clientver: "1.10.2",
		Type:      2,
		Key:       token,
	}
	for _, h := range hostList {
		d.conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s:443/sub", h), nil)
		if err != nil {
			fmt.Println("websocket.Dial: ", err)
			continue
		}
		fmt.Printf("连接弹幕服务器[%s]成功\n", hostList[0])
		break
	}
	if err != nil {
		fmt.Println("websocket.Dial Error")
	}
	jm, err := json.Marshal(hsInfo)
	if err != nil {
		fmt.Println("json.Marshal: ", err)
	}
	err = d.sendPackage(0, 16, 1, 7, 1, jm)
	if err != nil {
		fmt.Println("Conn SendPackage: ", err)
	}
	fmt.Printf("连接房间[%d]成功\n", d.roomID)
}

func (d *DanmuClient) heartBeat() {
	for {
		obj := []byte("5b6f626a656374204f626a6563745d")
		if err := d.sendPackage(0, 16, 1, 2, 1, obj); err != nil {
			fmt.Println("heart beat err: ", err)
			continue
		}
		time.Sleep(30 * time.Second)
	}
}
func (d *DanmuClient) receiveRawMsg() {
	for {
		_, msg, _ := d.conn.ReadMessage()
		if msg[7] == 2 {
			// fmt.Println("UnZlib..")
			msgs := splitMsg(zlibUnCompress(msg[16:]))
			for _, m := range msgs {
				d.unzlibChannel <- m
				// uz := m[16:]
				// js := new(receivedInfo)
				// json.Unmarshal(uz, js)
				// fmt.Println(js.Cmd)
				// fmt.Println(js)
				// NOTICE_MSG
				// if js.Cmd != "INTERACT_WORD" && js.Cmd != "ONLINE_RANK_V2" && js.Cmd != "DANMU_MSG" && js.Cmd != "SEND_GIFT" && js.Cmd != "COMBO_SEND" && js.Cmd != "ACTIVITY_BANNER_UPDATE_V2" && js.Cmd != "PANEL" && js.Cmd != "ONLINE_RANK_COUNT" && js.Cmd != "ENTRY_EFFECT" && js.Cmd != "LIVE_INTERACTIVE_GAME" && js.Cmd != "ONLINE_RANK_TOP3" && js.Cmd != "LIVE" {
				// if js.Cmd == "DANMU_MSG" {
				// 	fmt.Println(js.Cmd)
				// 	fmt.Println(BytesToStringFast(uz))
				// 	fmt.Println(js)
				// }
				// fmt.Println(zlibUnCompress(msg[16:]))
			}
		} else if msg[11] == 3 {
			d.heartBeatChannel <- msg
			// fmt.Println("HeartBeat..")
			// fmt.Println(ByteArrToDecimal(msg[16:]))
		} else {
			d.serverNoticeChannel <- msg
			// fmt.Println("Sth..")
			// ROOM_REAL_TIME_MESSAGE_UPDATE
			// NOTICE_MSG
			// ROOM_RANK
			// if msg[7] == 0 {
			// 	uz := msg[16:]
			// 	js := new(receivedInfo)
			// 	json.Unmarshal(uz, js)
			// 	if js.Cmd != "NOTICE_MSG" && js.Cmd != "ROOM_RANK" && js.Cmd != "ROOM_REAL_TIME_MESSAGE_UPDATE" && js.Cmd != "HOT_RANK_CHANGED" {
			// 		fmt.Println(*(*string)(unsafe.Pointer(&uz)))
			// 	}
			// } else {
			// 	fmt.Println(msg)
			// }
		}
	}
}

func (d *DanmuClient) Run() {
	d.connect()
	go d.process()
	go d.heartBeat()
	go d.receiveRawMsg()
}
