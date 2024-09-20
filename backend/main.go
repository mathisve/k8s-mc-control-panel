package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

var NAMESPACE = "mc"
var DEPLOYMENT_NAME = "mc"
var IMAGE_NAME = "minecraft-server"

var deploymentClient appsv1.DeploymentInterface
var podsClient v1.PodInterface
var clientset *kubernetes.Clientset

func main() {
	fmt.Println("running backend control panel")

	config, err := rest.InClusterConfig()
	if err != nil {
		// If we're not in a cluster, use kubeconfig
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentClient = clientset.AppsV1().Deployments(NAMESPACE)
	podsClient = clientset.CoreV1().Pods(NAMESPACE)

	handleRequests()

}

func start(w http.ResponseWriter, r *http.Request) {
	log.Println("starting server")
	err := UpdateDeployment(DEPLOYMENT_NAME, 1)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func stop(w http.ResponseWriter, r *http.Request) {
	log.Println("stopping server")
	err := UpdateDeployment(DEPLOYMENT_NAME, 0)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func status(w http.ResponseWriter, r *http.Request) {
	status, err := GetDeploymentStatus(DEPLOYMENT_NAME)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Write([]byte(status))
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := getPodLogs()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, logs)
}

func getPodLogs() (logs string, err error) {
	pods, err := podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return logs, err
	}

	var podName string
	var containerName string

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if strings.Contains(container.Image, IMAGE_NAME) {
				podName = pod.Name
				containerName = container.Name
				break
			}
		}
	}

	if podName == "" {
		return logs, nil
	}

	log.Println("retrieving pod logs")
	req := clientset.CoreV1().Pods(NAMESPACE).GetLogs(podName, &corev1.PodLogOptions{Container: containerName})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return logs, err
	}

	defer podLogs.Close()

	bytes, err := io.ReadAll(podLogs)
	if err != nil {
		return logs, err
	}

	return string(bytes), err
}

func UpdateDeployment(name string, size int) (err error) {
	if size > 1 || size < 0 {
		return fmt.Errorf("size needs to be either 0 or 1")
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentClient.Get(context.TODO(), name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of deployment: %v", getErr))
		}

		log.Printf("updating deployment replicas: %v to %d\n", name, size)
		result.Spec.Replicas = int32Ptr(int32(size))
		_, updateErr := deploymentClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	return retryErr
}

func GetDeploymentStatus(name string) (status string, err error) {
	result, getErr := deploymentClient.Get(context.TODO(), name, metav1.GetOptions{})
	if getErr != nil {
		return "", getErr
	}

	var s string
	if result.Status.Replicas == 0 {
		s = "stopped"
	} else if result.Status.Replicas == 1 {
		s = "started"
	}

	return s, nil
}

func handleRequests() {

	mux := http.NewServeMux()

	mux.HandleFunc("/start", start)
	mux.HandleFunc("/stop", stop)
	mux.HandleFunc("/status", status)
	mux.HandleFunc("/logs", getLogs)
	log.Println("serving :80")
	log.Fatal(http.ListenAndServe(":80", enableCors(mux)))
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-headers", "Content-Type")

		next.ServeHTTP(w, r)
	})
}

func int32Ptr(i int32) *int32 { return &i }
