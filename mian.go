package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	create := false
	cm := &corev1.ConfigMap{}
	cm, err = clientset.CoreV1().
		ConfigMaps(os.Getenv("namespace")).
		Get(context.TODO(), os.Getenv("name"), metav1.GetOptions{})
	if errors.IsNotFound(err) {
		create = true
		err = nil
		cm.Namespace = os.Getenv("namespace")
		cm.Name = os.Getenv("name")
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	files := make([]string, 0)
	if os.Getenv("files") != "" {
		//get file content
		err = yaml.Unmarshal([]byte(os.Getenv("files")), &files)
		if err != nil {
			err = nil
			err = json.Unmarshal([]byte(os.Getenv("files")), &files)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	if os.Getenv("dir") != "" {
		filepathNames, err := filepath.Glob(filepath.Join(os.Getenv("dir"), "*"))
		if err != nil {
			log.Fatalln(err)
		}
		for i := range filepathNames {
			var exist bool
			for j := range files {
				if files[j] == filepathNames[i] {
					exist = true
				}
			}
			if exist {
				continue
			}
			files = append(files, filepathNames[i])
		}
	}
	for i := range files {
		rb, err := ioutil.ReadFile(files[i])
		if err != nil {
			log.Fatalln(err)
		}
		cm.Data[filepath.Base(files[i])] = string(rb)
	}
	if os.Getenv("data") != "" || os.Getenv("data") != "{}" {
		params := make(map[string]string)
		err = yaml.Unmarshal([]byte(os.Getenv("data")), &params)
		if err != nil {
			err = nil
			err = json.Unmarshal([]byte(os.Getenv("data")), &params)
			if err != nil {
				log.Fatalln(err)
			}
		}
		for k := range params {
			cm.Data[k] = params[k]
		}

	}

	if create {
		err = nil
		cm = &corev1.ConfigMap{}
		cm.Namespace = os.Getenv("namespace")
		cm.Name = os.Getenv("name")
		_, err = clientset.CoreV1().ConfigMaps(cm.Namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("create configmap(%s) success\n", cm.Name)
		return
	}
	cm.Namespace = os.Getenv("namespace")
	cm.Name = os.Getenv("name")
	_, err = clientset.CoreV1().ConfigMaps(cm.Namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("update configmap(%s) success\n", cm.Name)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
