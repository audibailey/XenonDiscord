package main

import (
	"time"
	"fmt"

	"github.com/audibailey/XenonDiscord"
	"github.com/audibailey/XenonDiscord/logger"

	"github.com/andersfylling/disgord"
	. "github.com/logrusorgru/aurora"
)

func main() {
	fmt.Println(Cyan(" _|      _|                                          _|_|_|                _|      "))
	fmt.Println(Cyan("   _|  _|      _|_|    _|_|_|      _|_|    _|_|_|    _|    _|    _|_|    _|_|_|_|  "))
	fmt.Println(Cyan("     _|      _|_|_|_|  _|    _|  _|    _|  _|    _|  _|_|_|    _|    _|    _|      "))
	fmt.Println(Cyan("   _|  _|    _|        _|    _|  _|    _|  _|    _|  _|    _|  _|    _|    _|      "))
	fmt.Println(Cyan(" _|      _|    _|_|_|  _|    _|    _|_|    _|    _|  _|_|_|      _|_|        _|_|  "))
	fmt.Println(Cyan("XenonBot is starting up!"))

	XD := XenonDiscord.New(nil)
	XD.Configure()
	XD.RegisterCommands(
		XenonDiscord.Command{
			Intent: "Ping",
			Command: Ping,
		},
	)
	XD.On()
}

func Ping(session disgord.Session, messageEvent *disgord.MessageCreate, params map[string]interface{}) {
	start := time.Now()
	msg, err := session.SendMsg(messageEvent.Ctx, messageEvent.Message.ChannelID, "Pong!")
	if err != nil {
		logger.Log.Error(err)
	} else {
		took := time.Since(start)
		session.SetMsgContent(messageEvent.Ctx, msg.ChannelID, msg.ID, fmt.Sprintf("Pong! - `%s`", took))
		logger.Log.Debug("The ping request was successful")
	}
} 
