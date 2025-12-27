package commands

import (
	"fmt"

	"github.com/markwharris77/atlas/internal/tools"
)

func FollowLogs(target tools.Target, unitName string) error {
	if unitName == "" {
		return fmt.Errorf("name required")
	}

	return nil
}
