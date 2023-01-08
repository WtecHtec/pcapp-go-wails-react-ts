package main

import (
	"go-pc-react/bot"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type UserInfo struct {
	StatusCode int
	Unid       string
	Info       *openwechat.Self
	Firends    *openwechat.Friends
	Groups     *openwechat.Groups
}

func consoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	q.WriteFile(256, QR_PATH+uuid+".png")
	runtime.EventsEmit(WailApp.ctx, "login:code", uuid)
}

var WxBot *openwechat.Bot
var LoginUser *UserInfo

func WxLogin() {
	WxBot = openwechat.DefaultBot(openwechat.Desktop)
	login := false
	WxBot.UUIDCallback = consoleQrCode
	WxBot.LoginCallBack = func(body []byte) {
		res := string(body)
		if strings.Contains(res, "200") {
			runtime.EventsEmit(WailApp.ctx, "login:code", "200")
			login = true
		}
	}
	runtime.EventsOn(WailApp.ctx, "info:ready", func(optionalData ...interface{}) {
		if login == true {
			self, err := WxBot.GetCurrentUser()
			userInfo := &UserInfo{}
			if err != nil {
				userInfo.StatusCode = 201
				runtime.EventsEmit(WailApp.ctx, "info:get", userInfo)
				return
			}
			userInfo.StatusCode = 200
			userInfo.Unid = self.ID()
			userInfo.Info = self
			firends, errf := self.Friends()
			if errf != nil {
				firends = nil
			}
			if firends != nil {
				for i := 0; i < len(firends); i++ {
					firends[i].UserName = firends[i].User.ID()
				}
			}
			userInfo.Firends = &firends
			// groups, errg := self.Groups()
			// if errg != nil {
			// 	groups = nil
			// }
			// userInfo.Groups = &groups
			// fmt.Println("info:get", userInfo)
			// fmt.Println("info:get User Sex               ", self.Sex)
			// fmt.Println("info:get User NickName               ", self.NickName)
			// fmt.Println("info:get User UserName                         ", self.UserName)
			// fmt.Println("info:get firends", firends)
			LoginUser = userInfo
			runtime.EventsEmit(WailApp.ctx, "info:get", userInfo)
		}
	})

	runtime.EventsOn(WailApp.ctx, "config:ready", func(optionalData ...interface{}) {
		WailStroe.InitStroe()
		vaule := WailStroe.GetDataById(optionalData[0].(string))
		if vaule != nil {
			runtime.EventsEmit(WailApp.ctx, "config:get", vaule)
			return
		}
		runtime.EventsEmit(WailApp.ctx, "config:get", 201)
	})

	runtime.EventsOn(WailApp.ctx, "config:save", func(optionalData ...interface{}) {
		WailStroe.Save(optionalData)
		time.Sleep(time.Microsecond * 3)
		WailStroe.Update(optionalData[0].(map[string]interface{})["unid"].(string))
	})

	// 注册消息处理函数
	WxBot.MessageHandler = func(msg *openwechat.Message) {
		if WailStroe.CurrentData == nil {
			return
		}
		if msg.IsSendByFriend() && WailStroe.CurrentData["auto_reply"] == true {
			replyText(msg)
		}
		if msg.IsAt() && WailStroe.CurrentData["auto_reply"] == true && WailStroe.CurrentData["auto_reply_group"] == true {
			replyText(msg)
		}
	}

	WxBot.Login()
	WxBot.Block()
}

func replyText(msg *openwechat.Message) {

	if WailStroe.CurrentData["auto_bot"] == "tuling" {
		msgs := strings.Split(msg.Content, " ")
		info := msg.Content
		if msg.IsAt() && len(msgs) > 1 {
			info = msgs[1]
		}
		appKey := ""
		if WailStroe.CurrentData["tuling_api_key"] != nil {
			appKey = WailStroe.CurrentData["tuling_api_key"].(string)
		}
		msg.ReplyText(bot.TlBot(info, appKey))
		return
	}

	if WailStroe.CurrentData["auto_bot"] == "chatgpt" {
		return
	}

	if WailStroe.CurrentData["auto_desc"] != "" {
		msg.ReplyText(WailStroe.CurrentData["auto_desc"].(string))
	}

}
