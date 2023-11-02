package gonfiguration

import "errors"

var (
	ErrParsingConfig  = errors.New("error parsing config")
	ErrDefaultsNotSet = errors.New("defaults not set. use gonfiguration.SetDefaults")
)
