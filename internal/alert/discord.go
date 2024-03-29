package alert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kitabisa/teler/pkg/errors"
)

// token can be webhook URL if webhooked is true
func toDiscord(token string, channel string, color string, version string, log map[string]string, webhooked bool) {
	col, err := strconv.Atoi(color)
	if err != nil {
		errors.Show(err.Error())
	}

	field := []*discordgo.MessageEmbedField{
		{
			Name: "Request",
			Value: fmt.Sprintf(
				"%s %s %s",
				log["request_method"], log["request_uri"], log["request_protocol"],
			),
		},
		{
			Name:   "Date",
			Value:  log["time_local"],
			Inline: true,
		},
		{
			Name:   "IP Address",
			Value:  log["remote_addr"],
			Inline: true,
		},
		{
			Name:  "User-Agent",
			Value: log["http_user_agent"],
		},
		{
			Name:  "Referrer",
			Value: log["http_referer"],
		},
		{
			Name:   "Status code",
			Value:  log["status"],
			Inline: true,
		},
		{
			Name:   "Bytes sent",
			Value:  log["body_bytes_sent"],
			Inline: true,
		},
	}
	embed := &discordgo.MessageEmbed{
		Title: ":warning: teler Alert",
		Description: fmt.Sprintf(
			"**%s**",
			log["category"],
		),
		URL:   "https://github.com/kitabisa/teler",
		Color: col,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://user-images.githubusercontent.com/25837540/97091757-7200d880-1668-11eb-82c4-e5c4971d2bc8.png",
		},
		Fields: field,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "teler " + version,
		},
	}

	if webhooked {
		msg := discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{embed},
		}

		data, err := json.Marshal(msg)
		if err != nil {
			errors.Exit(err.Error())
		}

		// nolint:errcheck
		http.Post(token, "application/json", strings.NewReader(string(data)))
	} else {
		discord, err := discordgo.New("Bot " + token)
		if err != nil {
			errors.Exit(err.Error())
		}

		// nolint:errcheck
		discord.ChannelMessageSendEmbed(channel, embed)
	}
}
