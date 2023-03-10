package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

// This flag Stringvar is used for adding things from CLI
func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

var day time.Weekday
var day_time_set = false
var trigger_time_hour int
var trigger_time_minute int
var timezone string
var channel_id string

func check_day_and_time() bool {
	var today = time.Now()
	if today.Weekday() == day && today.Hour() == trigger_time_hour && today.Minute() == trigger_time_minute {
		return true
	}
	return false
}

func hourly_day_check(s *discordgo.Session) {
	for {
		var is_time = check_day_and_time()
		if day_time_set {
			if is_time {
				s.ChannelMessageSend(channel_id, "YOU READY TO GAMER?!")
				time.Sleep(1 * time.Minute)
			}
		}
	}
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	hourly_day_check(dg)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	// Cleanly close down the Discord session.
	dg.Close()
}

func setup_bot(message string) (time.Weekday, int, int) {

	date_comparison := map[string]time.Weekday{
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
		"Sunday":    time.Sunday}
	var extracted_data = strings.Split(message, " ")
	day = date_comparison[extracted_data[1]]
	trigger_time_hour, hour_err := strconv.Atoi(extracted_data[2])
	if hour_err != nil {

		//executes if there is any error
		fmt.Println(hour_err)
	}
	trigger_time_minute, min_err := strconv.Atoi(extracted_data[3])
	if min_err != nil {

		//executes if there is any error
		fmt.Println(min_err)
	}
	day_time_set = true
	return day, trigger_time_hour, trigger_time_minute

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "setup_bot") {
		day, trigger_time_hour, trigger_time_minute = setup_bot(m.Content)
		var date_string = day.String()
		var str1 = "Bot has been configured to RollCall for " + date_string + "s"
		s.ChannelMessageSend(m.ChannelID, str1)
		channel_id = m.ChannelID
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
