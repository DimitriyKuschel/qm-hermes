package providers

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"queue-manager/internal/structures"
)

type TypeEnum string

const (
	TypeGet  TypeEnum = "get"
	TypePost          = "post"
	TypeApp           = "app"
)

type Level string

const (
	TraceLevel Level = "trace"
	DebugLevel       = "debug"
	InfoLevel        = "info"
	WarnLevel        = "warn"
	ErrorLevel       = "error"
	FatalLevel       = "fatal"
	PanicLevel       = "panic"
)

type Logger interface {
	Errorf(t TypeEnum, format string, args ...interface{})
	Warnf(t TypeEnum, format string, args ...interface{})
	Debugf(t TypeEnum, format string, args ...interface{})
	Infof(t TypeEnum, format string, args ...interface{})
	Fatalf(t TypeEnum, format string, args ...interface{})
	Close()
}

func GetLogTypeByRequestType(rType string) TypeEnum {
	lType := TypeGet
	if rType == "POST" {
		lType = TypePost
	}
	return lType
}

func NewLogProvider(conf *structures.Config) (Logger, error) {
	log := Zerolog{conf: conf}
	err := log.init()
	if err != nil {
		return nil, err
	}
	return &log, nil
}

type Type struct {
	logger zerolog.Logger
	file   *os.File
}

type Zerolog struct {
	conf    *structures.Config
	loggers map[TypeEnum]Type
}

func (z *Zerolog) init() error {

	switch z.conf.Logger.Level {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}

	if z.conf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = "02/Jan/2006:15:04:05"

	z.loggers = make(map[TypeEnum]Type)
	for _, t := range []TypeEnum{TypeGet, TypeApp, TypePost} {

		p := filepath.Clean(z.conf.Logger.Dir + "/" + string(t) + ".log")
		file, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(z.conf.Logger.Mode))
		if err != nil {
			z.Close()
			return fmt.Errorf("Can`t open/create \"%s\" log file | File [%s] | %v\n", t, p, err)
		}

		var logger zerolog.Logger
		if z.conf.Debug {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			logger = zerolog.New(zerolog.MultiLevelWriter(zerolog.SyncWriter(consoleWriter), file)).With().Timestamp().Logger()
		} else {
			logger = zerolog.New(file).With().Timestamp().Logger()
		}

		z.loggers[t] = Type{
			logger: logger,
			file:   file,
		}
	}
	return nil
}

func (z *Zerolog) Errorf(t TypeEnum, format string, args ...interface{}) {
	logger := z.loggers[t].logger
	logger.Error().Msgf(format, args...)
}

func (z *Zerolog) Warnf(t TypeEnum, format string, args ...interface{}) {
	logger := z.loggers[t].logger
	z.write(logger.Warn(), format, args...)
}

func (z *Zerolog) Debugf(t TypeEnum, format string, args ...interface{}) {
	logger := z.loggers[t].logger
	z.write(logger.Debug(), format, args...)
}

func (z *Zerolog) Infof(t TypeEnum, format string, args ...interface{}) {
	logger := z.loggers[t].logger
	z.write(logger.Info(), format, args...)
}

func (z *Zerolog) Fatalf(t TypeEnum, format string, args ...interface{}) {
	logger := z.loggers[t].logger
	z.write(logger.Fatal(), format, args...)
}

func (z *Zerolog) write(event *zerolog.Event, format string, args ...interface{}) {
	if len(args) == 0 {
		event.Msg(format)
		return
	}
	event.Msgf(format, args...)
}

func (z *Zerolog) Close() {
	for _, logger := range z.loggers {
		err := logger.file.Close()
		if err != nil {
			fmt.Printf("error via file %s close %s\n", logger.file.Name(), err)
		}
	}
}
