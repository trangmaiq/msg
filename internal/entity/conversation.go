package entity

type Conversation struct {
	// Inbound messages from the Clients
	Broadcast chan []byte
	// Register requests from the Clients
	Register chan *Client
	// Unregister requests from the Clients
	Unregister chan *Client
	// registered Clients
	Clients map[*Client]struct{}
}

func NewConversation() *Conversation {
	return &Conversation{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]struct{}),
	}
}

func (c *Conversation) Start() {
	for {
		select {
		case client := <-c.Register:
			c.Clients[client] = struct{}{}
		case client := <-c.Unregister:
			if _, ok := c.Clients[client]; ok {
				delete(c.Clients, client)
				close(client.Send)
			}
		case message := <-c.Broadcast:
			for client := range c.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(c.Clients, client)
				}
			}
		}
	}
}
