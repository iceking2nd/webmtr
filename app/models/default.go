package models

import "time"

type DefaultParameters struct {
	COUNT            int
	TIMEOUT          time.Duration
	INTERVAL         time.Duration
	HOP_SLEEP        time.Duration
	MAX_HOPS         int
	MAX_UNKNOWN_HOPS int
	RING_BUFFER_SIZE int
	PTR_LOOKUP       bool
	JsonFmt          bool
	SrcAddr          string
}

var DefParams DefaultParameters
