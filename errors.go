package gonfiguration

import "errors"

var (
	ErrParsingConfig          = errors.New("dafuq? error parsing config")
	ErrTargetNotPointer       = errors.New("yo, the destination ain't a pointer")
	ErrDestinationNotStruct   = errors.New("what the hell? expected a struct, but this ain't one")
	ErrFieldIsNotOfSimpleType = errors.New("yo, this crappy field is way too complex. we like to keep it simple here")
)
