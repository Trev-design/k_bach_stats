package webtypes

type NewAccountRepresentation struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Confirmation string `json:"confirmation"`
}

type VerifyRepresentation struct {
	Code string `json:"verify"`
}

type LoginRepresentation struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewPasswordRepresentation struct {
	Email string `json:"email"`
}

type ChangePasswordRepresentation struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Confirmation string `json:"confirmation"`
	VerifyCode   string `json:"verify"`
}
