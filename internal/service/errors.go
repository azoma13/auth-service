package service

import "fmt"

var (
	ErrDifferentUserAgent     = fmt.Errorf("user-agent is different")
	ErrDifferentXForwardedFor = fmt.Errorf("x-forwarded-for is different")
)
