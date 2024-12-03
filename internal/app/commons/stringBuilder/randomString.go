package stringBuilder

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

var (
	LowerCaseRune = []rune("abcdefghijklmnopqrstuvwxyz")
	UpperCaseRune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	NumericRune   = []rune("1234567890")
)

func Generate(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"

	return str
}

func GenerateNumber(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := NumericRune

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "1234567890"

	return str
}

func GenerateAlphabet(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := LowerCaseRune

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "123"

	return str
}

func GenerateCapitalAlphabet(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := UpperCaseRune

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "123"

	return str
}

func GenerateAlphanumeric(length, numOfNumeric int) (string, error) {
	if length < 3 {
		return "", errors.New("minimum len is 3")
	}

	if numOfNumeric > length-2 {
		return "", errors.New("number of numeric character exceed (should not be more than length-2)")
	}

	randomNumeric := GenerateNumber(numOfNumeric)

	remainingLength := length - numOfNumeric

	//determine random number of uppercase alphabet
	upperCaseLen := rand.Intn(remainingLength/2) + 1
	randomUpperCase := GenerateCapitalAlphabet(upperCaseLen)

	//determine random number of lowercase alphabet
	lowerCaseLen := remainingLength - upperCaseLen
	randomLowerCase := GenerateAlphabet(lowerCaseLen)

	str := randomNumeric + randomUpperCase + randomLowerCase

	//shuffle string
	inRune := []rune(str)
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune), nil
}

func BackDate(from, to string) (yesterday1 time.Time, yesterday2 time.Time, err error) {
	date1, err := time.Parse("2006-01-02", from)
	if err != nil {
		err = errors.New("Error, parse date")
		return
	}
	date2, err := time.Parse("2006-01-02", to)
	if err != nil {
		err = errors.New("Error, parse date")
		return
	}
	duration := date2.Sub(date1)
	days := int(duration.Hours())
	if days == 0 {
		days = 24
	}
	yesterday1 = date1.Add(-time.Duration(days) * time.Hour)
	yesterday2 = date1.Add(-24 * time.Hour)
	return
}
