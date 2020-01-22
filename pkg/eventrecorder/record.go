package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/api/core/v1"
	clientset "github.com/openebs/maya/pkg/client/generated/clientset/versioned"
	openebsScheme "github.com/openebs/maya/pkg/client/generated/clientset/versioned/scheme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Path for kube config")
)

type Record interface{
	RecordOnPVC()
}

type Recorder struct {
	Recorder record.EventRecorder
	clients *clients
}

type clients struct {
	kubeClient kubernetes.Interface
	openebsClient clientset.Interface
}

func main(){
	fmt.Println("Hello")
	NewRecorderWithComponentName("Maya").RecordOnSPC("overprovisioning-disabled-sparse-pool")
	//NewRecorderWithComponentName("Maya").RecordOnPVCWithPVCObject("cspc-operator-79958787b5-brrm9")
}
func NewRecorderWithComponentName(componentName string) *Recorder {
	c,err:=GetKubeClient()
	if err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err = openebsScheme.AddToScheme(scheme.Scheme)
	if err != nil {
		klog.Errorf("failed to add to scheme: error {%v}", err)
	}
	broadcaster := record.NewBroadcaster()
	broadcaster.StartLogging(klog.Infof)
	broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: c.kubeClient.CoreV1().Events("")})
	eventRecorder := broadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: componentName})
	return &Recorder{
		Recorder:eventRecorder,
		clients:c,
	}
}

func (r *Recorder)RecordOnPVCWithPVCObject(name string)  {

	pvc,err:=r.clients.kubeClient.CoreV1().Pods("openebs").Get(name,v12.GetOptions{})
	fmt.Println(pvc.UID)
	if err!=nil{
		fmt.Println("Error in getting PVC",err.Error())
	}
	klog.Info("cStorPool Added event")
	r.Recorder.Event(pvc,"Warning","Creation SUcc","This shall pass")
}

func (r *Recorder)RecordOnSPC(name string)  {

	pvc,err:=r.clients.openebsClient.OpenebsV1alpha1().StoragePoolClaims().Get(name,v12.GetOptions{})
	if err!=nil{
		fmt.Println("Error in getting PVC",err.Error())
	}

	r.Recorder.Event(pvc,v1.EventTypeWarning, "Update","This shall pass")
}

func GetKubeClient()(*clients,error){
	klog.InitFlags(nil)
	err := flag.Set("logtostderr", "true")
	if err != nil {
		return nil,errors.Wrap(err, "failed to set logtostderr flag")
	}
	flag.Parse()

	cfg, err := getClusterConfig(*kubeconfig)
	if err != nil {
		return nil,errors.Wrap(err, "error building kubeconfig")
	}
	// Building OpenEBS Clientset
	openebsClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil,errors.Wrap(err, "error building openebs clientset")
	}
	// Building Kubernetes Clientset
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil,errors.Wrap(err, "error building kubernetes clientset")
	}
	return &clients{kubeClient:kubeClient,openebsClient:openebsClient},nil
}

// GetClusterConfig return the config for k8s.
func getClusterConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	klog.V(2).Info("Kubeconfig flag is empty")
	return rest.InClusterConfig()
}