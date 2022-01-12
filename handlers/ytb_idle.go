package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleIdle(bot *bot.Bot, info *youtube.LiveInfo) error {

	go discord.SendNewsEmbed(&discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:  youtubeChannelLink(info.ChannelId),
			Name: info.ChannelName,
		},
		Description: fmt.Sprintf("%s 的油管直播已结束。", info.ChannelName),
	})

	msg := message.NewSendingMessage().Append(qq.NewTextf("%s 的油管直播已结束。", info.ChannelName))
	return qq.SendGroupMessage(msg)
}

func init() {
	youtube.RegisterDataHandler(youtube.Idle, HandleIdle)
}
