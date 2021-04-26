package peajob

import (
	"fmt"
	"log"
)

// 主逻辑： 建立模块，分配策略，接口

// handle 处理函数
type handler func(interface{}) error

// callback 回调函数
type callback func(error, interface{})

// job 任务
type job struct {
	id    string
	retry int
	data  interface{}
}

type Pea struct {
	ch      chan job
	times   int
	handler handler
	caller  callback
	// The logger used for this table.
	logger *log.Logger
}

func New(maxGoroutine int, fun handler) *Pea {
	p := &Pea{
		ch:      make(chan job, 2),
		times:   3,
		handler: fun,
	}
	for i := 0; i < maxGoroutine; i++ {
		name := fmt.Sprintf("worker %d", i)
		worker := newWorker(1, name)
		go p._process(worker)
	}

	return p
}

func (p *Pea) Run(data interface{}) {
	job := job{
		id:    "xx",
		retry: 0,
		data:  data,
	}
	p.ch <- job
}

func (p *Pea) _run(job job) {
	p.ch <- job
}

func (p *Pea) _process(w *worker) {
	// 处理循环多次的往复逻辑
	for job := range p.ch {
		fmt.Println("读取到队列数据", job)
		err := w.execute(job.data, p.handler)
		job.retry++
		fmt.Println("执行了当前job的第次：", job.retry, job.data)

		if job.retry < p.times && err != nil {
			// 重新放入队列
			fmt.Println("重新放入队列", job)
			p._run(job)
		} else if p.caller != nil {
			p.caller(err, job.data)
		}

	}
}

func (p *Pea) SetCallback(fun callback) {
	p.caller = fun
}

// SetLogger sets the logger to be used by this cache table.
func (p *Pea) SetLogger(logger *log.Logger) {
	p.logger = logger
}

// Internal logging method for convenience.
func (p *Pea) log(v ...interface{}) {
	if p.logger == nil {
		return
	}

	p.logger.Println(v...)
}
