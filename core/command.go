package core

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/hako/durafmt"
)

type Environment int

const (
	twitch_env Environment = iota
	discord_env
)

type Command struct {
	Bot          *ZjeBot
	De           *DiscordEnvironment
	Te           *TwitchEnvironment
	PlainMessage string
	TMessage     *twitch.PrivateMessage
	DMessage     *discordgo.MessageCreate
}

func HandleCommand(env Environment, c *Command) {
	if !strings.HasPrefix(c.PlainMessage, COMMAND_PREFIX) {
		return
	}

	command := strings.Split(strings.Split(c.PlainMessage, COMMAND_PREFIX)[1], " ")[0]
	args := strings.Split(strings.Split(c.PlainMessage, COMMAND_PREFIX)[1], " ")[1:]

	switch command {
	case "ping":
		SendResponse("pong!", env, c)

	case "today":
		SendResponse(c.Bot.BotData.Today.Text, env, c)

	case "project":
		SendResponse(c.Bot.BotData.Project.Text, env, c)

	case "settoday":
		if !IsZielino(env, c) {
			SendResponse("You cannot use that command.", env, c)
		} else {
			todayData := strings.Join(args, " ")

			data := &Data{
				Today:   TodayData{Text: todayData},
				Project: c.Bot.BotData.Project,
			}

			err := c.Bot.WriteData(data)
			if err != nil {
				SendResponse("SetToday failed. See logs for more info.", env, c)
				log.Printf("SetToday failed: %s", err)
			} else {
				SendResponse(fmt.Sprintf("Today: %s", todayData), env, c)
			}
		}

	case "setproject":
		if !IsZielino(env, c) {
			SendResponse("You cannot use that command.", env, c)
		} else {
			projectData := strings.Join(args, " ")

			data := &Data{
				Today:   c.Bot.BotData.Today,
				Project: ProjectData{Text: projectData},
			}

			err := c.Bot.WriteData(data)
			if err != nil {
				SendResponse("SetProject failed. See logs for more info.", env, c)
				log.Printf("SetProject failed: %s", err)
			} else {
				SendResponse(fmt.Sprintf("Project: %s", projectData), env, c)
			}
		}

	case "time":
		SendResponse(time.Now().Format(time.RFC1123), env, c)

	case "shoutout":
		if !IsZielino(env, c) {
			SendResponse("You cannot use that command.", env, c)
		} else {
			channel := args[0]

			SendResponse(fmt.Sprintf("Follow %s @ https://twitch.tv/%s", channel, channel), env, c)
		}

	case "id":
		SendResponse(GetIdString(env, c), env, c)

	case "zjebot":
		SendResponse(PROJECT_ZJEBOT, env, c)

	case "website":
		SendResponse(WEBSITE_URL, env, c)

	case "os":
		SendResponse("Arch Linux", env, c)

	case "wm":
		SendResponse("Hyprland + Waybar", env, c)

	case "env":
		switch env {
		case twitch_env:
			SendResponse("Twitch", env, c)
		case discord_env:
			SendResponse("Discord", env, c)
		}

	case "uptime":
		t := time.Now()
		elapsed := t.Sub(*c.Bot.Start)

		duration, err := durafmt.ParseString(elapsed.String())
		if err != nil {
			log.Printf("Uptime failed: %s", err)
		}

		SendResponse(duration.String(), env, c)
	}
}

func SendResponse(response string, env Environment, c *Command) {
	switch env {
	case twitch_env:
		c.Te.Tc.Reply(c.TMessage.Channel, c.TMessage.ID, response)

	case discord_env:
		c.De.Dg.ChannelMessageSend(c.DMessage.ChannelID, fmt.Sprintf("<@%v> %s", c.DMessage.Author.ID, response))
	}
}

func IsZielino(env Environment, c *Command) bool {
	switch env {
	case twitch_env:
		return c.TMessage.User.Name == ZIELINO_TWITCH
	case discord_env:
		return c.DMessage.Author.ID == ZIELINO_DISCORD
	}

	return false
}

func GetIdString(env Environment, c *Command) string {
	var id string
	var username string

	switch env {
	case twitch_env:
		id = c.TMessage.User.ID
		username = c.TMessage.User.DisplayName
	case discord_env:
		id = c.DMessage.Author.ID
		username = c.DMessage.Author.Username
	}

	return fmt.Sprintf("%s (%s)", id, username)
}
