package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

func main() {
	sess, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Panicf("Unable to create Discord session: %v", err)
	}

	sess.AddHandler(presenceEvent)

	if err := sess.Open(); err != nil {
		log.Panicf("Failed to Open Websocket Connection: %v", err)
	}

	// Handle Graceful Shutdown
	log.Info("Discord Bot now running. Exit with CTRL+C")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	if signal := <-sigChan; signal != nil {
		log.Infof("Recieved signal: %v", signal)
		log.Info("Shutting down Discord Bot")
	}

}

func presenceEvent(dSess *discordgo.Session, pUpdate *discordgo.PresenceUpdate) {
	g, err := dSess.State.Guild(pUpdate.GuildID)
	if err != nil {
		log.Error(err)
	}

	liveRole := ""
	for _, r := range g.Roles {
		if r.Name == "Live" {
			liveRole = r.ID
		}
	}

	p, err := dSess.State.Presence(g.ID, pUpdate.User.ID)
	if err != nil {
		log.Panicf("obtain presence error: %v", err)
	}

	if p.Game != nil && p.Game.Type == 1 {
		if err := dSess.GuildMemberRoleAdd(g.ID, pUpdate.User.ID, liveRole); err != nil {
			log.Errorf("tagging user error: %v", err)
		}
	} else {
		if err := dSess.GuildMemberRoleRemove(g.ID, pUpdate.User.ID, liveRole); err != nil {
			log.Errorf("tagging user error: %v", err)
		}
	}
}
