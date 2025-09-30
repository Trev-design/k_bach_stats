package jwtimpl

// the adapter to make jwts and give the user from accounts access
type Adapter interface {
	// signs a new json web token if the sign failed you get an error
	Sign(id, role string) (string, error)

	// tries to verify a json web token. if the verify failed you get an error
	Verify(jwtToken string) (string, string, error)

	// closes the jwt service by server shutdown
	CloseJWTService() error

	// registers and compute background services like intervalls for key change
	ComputeBackgroundService()
}
