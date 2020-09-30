package protocol

//go:generate stringer -type=Type -linecomment

type Type uint8

const (
	Invalid    Type = iota // invalid
	Api                    // api
	Entrypoint             // entrypoint
)
