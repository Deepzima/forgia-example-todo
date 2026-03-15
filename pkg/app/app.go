package app

import (
	"log/slog"

	"github.com/grafana/grafana-app-sdk/app"
	"github.com/grafana/grafana-app-sdk/k8s"
	"github.com/grafana/grafana-app-sdk/simple"
	todov1 "github.com/zima/forgia-example-todo/pkg/generated/todo/v1"
	"github.com/zima/forgia-example-todo/pkg/repository"
	"github.com/zima/forgia-example-todo/pkg/watcher"
)

// NewTodoApp creates a new simple.App for managing Todo resources.
func NewTodoApp(cfg app.Config) (app.App, error) {
	logger := slog.Default().With("component", "todo-operator")

	clientGen := k8s.NewClientRegistry(cfg.KubeConfig, k8s.DefaultClientConfig())
	todoClient, err := todov1.NewTodoClientFromGenerator(clientGen)
	if err != nil {
		return nil, err
	}

	repo := repository.NewK8sTodoRepository(todoClient)
	todoWatcher := watcher.NewTodoWatcher(repo, logger)

	return simple.NewApp(simple.AppConfig{
		Name:            "todo-operator",
		KubeConfig:      cfg.KubeConfig,
		ClientGenerator: clientGen,
		ManagedKinds: []simple.AppManagedKind{{
			Kind: todov1.Kind(),
			Watcher: &simple.Watcher{
				AddFunc:    todoWatcher.Add,
				UpdateFunc: todoWatcher.Update,
				DeleteFunc: todoWatcher.Delete,
				SyncFunc:   todoWatcher.Sync,
			},
		}},
	})
}
