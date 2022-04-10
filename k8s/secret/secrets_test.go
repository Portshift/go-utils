package secret

import (
	"fmt"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	client_testing "k8s.io/client-go/testing"
)

func createTestSecret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "name",
			Namespace: "namespace",
		},
		StringData: map[string]string{
			"key": "val",
		},
		Type: corev1.SecretTypeOpaque,
	}
}

func TestHandlerImpl_CreateOrUpdate(t *testing.T) {
	secret := createTestSecret()
	type fields struct {
		clientset kubernetes.Interface
	}
	type args struct {
		secret *corev1.Secret
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *corev1.Secret
		wantErr bool
	}{
		{
			name: "successful create",
			fields: fields{
				clientset: successfulCreateRequestReactor(fake.NewSimpleClientset(), secret),
			},
			args: args{
				secret: secret,
			},
			want:    secret,
			wantErr: false,
		},
		{
			name: "failed create",
			fields: fields{
				clientset: failedCreateRequestReactor(fake.NewSimpleClientset(), nil),
			},
			args: args{
				secret: &corev1.Secret{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "successful update",
			fields: fields{
				clientset: successfulUpdateRequestReactor(fake.NewSimpleClientset(), secret),
			},
			args: args{
				secret: secret,
			},
			want:    secret,
			wantErr: false,
		},
		{
			name: "failed update",
			fields: fields{
				clientset: failedUpdateRequestReactor(fake.NewSimpleClientset(), nil),
			},
			args: args{
				secret: secret,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil secret",
			fields: fields{
				clientset: fake.NewSimpleClientset(),
			},
			args: args{
				secret: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HandlerImpl{
				clientset: tt.fields.clientset,
			}
			got, err := h.CreateOrUpdate(tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateOrUpdate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func failedCreateRequestReactor(clientset kubernetes.Interface, ret runtime.Object) kubernetes.Interface {
	clientset.(*fake.Clientset).Fake.ReactionChain = []client_testing.Reactor{
		&client_testing.SimpleReactor{
			Verb:     "create",
			Resource: "secrets",
			Reaction: func(action client_testing.Action) (bool, runtime.Object, error) {
				return true, ret, errors.NewBadRequest("test")
			},
		},
	}

	return clientset
}

func successfulCreateRequestReactor(clientset kubernetes.Interface, ret runtime.Object) kubernetes.Interface {
	clientset.(*fake.Clientset).Fake.ReactionChain = []client_testing.Reactor{
		&client_testing.SimpleReactor{
			Verb:     "create",
			Resource: "secrets",
			Reaction: func(action client_testing.Action) (bool, runtime.Object, error) {
				return true, ret, nil
			},
		},
	}

	return clientset
}

func successfulUpdateRequestReactor(clientset kubernetes.Interface, ret runtime.Object) kubernetes.Interface {
	clientset.(*fake.Clientset).Fake.ReactionChain = []client_testing.Reactor{
		&client_testing.SimpleReactor{
			Verb:     "create",
			Resource: "secrets",
			Reaction: func(action client_testing.Action) (bool, runtime.Object, error) {
				return true, nil, errors.NewAlreadyExists(schema.GroupResource{Resource: "secrets"}, "test")
			},
		},
		&client_testing.SimpleReactor{
			Verb:     "update",
			Resource: "secrets",
			Reaction: func(action client_testing.Action) (bool, runtime.Object, error) {
				return true, ret, nil
			},
		},
	}

	return clientset
}

func failedUpdateRequestReactor(clientset kubernetes.Interface, _ runtime.Object) kubernetes.Interface {
	clientset.(*fake.Clientset).Fake.ReactionChain = []client_testing.Reactor{
		&client_testing.SimpleReactor{
			Verb:     "create",
			Resource: "secrets",
			Reaction: func(action client_testing.Action) (bool, runtime.Object, error) {
				return true, nil, errors.NewAlreadyExists(schema.GroupResource{Resource: "secrets"}, "test")
			},
		},
		&client_testing.SimpleReactor{
			Verb:     "update",
			Resource: "secrets",
			Reaction: func(action client_testing.Action) (bool, runtime.Object, error) {
				return true, nil, errors.NewInternalError(fmt.Errorf("test-err"))
			},
		},
	}

	return clientset
}
