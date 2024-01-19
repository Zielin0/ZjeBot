package core

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gempir/go-twitch-irc/v4"
)

type ZjeBot struct {
	BotSecrets    Secrets
	BotData       Data
	BotDataLoader DataLoader
	Start         *time.Time
}

func CreateBot(secrets Secrets, data Data, dataLoader *DataLoader, start *time.Time) *ZjeBot {
	return &ZjeBot{
		BotSecrets:    secrets,
		BotData:       data,
		BotDataLoader: *dataLoader,
		Start:         start,
	}
}

func (bot *ZjeBot) WriteData(data *Data) error {
	err := bot.BotDataLoader.WriteData(DATA_PATH, data)
	if err != nil {
		return err
	}

	err = bot.BotDataLoader.LoadData(DATA_PATH)
	if err != nil {
		return err
	}

	bot.BotData = bot.BotDataLoader.GetData()

	return nil
}

type DiscordEnvironment struct {
	Bot *ZjeBot
	Dg  *discordgo.Session
}

func (de *DiscordEnvironment) HandleDiscordMessage(message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	command := &Command{
		Bot:          de.Bot,
		De:           de,
		Te:           nil,
		PlainMessage: message.Content,
		TMessage:     nil,
		DMessage:     message,
	}

	HandleCommand(discord_env, command)
}

func InitDiscord(bot *ZjeBot) (*DiscordEnvironment, error) {
	dg, err := discordgo.New("Bot " + bot.BotSecrets.Discord.Auth)
	if err != nil {
		return nil, err
	}

	de := &DiscordEnvironment{
		Bot: bot,
		Dg:  dg,
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.MessageID == PINGS_MESSAGE_ID {
			if m.MessageReaction.Emoji.Name == "🔴" {
				err := s.GuildMemberRoleAdd(m.MessageReaction.GuildID, m.MessageReaction.UserID, PINGS_ROLE_ID)
				if err != nil {
					log.Fatalf("Failed adding pings role: %s", err)
				}
			}
		}
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
		if m.MessageID == PINGS_MESSAGE_ID {
			if m.MessageReaction.Emoji.Name == "🔴" {
				err := s.GuildMemberRoleRemove(m.MessageReaction.GuildID, m.MessageReaction.UserID, PINGS_ROLE_ID)
				if err != nil {
					log.Fatalf("Failed removing pings role: %s", err)
				}
			}
		}
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		de.HandleDiscordMessage(m)
	})

	err = dg.Open()
	if err != nil {
		return nil, err
	}

	return de, nil
}

type TwitchEnvironment struct {
	Bot *ZjeBot
	Tc  *twitch.Client
}

func (te *TwitchEnvironment) HandleTwitchMessage(message twitch.PrivateMessage) {
	if message.User.Name == te.Bot.BotSecrets.Twitch.Username {
		return
	}

	command := &Command{
		Bot:          te.Bot,
		De:           nil,
		Te:           te,
		PlainMessage: message.Message,
		TMessage:     &message,
		DMessage:     nil,
	}

	HandleCommand(twitch_env, command)
}

func InitTwitch(bot *ZjeBot) (*TwitchEnvironment, error) {
	tc := twitch.NewClient(bot.BotSecrets.Twitch.Username, bot.BotSecrets.Twitch.Auth)

	te := &TwitchEnvironment{
		Bot: bot,
		Tc:  tc,
	}

	tc.OnPrivateMessage(func(m twitch.PrivateMessage) {
		te.HandleTwitchMessage(m)
	})

	tc.Join(ZIELINO_TWITCH)

	return te, nil
}
