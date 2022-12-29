package logutils

import (
	"bufio"
	"log"
	"os"

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

func Warn(msg string) {
	logr.Zerolog.Warn().Msg(msg)
}

func Error(msg string, err error) {
	logr.Zerolog.Error().Err(err).Msg(msg)
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
	zeroLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()

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
