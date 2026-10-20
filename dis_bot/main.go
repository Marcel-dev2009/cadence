package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const Prefix = "!"
func main() {
 token := os.Getenv("DISCORD_TOKEN")
 if token == ""{
  fmt.Println("ERROR:DISCORD_TOKEN envars not set")
  return	
 }

 dg, err := discordgo.New("Bot " + token)
 if err != nil{
	fmt.Println("Error creating a new discord bot", err)
	return
 }
 	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent
  
  dg.AddHandler(createMessage) // register the message handler
  err = dg.Open()   //Assumption that dg.Open creates the open discord connection trying to access a pointer's submethod by derefrencing it before init might be the cause of the nul pointer derefrence
 if err != nil{
	fmt.Println("Error opening connection to discord:", err)
	return
 }
 fmt.Println("Spider server bot is now running in the terminal ")
 sc := make(chan os.Signal, 1)
 signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt) 
 <-sc 
 dg.Close()

}
 
func createMessage (s *discordgo.Session, m *discordgo.MessageCreate){
 if m.Author.ID == s.State.User.ID{
  return
 }
 if m.Content == Prefix+"help" {
		embed := &discordgo.MessageEmbed{
			Title:       "🕸️ FUTO Engineering Team Sync Bot Guide",
			Description: "I process natural language messages and turn them into downloadable `.ics` calendar cards instantly!",
			Color:       0xDC2626, // Spider-Verse Red
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🚀 Core Command",
					Value:  "`!create-event [Title] in my dorm room on [Time Phase]`",
					Inline: false,
				},
				{
					Name:   "💡 Examples That Work",
					Value:  "• `!create-event Coding marathon in my dorm room on Friday night.`\n• `!create-event UI/UX Design Sprint in my dorm room on tomorrow at 4pm.`\n• `!create-event Systems Review in my dorm room on next Tuesday 10:00.`",
					Inline: false,
				},
				{
					Name:   "⚙️ Setup Info",
					Value:  "Simply download the attached file to inject the event straight into your Apple or Google Calendar app.",
					Inline: false,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Built with Go & DiscordGo",
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

 if strings.HasPrefix(m.Content, Prefix + "create-event"){
  eventDetails := strings.TrimPrefix(m.Content, Prefix+"create-event")
  if eventDetails == "" {
	s.ChannelMessageSend(m.ChannelID, 
	"⚠️ Yo! Provide some details. Example: `!create-event Coding marathon in my dorm room on Friday night.`")
	return
  }
  go HandleCalenderOutput(s, m, &EventData{})
}

}
