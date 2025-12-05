# Go RPC Chat Application 🚀

A real-time, command-line interface (CLI) chat application built using Go's standard `net/rpc` library. This project demonstrates how to implement a concurrent client-server architecture using Remote Procedure Calls (RPC), Goroutines, and Channels.

## 📋 Features

* **RPC-Based Communication:** Uses Go's native `net/rpc` for robust client-server interaction.
* **Real-time Messaging:** Instant message broadcasting to all connected clients.
* **System Notifications:** Automatic alerts when users join or leave the chat.
* **User List Management:** Updates users on who is currently online.
* **Concurrent Handling:** * **Server:** Handles multiple clients simultaneously using Mutexes and Goroutines.
    * **Client:** Separates message receiving (rendering) and message sending (input) to prevent UI blocking.
* **Terminal UI:** Color-coded output to distinguish between "You", "System", and "Other Users".

## 🛠️ Technology Stack

* **Language:** Go (Golang)
* **Networking:** TCP, `net/rpc`
* **Concurrency:** Goroutines (`go func`), Channels (`chan`), `sync.Mutex`

## 📂 Project Structure

```text
.
├── client.go   # The client-side application (UI and RPC calls)
├── server.go   # The centralized RPC server (State management and broadcasting)
└── README.md   # Project documentation
```
## 🚀 Getting Started

### Prerequisites
* [Go](https://go.dev/dl/) (version 1.18 or higher recommended)

### Installation

1.  Clone the repository:
    ```bash
    git clone (https://github.com/Abdallazayed2004/RPC---Chat-system.git)
    ```

2.  Ensure your directory contains both `server.go` and `client.go`.

## 💻 Usage

To run the chat application, you will need at least two terminal windows (one for the server, and one or more for clients).

### Step 1: Start the Server
Open a terminal and run the server. It will listen on port `1234`.

```bash
go run server.go
```
### Step 2: Start a Client
Open a new terminal window and run the client.
```bash
go run client.go
```
1- Enter your User ID when prompted.

2- The chat interface will load.

3- Type a message and press Enter to send.
### Step 3: Add More Users
Open additional terminal windows and run go run client.go again to simulate multiple users chatting with each other.

## 🧩 Architecture Overview
### Registration:
When a client starts, it calls ChatService.Register via RPC. The server creates a dedicated channel for that user.

### Broadcasting:

The server maintains a broadcast channel.

When a message is received (via SendMessage), it is pushed to the broadcast channel.

A dedicated goroutine on the server loops through all connected clients and forwards the message to their individual channels.

### Receiving:

The client calls ChatService.Receive. This is a blocking call on the server side (it waits until a message arrives in the user's channel).

Once a message arrives, the RPC call returns, and the client renders the message.

UI Rendering: The client uses a map userMessages to group messages by sender and redraws the screen on every update to keep the UI clean.
## Architecture

```
┌─────────┐     RPC      ┌─────────────┐     RPC      ┌─────────┐
│ Client1 │◄────────────►│             │◄────────────►│ Client2 │
└─────────┘              │   Server    │              └─────────┘
                         │             │
┌─────────┐     RPC      │  - Mutex    │
│ Client3 │◄────────────►│  - Channels │
└─────────┘              │  - History  │
                         └─────────────┘
```
