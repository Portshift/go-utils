package secret

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//go:generate $GOPATH/bin/mockgen -destination=./mock_secret.go -package=secret wwwin-github.cisco.com/eti/go-utils/k8s/secret Handler
type Handler interface {
	CreateOrUpdate(secret *corev1.Secret) (*corev1.Secret, error)
	Delete(meta metav1.ObjectMeta)
	Get(meta metav1.ObjectMeta) (*corev1.Secret, error)
}

// ensure type implement the requisite interface
var _ Handler = &HandlerImpl{}

type HandlerImpl struct {
	clientset kubernetes.Interface
}

func NewHandler(clientset kubernetes.Interface) *HandlerImpl {
	return &HandlerImpl{
		clientset: clientset,
	}
}

func (h HandlerImpl) CreateOrUpdate(secret *corev1.Secret) (*corev1.Secret, error) {
	if secret == nil {
		return nil, nil
	}

	var ret *corev1.Secret
	var err error

	if ret, err = h.clientset.CoreV1().Secrets(secret.GetNamespace()).Create(context.TODO(), secret, metav1.CreateOptions{}); err == nil {
		log.Infof("Secret was created successfully. name=%v, namespace=%v", secret.GetName(), secret.GetNamespace())
		return ret, nil
	}

	if !errors.IsAlreadyExists(err) {
		return nil, fmt.Errorf("failed to create secret: %v", err)
	}

	// secret already exists - update
	if ret, err = h.clientset.CoreV1().Secrets(secret.GetNamespace()).Update(context.TODO(), secret, metav1.UpdateOptions{}); err != nil {
		return nil, fmt.Errorf("failed to update secret: %v", err)
	}

	log.Infof("Secret was updated successfully. name=%v, namespace=%v", secret.GetName(), secret.GetNamespace())
	return ret, nil
}

func (h HandlerImpl) Delete(meta metav1.ObjectMeta) {
	if err := h.clientset.CoreV1().Secrets(meta.GetNamespace()).Delete(context.TODO(), meta.GetName(), metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		log.Errorf("Failed to delete secret. name=%v, namespace=%v: %v", meta.GetName(), meta.GetNamespace(), err)
	}
}

func (h HandlerImpl) Get(meta metav1.ObjectMeta) (*corev1.Secret, error) {
	return h.clientset.CoreV1().Secrets(meta.GetNamespace()).Get(context.TODO(), meta.GetName(), metav1.GetOptions{})
}
