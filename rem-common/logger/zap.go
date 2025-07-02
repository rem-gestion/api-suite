package logger

import (
	"os"

	"github.com/rem-gestion/rem-common/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New crea un *zap.Logger con:
//   - Nivel según cfg.Level
//   - ConsoleEncoder con colores
//   - Campo "service" para diferenciar la instancia
func New(cfg config.LoggerConfig, serviceName string) *zap.Logger {
	// 1) Parseo de nivel
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		lvl = zapcore.InfoLevel
	}

	// 2) Configuración del encoder para consola con colores
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // level en COLOR
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	// 3) Core: encoder + salida + nivel mínimo
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		lvl,
	)

	// 4) Armado del logger con caller y stacktrace a partir de ERROR
	lg := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// 5) Agregar campo "service" fijo
	return lg.With(zap.String("service", serviceName))
}
