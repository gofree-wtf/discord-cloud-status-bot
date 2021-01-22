package common

import "github.com/bwmarrin/discordgo"

const CreateMessageKey = "createMessage"

type CreateMessage struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
}
