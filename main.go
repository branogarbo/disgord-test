package main

import (
	"flag"
	"log"
	"os"

	dg "github.com/andersfylling/disgord"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot token")
	flag.Parse()
}

func main() {
	client := dg.New(dg.Config{
		BotToken: Token,
		// Intents:  32509,
	})
	gateway := client.Gateway()
	defer gateway.StayConnectedUntilInterrupted()

	gateway.MessageCreate(func(s dg.Session, m *dg.MessageCreate) {
		if m.Message.Content == "!db" {
			var (
				channelID dg.Snowflake
				guildID   = m.Message.GuildID
			)

			guildChannels, err := s.Guild(guildID).GetChannels()
			if err != nil {
				log.Fatal(err)
			}

			for _, channel := range guildChannels {
				if channel.Type == dg.ChannelTypeGuildVoice {
					break
				}

				for _, user := range channel.Recipients {
					if user.ID == m.Message.Author.ID {
						channelID = channel.ID
						break
					}
				}

				if channelID != 0 {
					break
				}
			}

			audioFile, err := os.Open("./lesGooo.mp3")
			if err != nil {
				log.Fatal(err)
			}

			voice, err := client.Guild(guildID).VoiceChannel(channelID).Connect(false, true)
			if err != nil {
				log.Fatal(err)
			}

			voice.StartSpeaking()
			voice.SendDCA(audioFile)
			voice.StopSpeaking()

			audioFile.Close()
			voice.Close()
		}
	})
}
