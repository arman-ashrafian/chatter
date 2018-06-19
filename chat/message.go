package chat

// Message ...
// username, body, and time of message
type Message struct {
	UserName  string `json:"userName"`
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}

// String ...
// print message
func (m *Message) String() string {
	return m.UserName + " at " + m.Timestamp + " says " + m.Body
}
