package main

import (
    "flag"
    "fmt"
    "log"
    "math"
    "net/http"
    "sync"
    "time"

    kclient "k8s.io/client-go/kubernetes"
    kappslisters "k8s.io/client-go/listers/apps/v1"
    kinformers "k8s.io/client-go/informers"
    "k8s.io/client-go/rest"
    utilwait "k8s.io/apimachinery/pkg/util/wait"

    "bitbucket.org/bertimus9/systemstat"
)

const sleep = time.Duration(10) * time.Millisecond

var mutex = &sync.Mutex{}
var requests = 0
var requestDuration = flag.Duration("request-duration", 400 * time.Millisecond, "Duration of request processing")
var myDeployment = flag.String("deployment-name", "", "Name of deployment controlling this pod")
var myNamespace = flag.String("pod-namepace", "", "Namespace of this pod")
var targetMillicores = flag.Int("target-millicores", 200, "HPA target millicores")

type sampleServer struct {
	dpments kappslisters.DeploymentLister
	initialized bool
	first systemstat.ProcCPUSample
	uptime float64
}

func (b *sampleServer) reset() {
	mutex.Lock()
	defer mutex.Unlock()

	log.Print("Resetting CPU")
	b.initialized = false
}


func (b *sampleServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Process one request at a time
	// so that latency grows when system is overloaded
	// (aka. pretend we're python)
	mutex.Lock()
	defer mutex.Unlock()

	requests += 1
	// Gather initial CPU stats, so we can ignore anything before first request
	if !b.initialized {
		b.initialized = true
		b.first = systemstat.GetProcCPUSample()
		b.uptime = systemstat.GetUptime().Uptime
	}

	dep, err := b.dpments.Deployments(*myNamespace).Get(*myDeployment)
	if err != nil {
		log.Printf("failed to get parent deployment: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	replicas := dep.Status.Replicas

	millicoresPct := float64(*targetMillicores) / float64(10)
	if replicas < 2 {
		millicoresPct *= 2.1
	}

	start := time.Now()

	for time.Now().Sub(start) < *requestDuration {
		cpu := systemstat.GetProcCPUAverage(b.first, systemstat.GetProcCPUSample(), systemstat.GetUptime().Uptime - b.uptime)
		if cpu.TotalPct < millicoresPct {
			// Burn some CPU on pointless math.
			x := 0.0001
			for i := 0; i < 10000000; i++ {
				x += math.Sqrt(x)
			}
		} else {
			time.Sleep(sleep)
		}
	}
	fmt.Fprintf(w, "OK")
}

func metrics(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "# HELP http_requests_total The amount of requests served by the server in total\n# TYPE http_requests_total counter\nhttp_requests_total %v\n", requests)
}


func main() {
    flag.Parse()

    config, err := rest.InClusterConfig()
    if err != nil {
	log.Fatal("Failed to get client config")
    }

    clientSet := kclient.NewForConfigOrDie(config)
    informers := kinformers.NewFilteredSharedInformerFactory(clientSet, 20*time.Minute, *myNamespace, nil)

    server := &sampleServer{
	dpments: informers.Apps().V1().Deployments().Lister(),
    }

    go informers.Start(utilwait.NeverStop)

    go utilwait.Forever(server.reset, 5 * time.Minute)

    http.Handle("/app", server)
    http.HandleFunc("/metrics", metrics)
    log.Fatal(http.ListenAndServe(":1234", nil))
}
