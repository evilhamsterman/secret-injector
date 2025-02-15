package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/evilhamsterman/secret-injector/pkg/signals"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log := slog.Default()
	ctx := signals.SetupSignalHandler()

	config := clientcmd.NewDefaultClientConfigLoadingRules()
	cfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", config.GetStartingConfig)
	if err != nil {
		log.Error("Failed to load kubeconfig", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("Successfully loaded kubeconfig file")

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Error("Failed to create kubernetes client", slog.Any("error", err))
	}

	l, _ := labels.NewRequirement("coder.com/project", selection.In, []string{"playground", "jow"})
	s := labels.NewSelector().Add(*l)
	log.Info("Selector", slog.Any("selector", s.String()))
	labelOptions := informers.WithTweakListOptions(func(options *metav1.ListOptions) {
		options.LabelSelector = s.String()
	})

	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		client,
		time.Second*30,
		informers.WithNamespace("default"),
		labelOptions,
	)

	secretInformer := informerFactory.Core().V1().Secrets().Informer()
	secretInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			object, ok := obj.(*v1.Secret)
			if !ok {
				log.Error("Failed to cast object to secret")
				return
			}
			log.Info(
				"Secret added",
				slog.String("secret", object.Name),
				slog.String("data", fmt.Sprint(object.Data)),
				slog.String("namespace", object.Namespace),
			)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			object, ok := newObj.(*v1.Secret)
			if !ok {
				log.Error("Failed to cast object to secret")
				return
			}
			log.Info(
				"Secret updated",
				slog.String("secret", object.Name),
				slog.String("data", fmt.Sprint(object.Data)),
				slog.String("namespace", object.Namespace),
			)
		},
	})

	informerFactory.Start(ctx.Done())
	log.Info("Syncing cache")
	informerFactory.WaitForCacheSync(ctx.Done())
	log.Info("Cache synced")

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
