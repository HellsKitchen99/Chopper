package domain

import "time"

type JWtConfig struct {
	Secret         []byte
	ExpirationTime time.Duration
	Issuer         string
	Audience       string
}
