package symmetricHash

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func ToBcrypt(plainText string) string {
	moduleName := "services.symmetricHash"
	password := []byte(plainText)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 16)
	if err != nil {
		log.Error().Msg(moduleName + " err : " + err.Error())
		return ""
	}

	return string(hashedPassword)
}

func CompareBcrypt(hashedString, plainString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(plainString))
	if err != nil {
		return false
	}

	return true
}
