package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	labelname *string
	labelvalue *string
	listenaddr *string
	kubeconfig *string
)
// initilization function to initialize flags
func init() {
	labelname = flag.String("label-name","","name of pod label to filter on")
	labelvalue = flag.String("label-value","","value corresponding to pod label name")
	listenaddr = flag.String("metrics-listen-addr","8080","name of pod label to filter on")
	// check if home directory is not empty then set kubeconfig in home directory
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "name of pod label to filter on")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "name of pod label to filter on")
	}

}

func main() {
	flag.Parse()
	updatePodCounts()
	http.Handle("/metrics", promhttp.Handler())
	log.Print("Starting prometheus server on localhost:" + *listenaddr)
	log.Fatal(http.ListenAndServe(":" + *listenaddr, nil))
	
}

// update pod counts every 10 seconds
func updatePodCounts() {
	// creating prometheus gauge object
	opsQueued := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:"pod_count",
	},[]string{"label_name","label_value","phase"})
	prometheus.MustRegister(opsQueued)
	// creating the context for the client command using kubeconfig path
	config, err := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Setting kubeconfig to set kubernetes environment")
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// map to have a count of pods for each pod lifecycle like Running,
	// CrashLoopBackoff, ImagePullBackOff. key is the lifecycle and value is
	// the count
	countmap := make(map[string]float64)
	// running this in a background go routine so that the rest api endpoint
	// can run in the main go routine
	go func() {
		for {
			// getting all the pods that match label selector "labelname=labelvalue"
			// which is gotten from the command line
			log.Print("getting pods from kubernetes cluster")
			pods , err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
				LabelSelector: *labelname + "=" + *labelvalue,
			})
			if err != nil {
				panic(err.Error())
			}
			// iterate through pods and adding it to the map
			// could have set the prometheus metric here but then I thought we would 
			// be reassiginging it everytime so I didn't put it here
			log.Print("Getting pod status and counts for statuses")
			for _, pod := range pods.Items {
				if val, ok := countmap[string(pod.Status.Phase)]; ok {
					countmap[string(pod.Status.Phase)] = val + 1
				} else {
					countmap[string(pod.Status.Phase)] = 1
				}
			}
			log.Print("Getting pod status and counts for statuses")
			for key, value := range countmap {
				// creating metrics for prometheus for each of the different pod lifecycle type
				log.Printf("adding prometheus metric with label %s, label value %s, and status %s",*labelname,*labelvalue,key)
				opsQueued.WithLabelValues(*labelname,*labelvalue,key).Set(value)
			}
			log.Print("updating prometheus pod count gauge")
			// checking the status of pods every 10 seconds
			time.Sleep(10 * time.Second)
			countmap = make(map[string]float64)
		}
	}()
}