package config

import "time"

const (
	// CookieLife tis the lifetime for a cookie (s)
	CookieLife int = int(24 * time.Hour)
	// JWTLife : JWT lifetime
	JWTLife time.Duration = 7 * 24 * time.Hour
)
