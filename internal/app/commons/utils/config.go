package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	_ "github.com/joho/godotenv/autoload" // buat jaga2
)

var (
	defclient = http.DefaultClient
)

//JSONMarshal is func
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
