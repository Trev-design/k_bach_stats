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

type InvitationArchivePayload struct {
	Invitor      string `json:"invitor"`
	Workspace    string `json:"string"`
	InvitedGuest string `json:"invited_guest"`
	ID           string `json:"id"`
	InvitorID    string `json:"invitor_id"`
	Subject      string `json:"subject"`
	Message      string `json:"message"`
}

type JoinRequestArchivePayload struct {
	Requester   string `json:"requester"`
	Workspace   string `json:"workspace"`
	ID          string `json:"id"`
	RequesterID string `json:"requester_id"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

type StreamPayload struct {
	Kind    string
	Payload [][]byte
}
