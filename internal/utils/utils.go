package utils

import (
	"github.com/gofrs/uuid"
)

func NewUUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

func VaildUUID(str string) (uuid.UUID, error) {
	return uuid.FromString(str)
}
