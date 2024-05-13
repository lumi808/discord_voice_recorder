package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gordonklaus/portaudio"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	disToken := os.Getenv("DISCORD_TOKEN")

	session, err := discordgo.New("Bot " + disToken)
	if err != nil {
		log.Fatal(err)
	}

	session.AddHandler(pingPong)
	session.AddHandler(channelCreate)

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func pingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func channelCreate(s *discordgo.Session, cc *discordgo.ChannelCreate) {
	if cc.Type == discordgo.ChannelTypeGuildVoice {
		s.ChannelMessageSend(cc.Channel.ID, "Hello, recording started")
	}
}
