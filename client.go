package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type Message struct {
	Sender  string
	Content string
}

// Local message format (for UI)
type LocalMsg struct {
	Sender  string
	Content string
	Time    string
	System  bool
}

func main() {

	//----------------------------------
	// CONNECT TO RPC SERVER
	//----------------------------------
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	//----------------------------------
	//----------------------------------
	fmt.Print("Enter your ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	// Register with RPC server
	var registerReply string
	client.Call("ChatService.Register", id, &registerReply)

	//----------------------------------
	// LOCAL MESSAGE STORAGE
	//----------------------------------
	userMessages := make(map[string][]LocalMsg)

	clearScreen := func() {
		fmt.Print("\033[H\033[2J")
	}

	redrawChat := func() {
		clearScreen()
		fmt.Println("====== Live Chat ======")

		displayOrder := []string{"YOU", "SYSTEM"}

		for user := range userMessages {
			if user != "YOU" && user != "SYSTEM" {
				displayOrder = append(displayOrder, user)
			}
		}

		for _, user := range displayOrder {
			msgs, ok := userMessages[user]
			if !ok {
				continue
			}

			fmt.Printf("----- %s -----\n", user)
			for _, m := range msgs {

				if m.System {
					fmt.Printf("\033[1;34m[%s] %s\033[0m\n", m.Time, m.Content)

				} else if user == "YOU" {
					content := strings.Replace(m.Content, "["+id+"]", "[you]", 1)
					fmt.Printf("\033[1;32m[%s] %s\033[0m\n", m.Time, content)

				} else {
					fmt.Printf("\033[1;33m[%s] %s\033[0m\n", m.Time, m.Content)
				}
			}
			fmt.Println()
		}

		fmt.Printf("[%s] > ", id)
	}

	//----------------------------------
	// RECEIVING GOROUTINE
	//----------------------------------
	go func() {
		for {
			var incoming string
			err := client.Call("ChatService.Receive", id, &incoming)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nDisconnected from server. Press Enter to exit.\n")
				os.Exit(0)
			}

			incoming = strings.TrimSpace(incoming)
			if incoming == "" {
				continue
			}

			now := time.Now().Format("15:04:05")

			// System message (starts with **)
			if strings.HasPrefix(incoming, "**") {
				userMessages["SYSTEM"] = append(userMessages["SYSTEM"], LocalMsg{
					Sender:  "SYSTEM",
					Content: incoming,
					Time:    now,
					System:  true,
				})

			} else {
				// Normal message: [user]: text
				parts := strings.SplitN(incoming, ":", 2)
				sender := "OTHER"

				if len(parts) == 2 && strings.HasPrefix(parts[0], "[") && strings.HasSuffix(parts[0], "]") {
					sender = strings.Trim(parts[0], "[]")
				}

				if sender != id {
					userMessages[sender] = append(userMessages[sender], LocalMsg{
						Sender:  sender,
						Content: incoming,
						Time:    now,
						System:  false,
					})
				}
			}

			redrawChat()
		}
	}()

	//----------------------------------
	// SEND LOOP
	//--------------------------------
