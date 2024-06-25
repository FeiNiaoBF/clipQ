package cmd

type ClipQ struct {
	// Registered clients.
	Clients map[*Client]bool
	// ClipMsg clipboard message from client.
	ClipMsg chan []byte
	// Register requests from the clients.
	Register chan *Client

	Unregister chan *Client
}

func NewClipQ() *ClipQ {
	return &ClipQ{
		Clients:    make(map[*Client]bool),
		ClipMsg:    make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (c *ClipQ) Run() {
	for {
		select {
		case client := <-c.Register:
			c.Clients[client] = true
		case client := <-c.Unregister:
			if _, ok := c.Clients[client]; ok {
				delete(c.Clients, client)
				close(client.Write)
			}
		case msg := <-c.ClipMsg:
			for client := range c.Clients {
				select {
				case client.Write <- msg:
				default:
					delete(c.Clients, client)
					close(client.Write)
				}
			}
		}
	}
}
