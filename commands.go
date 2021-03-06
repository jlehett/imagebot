package main

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const outputPath = "image_script/images/output.png"
const weatherOutputPath = "output.txt"

// Test assures the client that the bot is up and running
func Test(session *discordgo.Session, msg *discordgo.MessageCreate) {
	session.ChannelMessageSend(msg.ChannelID, "System is up")
}

// Fallback will be called if a messgae references the bot
// but does not match any of the existing commands
func Fallback(session *discordgo.Session, msg *discordgo.MessageCreate) {
	errMsg := "Command does not exist or was improperly used"
	session.ChannelMessageSend(msg.ChannelID, errMsg)
}

// Collage runs the process to create a collage and send it back to the user
func Collage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	removeFile(outputPath)
	msgText := strings.ToLower(msg.Content)

	filename, err := downloadMessageAttachment(msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	runCollageScript(filename, getKeyword(msgText))

	sendMessageFile(session, msg.ChannelID, outputPath)
}

// Minecraft runes the collage process but with only minecraft blocks
func Minecraft(session *discordgo.Session, msg *discordgo.MessageCreate) {
	filename, err := downloadMessageAttachment(msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	runMinecraftScript(filename)

	sendMessageFile(session, msg.ChannelID, outputPath)
}

// Help sends a help message
func Help(session *discordgo.Session, msg *discordgo.MessageCreate) {
	helpMsg := `
**Commands:**
* @ImageBot minecraft [input_file]                 |    Minecraft command takes [input_file] and reconstructs it using only Minecraft block textuers
* @ImageBot collage [search keyword] [input_file]  |    Collage command takes [input_file] and reconstructs it using images found in Bing's image search for [search keyword]
***note:*** [input_file] signifies that a file should be sent as an attachment
`

	session.ChannelMessageSend(msg.ChannelID, helpMsg)
}

// WeatherToday sends the channel the current weather today
func WeatherToday(session *discordgo.Session, msg *discordgo.MessageCreate) {
	result := runWeatherCommand("today")

	session.ChannelMessageSend(msg.ChannelID, result)
}

// WeatherTomorrow sends the weather forcast for tomorrow
func WeatherTomorrow(session *discordgo.Session, msg *discordgo.MessageCreate) {
	result := runWeatherCommand("tomorrow")

	session.ChannelMessageSend(msg.ChannelID, result)
}

// WeatherWeek sends forcast for week
func WeatherWeek(session *discordgo.Session, msg *discordgo.MessageCreate) {
	result := runWeatherCommand("week")

	session.ChannelMessageSend(msg.ChannelID, result)
}

func runWeatherCommand(cmdStr string) string {
	cmd := exec.Command("python3", "weather_script/ShortDesc.py", cmdStr)

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return "Error reading weather information"
	}

	return ReadFile(weatherOutputPath)
}

func downloadMessageAttachment(msg *discordgo.MessageCreate) (string, error) {
	url := msg.Attachments[0].URL

	return DownloadImage(url)
}

func runCollageScript(filename string, keyword string) {
	cmd := exec.Command("python3", "image_script/collage.py", filename, keyword)

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}

func runMinecraftScript(filename string) {
	cmd := exec.Command("python3", "image_script/minecraft.py", filename)

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}

func getKeyword(msgText string) string {
	words := strings.Split(msgText, " ")
	return words[2]
}

func sendMessageFile(session *discordgo.Session, channelID, filename string) {
  file, err := os.Open(filename)

  if err != nil {
    return
  }

  defer file.Close()

  session.ChannelFileSend(channelID, filename, file)
}

func removeFile(filepath string) {
	cmd := exec.Command("rm", filepath)

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}