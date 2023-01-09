package logutils

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type Logger struct {
	Zerolog *zerolog.Logger
	Stdout  *os.File
	Stderr  *os.File
	Pipeout *os.File
	Pipeerr *os.File
}

func Debug(msg string) {
	logr.Zerolog.Debug().Msg(msg)
}

func Info(msg string) {
	logr.Zerolog.Info().Msg(msg)
}

func Infof(msg string, args ...interface{}) {
	logr.Zerolog.Info().Msgf(msg, args...)
}

func Warn(msg string) {
	logr.Zerolog.Warn().Msg(msg)
}

func Error(err error, msg string, args ...interface{}) {
	logr.Zerolog.Error().Err(err).Msgf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	logr.Zerolog.Fatal().Msgf(msg, args...)
}

var logr *Logger

func Close() error {
	os.Stdout = logr.Stdout
	os.Stderr = logr.Stderr
	log.SetOutput(logr.Stderr)
	return logr.Pipeout.Close()
}

func InitLogger() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		parts := strings.Split(file, "/")
		n := len(parts)
		if n < 3 {
			n = 3
		}
		file = strings.Join(parts[n-3:], "/")
		return file + ":" + strconv.Itoa(line)
	}

	zeroLogger := zerolog.New(os.Stderr).With().Timestamp().CallerWithSkipFrameCount(3).Logger()

	outR, outW, _ := os.Pipe()
	errR, errW, _ := os.Pipe()
	logr = &Logger{
		Zerolog: &zeroLogger,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Pipeout: outW,
		Pipeerr: errW,
	}
	os.Stdout = outW
	os.Stderr = errW
	log.SetOutput(errW)

	go func() {
		buf := bufio.NewReader(outR)
		for {
			line, _, _ := buf.ReadLine()
			if len(line) > 0 {
				logr.Zerolog.Info().Msg(string(line))
			}
		}
	}()

	go func() {
		buf := bufio.NewReader(errR)
		for {
			line, _, _ := buf.ReadLine()
			if len(line) > 0 {
				logr.Zerolog.Error().Msg(string(line))
			}
		}
	}()

	zeroLogger.Output(zerolog.ConsoleWriter{Out: logr.Stderr})
}
