package main

import (
	"log/slog"
	"os"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	startedTime := time.Now()

	// initial zap logger
	zCfg := zap.NewProductionConfig()
	zCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zLogger, _ := zCfg.Build(
		// zap fields will be converted to slog Attr
		zap.Fields(zap.String("app", "my-awesome-app")),
		zap.Fields(zap.Time("event_time", startedTime)),
		zap.Fields(zap.Int64("event_time_ns", startedTime.UnixNano())),
	)
	defer zLogger.Sync()

	// initial slog logger from zap-slog handler
	zHandler := zapslog.NewHandler(zLogger.Core(), nil)
	logger := slog.New(zHandler)

	version := os.Getenv("APP_VERSION")

	logger.Info("Hello world")
	// {"level":"info","ts":1699978120.9716039,"msg":"Hello world","app":"my-awesome-app","event_time":1699978120.9714801,"event_time_ns":1699978120971480000}

	logger.WithGroup("details").
		Info(
			"Hello world",
			// every field will be added to the group
			slog.Int("cpus", runtime.NumCPU()),
			slog.String("app_version", version),
			slog.Bool("is_debug", true),
			slog.Bool("is_prod", false),
			slog.Float64("pi", 3.14),
			slog.String("event_time_iso", time.Now().Format(time.RFC3339Nano)),
		)
	// {"level":"info","ts":1699978120.9718652,"msg":"Hello world","app":"my-awesome-app","event_time":1699978120.9714801,"event_time_ns":1699978120971480000,"details":{"cpus":8,"app_version":"5","is_debug":true,"is_prod":false,"pi":3.14,"event_time_iso":"2023-11-14T23:08:40.971632+07:00"}}

	logger.Debug("Hello world", slog.Int("cpus", runtime.NumCPU()), slog.String("app_version", version))
	// {"level":"debug","ts":1699978120.971888,"msg":"Hello world","app":"my-awesome-app","event_time":1699978120.9714801,"event_time_ns":1699978120971480000,"cpus":8,"app_version":"5"}

	logger.Info("Finished", slog.Duration("elapsed", time.Since(startedTime)))
	// {"level":"info","ts":1699978120.971893,"msg":"Finished","app":"my-awesome-app","event_time":1699978120.9714801,"event_time_ns":1699978120971480000,"elapsed":0.000412042}
}
