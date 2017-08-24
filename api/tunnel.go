package api

// Will need to add this so that I can access the same database through my mac

import (
	"fmt"
)

// Endpoint is for any connection point
type Endpoint struct {
	Host string
	Port int
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}
