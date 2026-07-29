package utils

import (
	"encoding/base64"

	"github.com/gabriel-vasile/mimetype"
)

func BytesToUrl(bytes []byte) string {
	b64 := base64.StdEncoding.EncodeToString(bytes)
	mime := mimetype.Detect(bytes)

	return "data:" + mime.String() + ";base64," + b64
}
