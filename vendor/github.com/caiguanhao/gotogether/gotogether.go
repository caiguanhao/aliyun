// To run go code concurrently.
package gotogether

import "sync"

// Array of objects of any type.
type Enumerable []interface{}

// Run immediately and wait until all functions for each element return.
func (e Enumerable) Each(f func(item interface{})) {
	var wg sync.WaitGroup
	wg.Add(len(e))
	for _, item := range e {
		go func(item interface{}) {
			f(item)
			wg.Done()
		}(item)
	}
	wg.Wait()
}

func (e Enumerable) EachWithIndex(f func(item interface{}, i int)) {
	var wg sync.WaitGroup
	wg.Add(len(e))
	for i, item := range e {
		go func(item interface{}, i int) {
			f(item, i)
			wg.Done()
		}(item, i)
	}
	wg.Wait()
}

// Make Parallel{} for each element that can be Run() later.
func (e Enumerable) Parallel(f func(item interface{})) (p Parallel) {
	for _, item := range e {
		p = append(p, func(item interface{}) func() {
			return func() {
				f(item)
			}
		}(item))
	}
	return
}

func (e Enumerable) ParallelWithIndex(f func(item interface{}, i int)) (p Parallel) {
	for i, item := range e {
		p = append(p, func(item interface{}, i int) func() {
			return func() {
				f(item, i)
			}
		}(item, i))
	}
	return
}

// Convert Enumerable to Queue
func (e Enumerable) Queue(f func(item interface{})) (q Queue) {
	q.AddJob = func(jobs *chan interface{}) {
		for _, item := range e {
			*jobs <- item
		}
	}
	q.DoJob = func(job *interface{}) {
		f(*job)
		return
	}
	return
}

func (e Enumerable) QueueWithIndex(f func(item interface{}, i int)) (q Queue) {
	q.AddJob = func(jobs *chan interface{}) {
		for i, item := range e {
			*jobs <- []interface{}{item, i}
		}
	}
	q.DoJob = func(job *interface{}) {
		jobInfo := (*job).([]interface{})
		f(jobInfo[0], jobInfo[1].(int))
		return
	}
	return
}

// Array of functions to run concurrently.
type Parallel []func()

// Run immediately and wait until all functions return.
func (p Parallel) Run() {
	var wg sync.WaitGroup
	wg.Add(len(p))
	for _, f := range p {
		go func(f func()) {
			f()
			wg.Done()
		}(f)
	}
	wg.Wait()
}

// Run no more than number of concurrency of jobs together.
type Queue struct {
	// required max number of jobs to run together
	Concurrency int

	// required function to send job params that will be used in DoJob() to jobs channel, and error to errs channel
	AddJob func(jobs *chan interface{})

	// required function to process job with job params from AddJob() and return results and errors,
	DoJob func(job *interface{})
}

// Run immediately and wait until all jobs complete.
func (q Queue) Run() {
	if q.Concurrency < 1 {
		panic("concurrency must not be less than 1")
	}

	if q.AddJob == nil {
		panic("AddJob() must not be nil")
	}

	if q.DoJob == nil {
		panic("DoJob() must not be nil")
	}

	jobs := make(chan interface{})
	go func() {
		defer close(jobs)
		q.AddJob(&jobs)
	}()

	var wg sync.WaitGroup
	wg.Add(q.Concurrency)
	for i := 0; i < q.Concurrency; i++ {
		go func() {
			defer wg.Done()
			for job := range jobs {
				q.DoJob(&job)
			}
		}()
	}
	wg.Wait()
}

// Set concurrency
func (q Queue) WithConcurrency(concurrency int) Queue {
	q.Concurrency = concurrency
	return q
}
