package config

import "time"

const (
	HTTPClientTimeoutShort  = 2 * time.Second
	HTTPClientTimeoutMedium = 5 * time.Second
	HTTPClientTimeoutLong   = 10 * time.Second

	TokenExpiration = 24 * time.Hour

	PRDefaultReviewersCap = 2

	DBPingTimeout  = 5 * time.Second
	FactoryTimeout = 10 * time.Second

	ExitCodeConfigError  = 1
	ExitCodeRuntimeError = 2
)
