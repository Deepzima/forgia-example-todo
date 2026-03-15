package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/grafana/grafana-app-sdk/operator"
	"github.com/grafana/grafana-app-sdk/simple"
	"github.com/zima/forgia-example-todo/pkg/app"
	"github.com/zima/forgia-example-todo/pkg/generated/manifestdata"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfigPath := flag.String("kubeconfig", "", "path to kubeconfig file (uses in-cluster config if empty)")
	metricsPort := flag.Int("metrics-port", 9090, "port for metrics/health server")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	logger.Info("starting todo operator")

	kubeConfig, err := buildKubeConfig(*kubeconfigPath)
	if err != nil {
		logger.Error("failed to build kubeconfig", "error", err)
		os.Exit(1)
	}

	runner, err := operator.NewRunner(operator.RunnerConfig{
		KubeConfig: *kubeConfig,
		MetricsConfig: operator.RunnerMetricsConfig{
			Enabled: true,
			MetricsServerConfig: operator.MetricsServerConfig{
				Port: *metricsPort,
			},
		},
	})
	if err != nil {
		logger.Error("failed to create operator runner", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	provider := simple.NewAppProvider(manifestdata.LocalManifest(), nil, app.NewTodoApp)
	logger.Info("operator runner starting", "metricsPort", *metricsPort)

	if err := runner.Run(ctx, provider); err != nil {
		logger.Error("operator exited with error", "error", err)
		os.Exit(1)
	}

	logger.Info("operator stopped")
}

func buildKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	// Try in-cluster config first, fall back to default kubeconfig location
	cfg, err := rest.InClusterConfig()
	if err != nil {
		home, _ := os.UserHomeDir()
		defaultPath := home + "/.kube/config"
		return clientcmd.BuildConfigFromFlags("", defaultPath)
	}
	return cfg, nil
}
