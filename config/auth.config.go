package config

import "time"

const (
	// CookieLife tis the lifetime for a cookie (s)
	CookieLife int = 3600
	// JWTLife : JWT lifetime
	JWTLife time.Duration = 2 * time.Hour
)
