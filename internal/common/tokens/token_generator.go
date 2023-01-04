package tokens

import "time"

type Generator[T InstanceCredentials] interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(instance T, duration time.Duration) (string, *Payload[T], error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload[T], error)
}
