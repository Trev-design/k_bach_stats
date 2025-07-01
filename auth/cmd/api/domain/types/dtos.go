package types

type NewAccountRequestDTO struct {
	Name         string
	Email        string
	Password     string
	Confirmation string
}

type VerifyAccountDTO struct {
	Cookie     string
	VerifyCode string
	UserAgent  string
	IPAddress  string
}

type NewAccountSessionDTO struct {
	Access  string
	Refresh string
}

type LoginAccountDTO struct {
	Email     string
	Password  string
	UserAgent string
	IPAddress string
}

type ChangePasswordDTO struct {
	Email        string
	Password     string
	Confirmation string
	VerifyCode   string
	Cookie       string
	UserAgent    string
	IPAddress    string
}

type VerifySessionDTO struct {
	Name   string
	Cookie string
}

type RefreshSessionDTO struct {
	JWT       string
	Cookie    string
	IPAddress string
	UserAgent string
}

type RemoveSessionDTO struct {
	Cookie    string
	IPAddress string
	UserAgent string
}
