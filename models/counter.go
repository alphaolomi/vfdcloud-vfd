package models

import "time"

type (
	Counter struct {
		TaxID string
		Env   string
		Date  time.Time
		GC    *int64
		DC    *int64
	}
)
