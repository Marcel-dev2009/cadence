package main

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)
// cleanTimeStr normalizes trickier natural phrases like "by 7PM" to "at 7PM"
func HandleCalenderOutput(s *discordgo.Session, m *discordgo.MessageCreate, event *EventData) {
	// 1. Regex Extraction matching incoming live message content
	re := regexp.MustCompile(`^!create-event\s+(.+?)\s+(?:in|at)\s+(.+)$`)
	matches := re.FindStringSubmatch(m.Content)
	if len(matches) == 0 {
		s.ChannelMessageSend(m.ChannelID, "❌ **Syntax Error:** Use format `!create-event [Title] in my dorm room on [Time Phase]`")
		return
	}

	event.Title = matches[1]
	rawTimeStr := matches[2]
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	result, err := w.Parse(rawTimeStr, time.Now())
	if err != nil || result == nil {
		s.ChannelMessageSend(m.ChannelID, "❌ **Parsing Error:** I couldn't interpret that date/time phrase. Try something like `Friday at 7PM`.")
		return
	}
	event.Time = result.Time

	// 3. Format strictly compliant UTC timestamps for iCalendar profiles
	formatStr := "20060102T150405Z"
	timeStmp := time.Now().UTC().Format(formatStr)
	startTimeStr := event.Time.UTC().Format(formatStr)
	endTimeStr := event.Time.Add(1 * time.Hour).UTC().Format(formatStr)

	// 4. Formatted legal iCalendar string buffer utilizing official CRLF (\(\r\\)n) line breaks
	var icsBuffer bytes.Buffer
	icsBuffer.WriteString("BEGIN:VCALENDAR\n") // Fixed: Removed trailing colon
	icsBuffer.WriteString("VERSION:2.0\n")
	icsBuffer.WriteString("PRODID:-//FUTO Team Engineering//marc64-x86bot//EN\n")
	icsBuffer.WriteString("BEGIN:VEVENT\n")
	icsBuffer.WriteString(fmt.Sprintf("UID:%d@futo.team\n", time.Now().UnixNano()))
	icsBuffer.WriteString(fmt.Sprintf("DTSTAMP:%s\n", timeStmp))
	icsBuffer.WriteString(fmt.Sprintf("SUMMARY:%s\n", event.Title))
	icsBuffer.WriteString(fmt.Sprintf("DTSTART:%s\n", startTimeStr))
	icsBuffer.WriteString(fmt.Sprintf("DTEND:%s\n", endTimeStr))
	icsBuffer.WriteString("END:VEVENT\n")
	icsBuffer.WriteString("END:VCALENDAR\n")

	// 5. Discord Presentation Layer Layout
	embed := &discordgo.MessageEmbed{
		Title: "🕸️ Event Transcribed & Synchronized!",
		Description: fmt.Sprintf(
			"**Title:** %s\n**Time:** %s\n\nClick the attached `.ics` file below to instantly inject this event straight into your phone or laptop calendar app!",
			event.Title,
			event.Time.Format("Monday, Jan 02 at 3:04 PM"), // Display in local human-readable layout
		),
		Color: 0xDC2626, // Spider-Man themed neon red hue
		Footer: &discordgo.MessageEmbedFooter{
			Text: "FUTO Engineering Team Sync Engine",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 6. Complete File dispatch payload assembly
	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		Files: []*discordgo.File{
			{
				Name:        "futo_eng_team_events.ics",
				ContentType: "text/calendar",
				Reader:      &icsBuffer,
			},
		},
	})

	if err != nil {
		fmt.Printf("Error Dispatching Message Package to Channel %s: %v\n", m.ChannelID, err)
	}
}
