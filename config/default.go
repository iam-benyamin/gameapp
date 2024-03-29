package config

import "time"

var defaultConfig = map[string]interface{}{
	"auth.refresh_subject":                  RefreshTokenSubject,
	"auth.access_subject":                   AccessTokenSubject,
	"auth.access_expiration_time":           AccessTokenExpireDuration,
	"auth.refresh_expiration_time":          RefreshTokenExpireDuration,
	"application.graceful_shutdown_timeout": 5 * time.Second,
}
