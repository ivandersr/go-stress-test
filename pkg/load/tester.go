package load

import (
	"net/http"
	"sync"
	"time"
)

type Report struct {
	ExecutionTime           time.Duration
	TotalRequests           int
	SuccessfulRequests      int
	FailedRequests          int
	ErrorStatusDistribution map[int]int

	mu sync.Mutex
}

type Tester struct {
	URL         string
	Requests    int
	Concurrency int
}

func NewTester(url string, requests int, concurrency int) *Tester {
	return &Tester{
		URL:         url,
		Requests:    requests,
		Concurrency: concurrency,
	}
}

func (t *Tester) Run() *Report {
	start := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}
	report := &Report{
		ErrorStatusDistribution: make(map[int]int),
		TotalRequests:           t.Requests,
	}

	requestsChannel := make(chan struct{}, t.Requests)
	resultsChannel := make(chan int, t.Requests)

	var wg sync.WaitGroup

	for i := 0; i < t.Concurrency; i++ {
		wg.Add(1)
		go t.threadWorker(client, requestsChannel, resultsChannel, &wg)
	}

	for i := 0; i < t.Requests; i++ {
		requestsChannel <- struct{}{}
	}

	close(requestsChannel)
	wg.Wait()
	close(resultsChannel)

	for status := range resultsChannel {
		report.processResult(status)
	}

	report.ExecutionTime = time.Since(start)
	return report
}

func (t *Tester) threadWorker(client *http.Client, requestsChannel chan struct{}, resultsChannel chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for range requestsChannel {
		res, _ := client.Get(t.URL)
		resultsChannel <- res.StatusCode
		res.Body.Close()
	}
}

func (r *Report) processResult(status int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if status >= 200 && status < 300 {
		r.SuccessfulRequests++
	} else {
		r.FailedRequests++
		r.ErrorStatusDistribution[status]++
	}
}
