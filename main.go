package main

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/docker/docker/pkg/units"
	"github.com/fsouza/go-dockerclient"
	"github.com/namsral/flag"
)

func main() {
	var dockerHost string
	var dockerTLSVerify bool
	var dockerCertPath string

	flag.StringVar(&dockerHost, "docker-host", "unix:///var/run/docker.sock", "address of Docker host")
	flag.BoolVar(&dockerTLSVerify, "docker-tls-verify", false, "use TLS client for Docker")
	flag.StringVar(&dockerCertPath, "docker-cert-path", "", "path to the cert.pem, key.pem, and ca.pem for authenticating to Docker")
	flag.Parse()

	client := dockerClient(dockerHost, dockerTLSVerify, dockerCertPath)

	s := newSelector()

	log.Println("Querying for stats...")

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		panic(err)
	}
	listenChan := make(chan *docker.APIEvents)
	client.AddEventListener(listenChan)
	go func() {
		for {
			event := <-listenChan

			if event.Status == "start" {
				w := newWatcher(event.ID[:12], client)
				s.Add(w)
				go w.Watch()
			}
		}
	}()

	for _, cont := range containers {
		w := newWatcher(cont.ID[:12], client)
		s.Add(w)
		go w.Watch()
	}

	log.Println("Waiting for stats...")
	s.Select()
}

type watcher struct {
	Name    string
	Stats   *stats
	Updates chan *docker.Stats
	client  *docker.Client
}

func newWatcher(name string, client *docker.Client) *watcher {
	return &watcher{
		Name:    name,
		Stats:   newStats(),
		Updates: make(chan *docker.Stats, 0),
		client:  client,
	}
}

func (w *watcher) Watch() {
	w.client.Stats(docker.StatsOptions{
		ID:    w.Name,
		Stats: w.Updates,
	})
}

type selector struct {
	cases    []reflect.SelectCase
	watchers []*watcher
	mu       sync.RWMutex
}

func newSelector() *selector {
	return &selector{
		cases:    make([]reflect.SelectCase, 0),
		watchers: make([]*watcher, 0),
	}
}

func (s *selector) Add(w *watcher) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cases = append(s.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(w.Updates)})
	s.watchers = append(s.watchers, w)
}

func (s *selector) remove(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cases = append(s.cases[:i], s.cases[i+1:]...)
	s.watchers = append(s.watchers[:i], s.watchers[i+1:]...)
}

func (s *selector) Select() {
	for {
		chosen, value, ok := reflect.Select(s.cases)
		s.mu.Lock()
		w := s.watchers[chosen]
		s.mu.Unlock()

		if !ok {
			s.remove(chosen)
			continue
		}

		switch value.Kind() {
		case reflect.Ptr:
			ds := value.Elem().Interface().(docker.Stats)
			w.Stats.Update(&ds)
			w.Display()
		}

	}
}

type stats struct {
	CPUPercentage    float64
	Memory           float64
	MemoryLimit      float64
	MemoryPercentage float64
	NetworkRx        float64
	NetworkTx        float64

	previousCPUUsage       float64
	previousSystemCPUUsage float64
	cpuUsage               float64
	numberCPUs             float64
	systemCPUUsage         float64
	mu                     sync.RWMutex
}

func newStats() *stats {
	return &stats{
		cpuUsage:       0.0,
		systemCPUUsage: 0.0,
		CPUPercentage:  0.0,
	}
}

func (s *stats) Update(d *docker.Stats) {
	s.mu.Lock()
	s.previousCPUUsage = s.cpuUsage
	s.previousSystemCPUUsage = s.systemCPUUsage
	s.numberCPUs = float64(len(d.CPUStats.CPUUsage.PercpuUsage))
	s.cpuUsage = float64(d.CPUStats.CPUUsage.TotalUsage)
	s.systemCPUUsage = float64(d.CPUStats.SystemCPUUsage)
	s.calculateCPUPercentage()
	s.Memory = float64(d.MemoryStats.Usage)
	s.MemoryLimit = float64(d.MemoryStats.Limit)
	s.calculateMemoryPercentage()
	s.NetworkRx = float64(d.Network.RxBytes)
	s.NetworkTx = float64(d.Network.TxBytes)
	s.mu.Unlock()
}

func (s *stats) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return fmt.Sprintf("%.2f%%\t%s/%s\t%.2f%%\t%s/%s",
		s.CPUPercentage,
		units.HumanSize(s.Memory), units.HumanSize(s.MemoryLimit),
		s.MemoryPercentage,
		units.HumanSize(s.NetworkRx), units.HumanSize(s.NetworkTx))
}

func (w *watcher) Display() {
	fmt.Printf("%s\t%s\n", w.Name, w.Stats.String())
}

func (s *stats) calculateMemoryPercentage() {
	s.MemoryPercentage = s.Memory / s.MemoryLimit * 100.0
}

func (s *stats) calculateCPUPercentage() {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(s.cpuUsage - s.previousCPUUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(s.systemCPUUsage - s.previousSystemCPUUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * s.numberCPUs * 100.0
	}

	s.CPUPercentage = cpuPercent
}

func dockerClient(host string, tls bool, certPath string) *docker.Client {
	var client *docker.Client
	var err error

	if tls {
		cert := fmt.Sprintf("%s/cert.pem", certPath)
		key := fmt.Sprintf("%s/key.pem", certPath)
		ca := fmt.Sprintf("%s/ca.pem", certPath)
		client, err = docker.NewTLSClient(host, cert, key, ca)
	} else {
		client, err = docker.NewClient(host)
	}

	if err != nil {
		log.Fatal(err)
	}

	return client
}
