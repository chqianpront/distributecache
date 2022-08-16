package cache

import "encoding/json"

type Command struct {
	Type  int    `json:"type"`
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// command types
const (
	Get    = 0
	Add    = 1
	Update = 2
	Delete = 3
	GetOk  = 5
	Ok     = 6
	Ping   = 7
	Pong   = 8
)

func NewCommand(ty int, key string, value any) *Command {
	return &Command{
		Type:  ty,
		Key:   key,
		Value: value,
	}
}
func ParseCommand(cmd []byte) (*Command, error) {
	c := new(Command)
	err := json.Unmarshal(cmd, c)
	return c, err
}
