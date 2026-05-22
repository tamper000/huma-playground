package utils

import nanoid "github.com/matoous/go-nanoid/v2"

func RandomID() (string, error) {
	return nanoid.Generate("abcdefghijklmnopqrstuvwxyz1234567890", 7)
}
