package main

import (
	_ "github.com/apache/apisix-go-plugin-runner/cmd/go-runner/plugins"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
	"github.com/apache/apisix-go-plugin-runner/pkg/runner"
	"github.com/apisix-go-runner-plugin/plugins"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := runner.RunnerConfig{}
	cfg.LogLevel = zapcore.DebugLevel
	err := plugin.RegisterPlugin(&plugins.TimestampInserterGo{})
	if err != nil {
		log.Fatalf("failed to register plugin TimestampInserterGo, err: %s", err)
	}
	runner.Run(cfg)
}
