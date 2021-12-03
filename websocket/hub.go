package websocket

type Room struct {
	name string

	text string

	hub *Hub

	clients map[*Client]bool
}

type Hub struct {
	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client
	change     chan *Message
}

type Message struct {
	client *Client
	text   string
}

func NewHub() *Hub {
	return &Hub {
		rooms:    make(map[string]*Room),
		register: make(chan *Client),
		change:   make(chan *Message),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:

			roomName := client.roomName
			room, ok := hub.rooms[roomName]

			if !ok {
				room = &Room{name: roomName, hub: hub, clients: make(map[*Client]bool)}
				hub.rooms[roomName] = room
			}

			client.room = room
			room.clients[client] = true

			client.send <- room.text

		case client := <-hub.unregister:
			if _, ok := client.room.clients[client]; ok {
				delete(client.room.clients, client)
				close(client.send)
			}

		case message := <-hub.change:
			message.client.room.text = message.text
			for client := range message.client.room.clients {
				if client == message.client {
					continue
				}
				select {
				case client.send <- message.text:
				default:
					close(client.send)
					delete(client.room.clients, client)
				}
			}
		}
	}
}
