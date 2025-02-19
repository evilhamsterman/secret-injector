package kube

import (
	"context"
	"fmt"
	"log/slog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	informerFactory informers.SharedInformerFactory
	secretInformer  cache.SharedIndexInformer
}

// NewController initializes and creates a new controller
func NewController(
	ctx context.Context,
	informerFactory informers.SharedInformerFactory,
) *Controller {

	c := &Controller{
		informerFactory: informerFactory,
		secretInformer:  informerFactory.Core().V1().Secrets().Informer(),
	}

	c.secretInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addSecret,
		UpdateFunc: c.updateSecret,
		DeleteFunc: c.deleteSecret,
	})

	return c
}

// Run sets up the controller and starts the informers
func (c *Controller) Run(ctx context.Context) error {
	slog.Info("Starting controller")

	go c.secretInformer.Run(ctx.Done())

	if ok := cache.WaitForCacheSync(ctx.Done(), c.secretInformer.HasSynced); !ok {
		return fmt.Errorf("Failed to sync cache")
	}

	slog.Info("Cache synced")
	<-ctx.Done()
	slog.Info("Shutting down controller")

	return nil
}

func (c *Controller) updateSecret(oldObj, newObj interface{}) {
	newSecret := newObj.(*v1.Secret)

	slog.Info("Secret updated", slog.String("secret", newSecret.Name))
}

func (c *Controller) addSecret(obj interface{}) {
	secret := obj.(*v1.Secret)
	slog.Info("Secret added", slog.String("secret", secret.Name))
}

func (c *Controller) deleteSecret(obj interface{}) {
	secret := obj.(*v1.Secret)
	slog.Info("Secret deleted", slog.String("secret", secret.Name))
}
