package jwtimpl

type Adapter interface {
	Sign(id, role string) (string, error)
	Verify(jwtToken string) (string, string, error)
	CloseJWTService() error
	ComputeBackgroundService()
}
