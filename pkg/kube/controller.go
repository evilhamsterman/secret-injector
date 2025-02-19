package kube

import (
	"context"
	"fmt"
	"log/slog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	// Secret informer
	secretInformer cache.SharedIndexInformer
}

// NewController initializes and creates a new controller
func NewController(
	ctx context.Context,
	secretInformer cache.SharedIndexInformer,
) *Controller {

	c := &Controller{
		secretInformer: secretInformer,
	}

	c.secretInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleSecret,
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

func (c *Controller) handleSecret(obj interface{}) {
	secret := obj.(*v1.Secret)
	slog.Info("Secret added", slog.String("secret", secret.Name))
}

func (c *Controller) deleteSecret(obj interface{}) {
	secret := obj.(*v1.Secret)
	slog.Info("Secret deleted", slog.String("secret", secret.Name))
}
