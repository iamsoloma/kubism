package api

import (
	"context"
	"encoding/json"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateNamespace struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
}

func (s *Server) ListNamespaces(w http.ResponseWriter, r *http.Request) {
	list, err := s.KubernetesClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonList, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonList)
}

func (s *Server) CreateNamespace(w http.ResponseWriter, r *http.Request) {
	req := CreateNamespace{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace, err := s.KubernetesClient.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       req.Name,
			Namespace:                  req.Name,
			Labels:                     req.Labels,
		},
		Spec:   corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{},
	}, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		http.Error(w, err.Error(), http.StatusFound)
		return
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	jsonNameSpace, err := json.Marshal(namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonNameSpace)
}
