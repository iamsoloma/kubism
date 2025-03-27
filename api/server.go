package api

import (
	"log/slog"
	"net/http"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	ListenAddr string
}

type Server struct {
	*Config
	Started time.Time

	KubernetesClient *kubernetes.Clientset
}

func NewServer(config Config) (*Server, error) {
	// creates the in-cluster config
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	return &Server{
		Config:           &config,
		KubernetesClient: clientset,
	}, nil
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /namespaces", s.ListNamespaces)
	mux.HandleFunc("POST /namespaces", s.CreateNamespace)
	server := http.Server{
		Addr:    s.Config.ListenAddr,
		Handler: mux,
	}

	s.Started = time.Now().UTC()

	slog.Info("api is running", "address", s.Config.ListenAddr)
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("API stoped", "error", err)
	}

}
