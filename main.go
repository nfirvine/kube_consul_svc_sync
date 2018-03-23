package main

// from https://github.com/kubernetes/sample-controller/blob/master/main.go

import (
	"flag"
	"text/template"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"k8s.io/sample-controller/pkg/signals"
)

const consulNodeNameTmplStrDefault = `{{ .TargetNamespace }}-{{ .TargetName }}`
const consulServiceNameTmplStrDefault = `{{ .EndpointsNamespace }}-{{ .EndpointsName }}-{{ .PortName }}`

var (
	masterURL                string
	kubeconfig               string
	namespace                string
	consulNodeNameTmplStr    string
	consulServiceNameTmplStr string
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	var kubeInformerFactory kubeinformers.SharedInformerFactory
	if namespace != "" {
		kubeInformerFactory = kubeinformers.NewFilteredSharedInformerFactory(
			kubeClient, time.Second*30, namespace, func(_ *v1.ListOptions) {})
	} else {
		kubeInformerFactory = kubeinformers.NewSharedInformerFactory(
			kubeClient, time.Second*30)
	}

	consulNodeNameTmpl := template.New("consulNodeNameTmpl")
	if consulNodeNameTmplStr == "" {
		consulNodeNameTmplStr = consulNodeNameTmplStrDefault
	}
	consulNodeNameTmpl, err = consulNodeNameTmpl.Parse(consulNodeNameTmplStr)
	if err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}

	consulServiceNameTmpl := template.New("consulServiceNameTmpl")
	if consulNodeNameTmplStr == "" {
		consulServiceNameTmplStr = consulServiceNameTmplStrDefault
	}
	consulServiceNameTmpl, err = consulServiceNameTmpl.Parse(consulNodeNameTmplStr)
	if err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}

	controller := NewController(kubeClient, kubeInformerFactory, consulNodeNameTmpl, consulServiceNameTmpl)

	go kubeInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&namespace, "namespace", "", "Only watch this namespace; default: all")
	flag.StringVar(&consulNodeNameTmplStr, "nodenametmpl", consulNodeNameTmplStrDefault, "consul node name template (Go's text/template)")
	flag.StringVar(&consulServiceNameTmplStr, "servicenametmpl", consulServiceNameTmplStrDefault, "consul service name template (Go's text/template)")

}
