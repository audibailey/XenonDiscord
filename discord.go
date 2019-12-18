package XenonDiscord

import (
	_"time"

	_ "github.com/audibailey/XenonDiscord/config"
	"github.com/audibailey/XenonDiscord/nlp"
	"github.com/audibailey/XenonDiscord/logger"
	_"github.com/audibailey/XenonDiscord/utils"

	"github.com/andersfylling/disgord"
)

var commands map[string]*Command

func Master(session disgord.Session, messageEvent *disgord.MessageCreate) {
	message := messageEvent.Message
	logger.Log.Debug("A command has been requested from ", message.Author.Username, ":", message.Author.ID, " with the contents `", message.Content[1:], "`")
	
	intent, params, err := nlp.Nlp.Request(message.Content[1:], message.Author.ID.String())
	if err == nil {
		if command, exist := commands[intent]; exist {
			command.Command(session, messageEvent, params)
		} else {
			logger.Log.Error("Intent doesnt exist")
			message.Reply(messageEvent.Ctx, session, "An error has occured!")
		}
	} else {
		logger.Log.Error(err)
		message.Reply(messageEvent.Ctx, session, "An error has occured!")
	} 
}  

func AddCommand(command *Command) {
	commands[command.Intent] = command
}

func init() {
	commands = make(map[string]*Command)
}
