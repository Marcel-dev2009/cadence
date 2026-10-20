package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleCalenderOutput(s *discordgo.Session, m *discordgo.MessageCreate, event *EventData) {
	// 1. Correct RFC Blueprint Layout Properties for ICS Files
	formatStr := "20060102T150405Z" 
	timeStmp := time.Now().UTC().Format(formatStr)
	
	// Create an event end time (defaulting to 1 hour after the parsed start time)
	startTimeStr := event.Time.UTC().Format(formatStr)
	endTimeStr := event.Time.Add(1 * time.Hour).UTC().Format(formatStr)

	// 2. Formatted matching legal iCalendar RFC specs
	icsContent := fmt.Sprintf(
		"BEGIN:VCALENDAR:\n"+
		"VERSION:2.0\n"+
		"PRODID:-//FUTO Team Engineering//marc64-x86bot//EN\n"+
		"BEGIN:VEVENT\n"+
		"UID:%d@futo.team\n"+ // Unique tracking ID
		"DTSTAMP:%s\n"+       // Corrected: Removed underscore
		"SUMMARY:%s\n"+       // Corrected: Changed from EVENT_NAME
		"DTSTART:%s\n"+       // Corrected: Changed from EVENT_TIME
		"DTEND:%s\n"+         // Added mandatory event boundary
		"END:VEVENT\n"+
		"END:VCALENDAR\n",
		time.Now().UnixNano(), timeStmp, event.Title, startTimeStr, endTimeStr,
	)

	fileBuffer := bytes.NewBufferString(icsContent) // to stream the data stream chunk by chunk 

	// 3. Discord Presentation Layer Layout
	embed := &discordgo.MessageEmbed{
		Title: "🕸️ Event Transcribed & Synchronized!",
		Description: fmt.Sprintf(
			"**Title:** %s\n**Time:** %s\n\nClick the attached `.ics` file below to instantly inject this event straight into your phone or laptop calendar app!", 
			event.Title, 
			event.Time.Format("Monday, Jan 02 at 3:04 PM"), // Human-readable string display
		),
		Color: 0xDC2626, // Spider-Man themed neon red hue 
		Footer: &discordgo.MessageEmbedFooter{
			Text: "FUTO Engineering Team Sync Engine",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 4. File dispatch payload assembly
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		Files: []*discordgo.File{
			{
				Name:        "futo_eng_team_events.ics",
				ContentType: "text/calendar",
				Reader:      fileBuffer,
			},
		},
	})
	
	if err != nil {
		fmt.Println("Error Uploading calendar file to discord:", err)
	}
}
