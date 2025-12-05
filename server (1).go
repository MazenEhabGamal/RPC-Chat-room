package main

import (
	"fmt"
	"net"
	"net/rpc"
	"strings"
	"sync"
)

type Message struct {
	Sender  string
	Content string
}

type BroadcastMsg struct {
	Sender   string
	Content  string
	IsSystem bool
}

type ChatService struct {
	mu        sync.Mutex
	clients   map[string]chan string
	broadcast chan BroadcastMsg
}

func NewChatService() *ChatService {
	s := &ChatService{
		clients:   make(map[string]chan string),
		broadcast: make(chan BroadcastMsg, 100),
	}

	// Main broadcaster goroutine (same as TCP version)
	go func() {
		for msg := range s.broadcast {
			s.mu.Lock()

			// Generate updated user list like your TCP server
			var userList []string
			for user := range s.clients {
				userList = append(userList, user)
			}
			userListMsg := fmt.Sprintf("** Current users in chat: %s **",
				strings.Join(userList, ", "))

			for username, ch := range s.clients {

				// Skip sender (no self echo)
				if username == msg.Sender && !msg.IsSystem {
					continue
				}

				// Send content
				ch <- msg.Content

				// System messages also send user list
				if msg.IsSystem {
					ch <- userListMsg
				}
			}

			s.mu.Unlock()
		}
	}()

	return s
}

// ----------- RPC Methods ------------

// User registers (same behavior as TCP handleClient)
func (s *ChatService) Register(username string, reply *string) error {
	s.mu.Lock()
	s.clients[username] = make(chan string, 20)
	s.mu.Unlock()

	// Broadcast join message
	s.broadcast <- BroadcastMsg{
		Sender:   username,
		Content:  fmt.Sprintf("** User [%s] joined the chat **", username),
		IsSystem: true,
	}

	fmt.Println("[SERVER] User joined:", username)
	*reply = "ok"
	return nil
}

// User sends a message
func (s *ChatService) SendMessage(msg Message, reply *string) error {
	s.broadcast <- BroadcastMsg{
		Sender:   msg.Sender,
		Content:  fmt.Sprintf("[%s]: %s", msg.Sender, msg.Content),
		IsSystem: false,
	}

	*reply = "sent"
	return nil
}

// Blocking receive (RPC version of TCP ReadString)
func (s *ChatService) Receive(username string, reply *string) error {
	s.mu.Lock()
	ch, ok := s.clients[username]
	s.mu.Unlock()

	if !ok {
		*reply = ""
		return nil
	}

	// Blocks until message arrives
	msg := <-ch
	*reply = msg
	return nil
}

// User exits
func (s *ChatService) Logout(username string, reply *string) error {
	s.mu.Lock()
	delete(s.clients, username)
	s.mu.Unlock()

	// Notify others
	s.broadcast <- BroadcastMsg{
		Sender:   username,
		Content:  fmt.Sprintf("** User [%s] left the chat **", username),
		IsSystem: true,
	}

	fmt.Println("[SERVER] User left:", username)
	*reply = "ok"
	return nil
}

// ----------- MAIN ------------

func main() {
	service := NewChatService()
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	fmt.Println("RPC Chat Server running on port 1234...")

	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}
