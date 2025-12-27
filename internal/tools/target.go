package tools

import (
	"fmt"
)

type Target struct {
	User string
	Host string
	Port int
}

func (t Target) Addr() string {
	if t.User == "" {
		return t.Host
	}
	return fmt.Sprintf("%s@%s", t.User, t.Host)
}

func (t Target) PortOrDefault() int {
	if t.Port == 0 {
		return 22
	}
	return t.Port
}
