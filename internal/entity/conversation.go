package entity

type Conversation struct {
	Forwarder chan []byte

	Subscribers []Subscriber
}

func (c *Conversation) publish() {

}
