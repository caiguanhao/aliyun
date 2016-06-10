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
	// if done channel is closed, AddJob() should be cancelled
	AddJob func(jobs *chan interface{}, done *chan interface{}, errs *chan error)
	// optional function to deal with errors from AddJob()
	OnAddJobError func(err *error)

	// required function to process job with job params from AddJob() and return results and errors,
	DoJob func(job *interface{}) (ret interface{}, err error)
	// optional function to deal with errors from DoJob()
	OnJobError func(err *error)
	// optional function to deal with results from DoJob()
	OnJobSuccess func(ret *interface{})
}

func (q Queue) doJob(job *interface{}) []interface{} {
	ret, err := q.DoJob(job)
	return []interface{}{ret, err}
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

	done := make(chan interface{})
	defer close(done)

	jobs := make(chan interface{})
	addJobErrs := make(chan error)
	go func() {
		defer close(jobs)
		defer close(addJobErrs)
		q.AddJob(&jobs, &done, &addJobErrs)
	}()

	rets := make(chan []interface{})
	var jobsWG sync.WaitGroup
	jobsWG.Add(q.Concurrency)
	for i := 0; i < q.Concurrency; i++ {
		go func() {
			defer jobsWG.Done()
			for job := range jobs {
				select {
				case rets <- q.doJob(&job):
				case <-done:
					return
				}
			}
		}()
	}
	go func() {
		jobsWG.Wait()
		close(rets)
	}()

	Parallel{
		func() {
			for err := range addJobErrs {
				if q.OnAddJobError != nil && err != nil {
					q.OnAddJobError(&err)
				}
			}
		},
		func() {
			for ret := range rets {
				if err, errExists := ret[1].(error); errExists && err != nil {
					if q.OnJobError != nil {
						q.OnJobError(&err)
					}
				} else if q.OnJobSuccess != nil {
					q.OnJobSuccess(&ret[0])
				}
			}
		},
	}.Run()
}
