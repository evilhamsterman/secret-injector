package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/evilhamsterman/secret-injector/pkg/kube"
	"github.com/evilhamsterman/secret-injector/pkg/signals"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/informers"
)

func main() {
	ctx := signals.SetupSignalHandler()

	client, err := kube.GetKubeClient("")
	if err != nil {
		slog.Error("Failed to get kube client", slog.Any("error", err))
		os.Exit(1)
	}

	l, _ := labels.NewRequirement("coder.com/project", selection.In, []string{"playground"})
	s := labels.NewSelector().Add(*l)

	slog.Info("Label selector", slog.Any("selector", s.String()))
	labelOptions := informers.WithTweakListOptions(func(options *metav1.ListOptions) {
		options.LabelSelector = s.String()
	})

	namespaceOptions := informers.WithNamespace("default")
	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		client,
		time.Second*30,
		informers.WithNamespace("default"),
		labelOptions,
		namespaceOptions,
	)

	secretInformer := informerFactory.Core().V1().Secrets().Informer()

	controller := kube.NewController(ctx, secretInformer)
	if err := controller.Run(ctx); err != nil {
		slog.Error("Failed to run controller", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Controller stopped")

}

func init() {
	w := os.Stderr
	var handler slog.Handler
	handler = slog.NewTextHandler(w, nil)

	// If the output is a terminal, enable colors
	if isatty.IsTerminal(w.Fd()) {
		handler = tint.NewHandler(w, nil)
	}
	slog.SetDefault(slog.New(handler))

}
