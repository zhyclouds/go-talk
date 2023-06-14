package utils

import (
	"fmt"
	"github.com/google/uuid"
)

func GetUUID() string {
	u := uuid.New()
	return fmt.Sprintf("%x", u)
}
