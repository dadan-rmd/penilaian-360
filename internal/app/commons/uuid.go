package commons

import (
	"strings"

	guuid "github.com/satori/go.uuid"
)

type uuid struct{}

type Iuuid interface {
	V4() string
	V4Stripped() string
}

func NewUuid() Iuuid {
	return &uuid{}
}

func (c *uuid) V4() string {
	return guuid.NewV4().String()
}

func (c *uuid) V4Stripped() string {
	uuid := c.V4()
	return strings.ReplaceAll(uuid, "-", "")
}
