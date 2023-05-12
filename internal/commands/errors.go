package commands

import "errors"

var (
	ErrCommandHandlerNotFound = errors.New("command handler not found")
)

func NewErrInvalidCommandType(ctype Type) error {
	return errors.New("invalid command type: " + string(ctype))
}
