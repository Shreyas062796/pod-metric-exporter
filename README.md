# Pod Metric Exporter

The pod metric exporter is a kubernetes controller that periodically gathers the pod count for different pod lifecycles given a pod label and value as a parameter every 10 seconds. It also exposes the prometheus http api's /metrics route so that you will be able to call the api and be able to see pod counts for each pod lifecycle. We would be able to see what the status of
pods with specific labels are at anytime.

## Usage

Given that you have a working Kubernetes environment

```
kubectl get pods --all-namespaces
NAMESPACE         NAME                               READY   STATUS    RESTARTS   AGE
demo-production   demo-app                           1/1     Running   0          22h
demo-production   demo-db                            1/1     Running   0          22h
demo-staging      demo-app                           1/1     Running   0          22h
kube-system       coredns-74ff55c5b-dz668            1/1     Running   0          46h
kube-system       etcd-minikube                      1/1     Running   0          46h
kube-system       kube-apiserver-minikube            1/1     Running   0          46h
kube-system       kube-controller-manager-minikube   1/1     Running   0          46h
kube-system       kube-proxy-zjnjp                   1/1     Running   0          46h
kube-system       kube-scheduler-minikube            1/1     Running   0          46h
kube-system       storage-provisioner                1/1     Running   1          46h
```

Run command go run main.go -label-name <label_name> -label-value <label_value> -metrics-listen-addr <prometheus_server_port> -kubeconfig <kubeconfig_path>

Example output below
```
go run main.go -label-name app -label-value demo-app -metrics-listen-addr 8080
2021/05/06 20:48:40 Setting kubeconfig to set kubernetes environment
2021/05/06 20:48:40 getting pods from kubernetes cluster
2021/05/06 20:48:40 Starting prometheus server on localhost:8080
2021/05/06 20:48:40 Getting pod status and counts for statuses
2021/05/06 20:48:40 Getting pod status and counts for statuses
2021/05/06 20:48:40 adding prometheus metric with label app, label value demo-app, and status Running
2021/05/06 20:48:40 updating prometheus pod count gauge
2021/05/06 20:48:50 getting pods from kubernetes cluster
2021/05/06 20:48:50 Getting pod status and counts for statuses
2021/05/06 20:48:50 Getting pod status and counts for statuses
2021/05/06 20:48:50 adding prometheus metric with label app, label value demo-app, and status Running
2021/05/06 20:48:50 updating prometheus pod count gauge
```

You can also call the prometheus metrics endpoint to get pod count by doing
curl -s http://localhost:<port>/metrics | grep pod_count

Example Below running on port 8080

```
curl -s http://localhost:8080/metrics | grep pod_count
# HELP pod_count
# TYPE pod_count gauge
pod_count{label_name="app",label_value="demo-app",phase="Running"} 2
```