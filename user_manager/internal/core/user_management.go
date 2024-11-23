package core

type UserManagement interface {
	AddUser(payload []byte) error
	RemoveUser(payload []byte) error
}
