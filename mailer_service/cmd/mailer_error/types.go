package mailererror

type Error struct {
	Assign Assignment
	MSG    ErrorMessage
	Info   ServerInformation
}

type Assignment struct {
	CorrelationID string
	UserID        string
}

type ErrorMessage struct {
	ErrorType    string
	ErrorPositon string
	Reason       string
	TimeStamp    string
}

type ServerInformation struct {
	Location       string
	JobDescription string
}
