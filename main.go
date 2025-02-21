package main

import (
	"os"
	"path"
	"runtime"
	"time"

	"penilaian-360/cmd"
	"penilaian-360/internal/app/commons/logHelper"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Join(file, "..")
	log.Info().Msg("path env =>" + rootPath + "/params/.env")
	err := godotenv.Load(rootPath + "/params/.env")
	if err != nil {
		log.Error().Msg("Error loading .env file")
	}

	log.Logger = log.With().Caller().Logger()
	if err := setTimezone("Asia/Jakarta"); err != nil {
		log.Error().Msg("timezone, err :" + err.Error()) // most likely timezone not loaded in Docker OS
	}

	logHelper.DebugMode = cast.ToBool(os.Getenv("DEBUG_MODE"))
	logHelper.LogTimeZone = os.Getenv("LOG_TIME_ZONE")

	cmd.Execute()
}

func setTimezone(tz string) error {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}
	time.Local = loc
	return nil
}
