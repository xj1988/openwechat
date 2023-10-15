package openwechat

import (
	"fmt"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.LoginCallBack = func(body CheckLoginResponse) {
		t.Log("login")
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
	}
}

func TestLogout(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.LoginCallBack = func(body CheckLoginResponse) {
		t.Log("login")
	}
	bot.LogoutCallBack = func(bot *Bot) {
		t.Log("logout")
	}
	bot.MessageHandler = func(msg *Message) {
		if msg.IsText() && msg.Content == "logout" {
			bot.Logout()
		}
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	bot.Block()
}

func TestMessageHandle(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.MessageHandler = func(msg *Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	bot.Block()
}

func TestFriends(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	friends, err := user.Friends()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(friends)
}

func TestGroups(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	groups, err := user.Groups()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(groups)
}

func TestPinUser(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	friends, err := user.Friends()
	if err != nil {
		t.Error(err)
		return
	}
	if friends.Count() > 0 {
		f := friends.First()
		f.Pin()
		time.Sleep(time.Second * 5)
		f.UnPin()
	}
}

func TestSender(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.MessageHandler = func(msg *Message) {
		if msg.IsSendByGroup() {
			fmt.Println(msg.SenderInGroup())
		} else {
			fmt.Println(msg.Sender())
		}
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	bot.Block()
}

// TestGetUUID
// @description: 获取登录二维码(UUID)
// @param t
func TestGetUUID(t *testing.T) {
	bot := DefaultBot(Desktop)

	uuid, err := bot.Caller.GetLoginUUID(bot.Context())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(uuid)
}

// TestLoginWithUUID
// @description: 使用UUID登录
// @param t
func TestLoginWithUUID(t *testing.T) {
	uuid := "oZZsO0Qv8Q=="
	bot := DefaultBot(Desktop, WithUUIDOption(uuid))
	err := bot.Login()
	if err != nil {
		t.Errorf("登录失败: %v", err.Error())
		return
	}
}

func TestBot(t *testing.T) {
	// 桌面模式
	bot := DefaultBot(Desktop)

	// 注册消息处理函数
	bot.MessageHandler = func(msg *Message) {
		fmt.Println(msg.Content)

		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	// 注册登陆二维码回调
	bot.UUIDCallback = PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
