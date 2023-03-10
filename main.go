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
var trigger_time time.Time

func hourly_day_check() {
	for {
		if day_time_set {
			time.Sleep(1 * time.Second)
			fmt.Println("Done")
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
	hourly_day_check()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	// if time.Now().Hour() ==

	// Cleanly close down the Discord session.
	dg.Close()
}

func setup_bot(message string) (time.Weekday, int) {

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
	trigger_time, err := strconv.Atoi(extracted_data[2])
	if err != nil {

		//executes if there is any error
		fmt.Println(err)
	}
	day_time_set = true
	return day, trigger_time

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
		setup_bot(m.Content)
		var date_string = day.String()
		fmt.Println(trigger_time)
		var str1 = "Bot has been configured to RollCall for " + date_string + "s"
		s.ChannelMessageSend(m.ChannelID, str1)
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
