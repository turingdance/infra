package pkkit

import (
	"strings"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func UseNanoID(size ...int) string {
	l := 32
	if len(size) > 0 {
		l = size[0]
	}
	id, _ := gonanoid.Generate(chars, l)
	return id
}

func UseUUID() string {
	u1 := uuid.New()
	return strings.ReplaceAll(u1.String(), "-", "")
}
