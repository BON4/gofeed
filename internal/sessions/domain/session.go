package domain

import "time"

// TODO: maby make separate microservice.

type Session struct {
	id           string
	refreshtoken string
	userAgent    string
	clientIp     string
	isBlocked    string
	expiresAt    *time.Time
	instance     []byte
}

type FactoryConfig struct {
}
