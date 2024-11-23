package types

type UserMessagePayload struct {
	Entity   string `json:"entity"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type SessionMessagePayload struct {
	Name      string `json:"name"`
	Account   string `json:"account"`
	SessionID string `json:"id"`
	AboType   string `json:"abo_type"`
}
