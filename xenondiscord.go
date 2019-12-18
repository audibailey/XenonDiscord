package XenonDiscord

import (
	"context"
	"fmt"

	"github.com/audibailey/XenonDiscord/config"
	"github.com/audibailey/XenonDiscord/logger"
	"github.com/audibailey/XenonDiscord/nlp"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

type Command struct{
	Intent string
	Command func(disgord.Session, *disgord.MessageCreate, map[string]interface{})
}

type XenonDiscord struct {
	Log logger.Logger
	Ctx context.Context
	Client *disgord.Client
	Cmd []Command
	BotName string
	Conf *config.Config
}

func New(conf *config.Config) *XenonDiscord {
	config.Configure(conf)
	log := logger.Configure()
	nlp.Configure()
	return &XenonDiscord{
		Log: log,
		Ctx: context.Background(),
		Client: disgord.New(disgord.Config{
			Logger: log,
			BotToken: config.Conf.Discord.Token,
			LoadMembersQuietly: config.Conf.Discord.LoadMembersQuietly,
			ProjectName: config.Conf.Discord.ProjectName,
			DisableCache: config.Conf.Discord.DisableCache,
		}),
		BotName: config.Conf.Discord.ProjectName,
	}
}

func (xd *XenonDiscord) RegisterCommands(cmds ...Command) {
	for _, cmd := range cmds {
		AddCommand(&Command{
			Intent: cmd.Intent,
			Command: cmd.Command,
		})
	}
}

func (xd *XenonDiscord) Configure() {
	xd.Log.Info(fmt.Sprintf("Setting up %s", xd.BotName))

	botID, err := xd.Client.Myself(xd.Ctx)
	if err != nil {
		xd.Log.Panic("Failed to get the bot's ID.")
	}

	filter, err := std.NewMsgFilter(xd.Ctx, xd.Client)
	if err != nil {
		xd.Log.Panic("Failed to set the bot's filter.")
	}
	filter.SetPrefix(fmt.Sprintf("<@!%s>", botID.ID))

	xd.Client.On("MESSAGE_CREATE", filter.NotByBot, filter.HasPrefix, filter.StripPrefix, Master)
}

func (xd *XenonDiscord) On() {
	defer xd.Client.StayConnectedUntilInterrupted(xd.Ctx)

	xd.Client.On("READY", func(s disgord.Session, dR *disgord.Ready) {
		logger.Log.Info(fmt.Sprintf("%s online", xd.BotName))
		xd.Client.UpdateStatusString(fmt.Sprintf("@%s help", xd.BotName))
	})
}
