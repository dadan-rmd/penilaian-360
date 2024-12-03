package commons

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"
)

func FormattingIdPhoneNo(phoneNo string) string {
	num := phoneNo
	num = strings.TrimLeft(num, "62")
	num = strings.TrimLeft(num, "+62")

	if num != phoneNo {
		num = fmt.Sprintf("0%s", num)
	}

	return num
}

func GetSeed(salt string) int64 {
	h := fnv.New32a()
	h.Write([]byte(salt))
	salt = fmt.Sprint(h.Sum32())
	salt = "999999999" + salt
	salt = salt[len(salt)-9:]

	t := time.Now()
	timeStr := t.Format("0405.00000")
	timeStr = strings.Replace(timeStr, ".", "", -1)
	seedStr := salt + timeStr
	seed, _ := strconv.ParseInt(seedStr, 10, 64)
	return seed
}
