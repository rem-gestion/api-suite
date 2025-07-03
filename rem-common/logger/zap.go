package logger

import (
	"os"
	"strings"

	"github.com/rem-gestion/rem-common/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// tiny helpers para color ANSI
const (
	grey  = "\033[37m"
	red   = "\033[31m"
	yel   = "\033[33m"
	green = "\033[32m"
	reset = "\033[0m"
)

func colourise(level zapcore.Level, svc string) string {
	switch level {
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return red + "[" + svc + "]" + reset
	case zapcore.WarnLevel:
		return yel + "[" + svc + "]" + reset
	default:
		return green + "[" + svc + "]" + reset
	}
}

func New(cfg config.LoggerConfig, serviceName string) *zap.Logger {
	// 1) nivel
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		lvl = zapcore.InfoLevel
	}

	// 2) encoder legible
	encCfg := zapcore.EncoderConfig{
		TimeKey:      "ts",
		LevelKey:     "level",
		NameKey:      "svc",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeName: func(s string, enc zapcore.PrimitiveArrayEncoder) {
			// el nombre aún no tiene color: lo coloreamos al volcar
			enc.AppendString(s)
		},
	}
	encCfg.EncodeName = func(s string, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + strings.ToUpper(s) + "]")
	}
	// 3) core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg),
		zapcore.AddSync(os.Stdout),
		lvl,
	)

	// 4) logger base
	base := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	).Named(serviceName)

	// 5) envolvemos el core para colorear el “[svc]” según nivel
	coloured := base.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(c, func(entry zapcore.Entry) error {
			entry.LoggerName = colourise(entry.Level, entry.LoggerName)
			return nil
		})
	}))

	return coloured
}
