package core

type Sessionhandler interface {
	CheckSession(token string) error
	InitialAuth(token string) (string, error)
}
