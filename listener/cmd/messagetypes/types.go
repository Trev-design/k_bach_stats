package messagetypes

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type ValidationNumberMessage struct {
	Email            string `json:"email"`
	UserName         string `json:"name"`
	ValidationNumber string `json:"number"`
	UserId           string `json:"id"`
}

type ChangeUserMessage struct {
	UserName string `json:"name"`
	UserId   string `json:"id"`
}

type DeleteUserMessage struct {
	UserId string `json:"id"`
}
