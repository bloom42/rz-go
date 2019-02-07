package bench

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/bloom42/rz-go"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

func newDisabledLogrus() *logrus.Logger {
	logger := newLogrus()
	logger.Level = logrus.ErrorLevel
	return logger
}

func newLogrus() *logrus.Logger {
	return &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
}

func newRz() rz.Logger {
	return rz.New(
		rz.Writer(ioutil.Discard),
		rz.With(func(e *rz.Event) { e.Timestamp() }),
	)
}

func newZerolog() zerolog.Logger {
	return zerolog.New(ioutil.Discard).With().Timestamp().Logger()
}

func newDisabledZerolog() zerolog.Logger {
	return newZerolog().Level(zerolog.Disabled)
}

var _tenInts = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var _tenStrings = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
var _tenTimes = []time.Time{time.Now()}

type _testUser struct {
	Username string
	Name     string
	Phone    string
}

var _oneUser = _testUser{
	Username: "lol",
	Name:     "lol2",
	Phone:    "lollol",
}

var _tenUsers = []_testUser{_testUser{}, _testUser{}, _testUser{}, _testUser{}, _testUser{},
	_testUser{}, _testUser{}, _testUser{}, _testUser{}, _testUser{}}
var errExample = errors.New("lolerror")

func logrus10Fields() logrus.Fields {
	return logrus.Fields{
		"int":     _tenInts[0],
		"ints":    _tenInts,
		"string":  _tenStrings[0],
		"strings": _tenStrings,
		"time":    _tenTimes[0],
		"times":   _tenTimes,
		"user1":   _oneUser,
		"user2":   _oneUser,
		"users":   _tenUsers,
		"error":   errExample,
	}
}

func zerolog10Fields(e *zerolog.Event) *zerolog.Event {
	return e.
		Int("int", _tenInts[0]).
		Ints("ints", _tenInts).
		Str("string", _tenStrings[0]).
		Strs("strings", _tenStrings).
		Time("time", _tenTimes[0]).
		Times("times", _tenTimes).
		Interface("user1", _oneUser).
		Interface("user2", _oneUser).
		Interface("users", _tenUsers).
		Err(errExample)
}

func zerolog10Context(c zerolog.Context) zerolog.Context {
	return c.
		Int("int", _tenInts[0]).
		Ints("ints", _tenInts).
		Str("string", _tenStrings[0]).
		Strs("strings", _tenStrings).
		Time("time", _tenTimes[0]).
		Times("times", _tenTimes).
		Interface("user1", _oneUser).
		Interface("user2", _oneUser).
		Interface("users", _tenUsers).
		Err(errExample)
}

func rz10Fields(e *rz.Event) {
	e.Int("int", _tenInts[0]).
		Ints("ints", _tenInts).
		String("string", _tenStrings[0]).
		Strings("strings", _tenStrings).
		Time("time", _tenTimes[0]).
		Times("times", _tenTimes).
		Interface("user1", _oneUser).
		Interface("user2", _oneUser).
		Interface("users", _tenUsers).
		Err(errExample)
}

var _testMessage = "hello world"

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context.")
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
	b.Run("bloom42/rz-go", func(b *testing.B) {
		logger := newRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage, nil)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(_testMessage)
			}
		})
	})
}

func Benchmark10FieldsContext(b *testing.B) {
	b.Logf("Logging with 10 fields in context")
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		fields := logrus10Fields()
		l := logger.WithFields(fields)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info(_testMessage)
			}
		})
	})
	b.Run("bloom42/rz-go", func(b *testing.B) {
		logger := newRz().Config(rz.With(rz10Fields))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage, nil)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := zerolog10Context(newZerolog().With()).Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(_testMessage)
			}
		})
	})
}

func Benchmark10Fields(b *testing.B) {
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		fields := logrus10Fields()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fields).Info(_testMessage)
			}
		})
	})
	b.Run("bloom42/rz-go", func(b *testing.B) {
		logger := newRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage, rz10Fields)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				zerolog10Fields(logger.Info()).Msg(_testMessage)
			}
		})
	})
}

func BenchmarkZl(b *testing.B) {
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Warn().
					Str("Hello", "world").
					Str("Hello2", "world").
					Str("Hello3", "world").
					Str("Hello4", "world").
					Msg(_testMessage)
			}
		})
	})
	b.Run("bloom42/rz", func(b *testing.B) {
		logger := newRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Warn(_testMessage, func(e *rz.Event) {
					e.String("Hello", "world")
					e.String("Hello2", "world")
					e.String("Hello3", "world")
					e.String("Hello4", "world")
				})
			}
		})
	})
}

func BenchmarkZlNoFields(b *testing.B) {
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Warn().
					Msg(_testMessage)
			}
		})
	})
	b.Run("bloom42/rz", func(b *testing.B) {
		logger := newRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Warn(_testMessage, nil)
			}
		})
	})
}

func BenchmarkZlNoFieldsNoMessage(b *testing.B) {
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Warn().
					Msg("")
			}
		})
	})
	b.Run("bloom42/rz", func(b *testing.B) {
		logger := newRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Warn("", nil)
			}
		})
	})
}

func BenchmarkZlLotOfFields(b *testing.B) {
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				zerolog10Fields(logger.Info()).Msg(_testMessage)
			}
		})
	})
	b.Run("bloom42/rz", func(b *testing.B) {
		logger := newRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Warn(_testMessage, rz10Fields)
			}
		})
	})
}
