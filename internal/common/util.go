package common

import gonanoid "github.com/matoous/go-nanoid/v2"

func NewID() string {
	id, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz1234567890", 16)
	return id
}
