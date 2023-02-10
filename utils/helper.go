package utils

import (
	"encoding/hex"
	"fmt"
	"go-fiber-starter/constants"
	"log"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func TimeTrack() string {
	elapsed := time.Since(time.Now()).Seconds()
	return fmt.Sprintf("%f", elapsed)
}

func PadNumberWithZero(value int) string {
	return fmt.Sprintf("%02d", value)
}

func TimestampString() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

func FilenameWithTimestamp(filename string) string {
	timeStampFormat := TimestampString()
	rootPath, _ := os.Getwd()
	return rootPath + "/" + constants.UploadDir + "/" + timeStampFormat + "-" + filename
}

func CurrentTimestamp() string {
	now := time.Now()
	return now.Format(constants.TimestampFormat)
}

func RemoveAllSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func RemoveNonAlphaNumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func Retry(attempts int, sleep time.Duration, f func() error) error {
	type stop struct {
		error
	}

	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

func InArray(val interface{}, array interface{}) (exists bool) {
	exists = false

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				exists = true
				return
			}
		}
	}

	return
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func RandomNumber(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s", msg)
	}
}

func IsValidUUID(u string) bool {
	if u == "00000000-0000-0000-0000-000000000000" {
		return false
	}

	_, err := uuid.Parse(u)
	return err == nil
}

func StripTags(content string) string {
	re := regexp.MustCompile(`<(.|\n)*?>`)
	strips := re.ReplaceAllString(content, "")
	return strings.ReplaceAll(strips, "&nbsp;", " ")
}
