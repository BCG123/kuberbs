// Copyright © 2018 Aviv Laufer <aviv.laufer@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package utils

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	apps_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClient returns a k8s clientset to the request from inside of cluster
func GetClient() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Can not create kubernetes client: %v", err)
	}

	return clientset
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// GetClientOutOfCluster returns a k8s clientset to the request from outside of cluster
func GetClientOutOfCluster() kubernetes.Interface {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		logrus.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	return clientset
}

func isStrategySupported(strategy string) bool {
	if strings.ToLower(strategy) == "rollingupdate" {
		return true
	}
	return false
}

func DeploymentComplete(deployment *apps_v1.Deployment, newStatus *apps_v1.DeploymentStatus) bool {
	return newStatus.UpdatedReplicas == *(deployment.Spec.Replicas) &&
		newStatus.Replicas == *(deployment.Spec.Replicas) &&
		newStatus.AvailableReplicas == *(deployment.Spec.Replicas) &&
		newStatus.ObservedGeneration >= deployment.Generation
}

// GetObjectMetaData returns metadata of a given k8s object
func GetObjectMetaData(obj interface{}) meta_v1.ObjectMeta {

	var objectMeta meta_v1.ObjectMeta

	switch object := obj.(type) {
	case *apps_v1.Deployment:
		objectMeta = object.ObjectMeta
	case *api_v1.ReplicationController:
		objectMeta = object.ObjectMeta
	case *apps_v1.ReplicaSet:
		objectMeta = object.ObjectMeta
	case *apps_v1.DaemonSet:
		objectMeta = object.ObjectMeta
	case *api_v1.Service:
		objectMeta = object.ObjectMeta
	case *api_v1.Pod:
		objectMeta = object.ObjectMeta
	case *batch_v1.Job:
		objectMeta = object.ObjectMeta
	case *api_v1.PersistentVolume:
		objectMeta = object.ObjectMeta
	case *api_v1.Namespace:
		objectMeta = object.ObjectMeta
	case *api_v1.Secret:
		objectMeta = object.ObjectMeta
	case *ext_v1beta1.Ingress:
		objectMeta = object.ObjectMeta
	}
	return objectMeta
}
