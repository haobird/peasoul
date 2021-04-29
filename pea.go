package peajob

import (
	"fmt"
	"log"
	"time"
)

// 主逻辑： 建立模块，分配策略，接口

// handle 处理函数
type handler func(interface{}) error

// callback 回调函数
type callback func(error, interface{})

// job 任务
type Job struct {
	id      string
	tag     string
	retry   int
	data    interface{}
	handler handler
}

type Pea struct {
	ch      chan Job
	cache   map[string]string
	times   int // 每个任务重复执行次数
	retry   int // 每个任务失败后，重复放队列次数
	handler handler
	caller  callback
	// The logger used for this table.
	logger *log.Logger
}

func New(maxGoroutine int) *Pea {
	fun := func(data interface{}) error {
		fmt.Println(data)
		return nil
	}

	p := &Pea{
		ch:      make(chan Job, 2),
		cache:   make(map[string]string),
		times:   1,
		retry:   2,
		handler: fun,
	}

	go p._start(maxGoroutine)

	return p
}

func (p *Pea) SetHandler(fun handler) *Pea {
	p.handler = fun
	return p
}

func (p *Pea) SetCallback(fun callback) *Pea {
	p.caller = fun
	return p
}

// SetLogger sets the logger to be used by this cache table.
func (p *Pea) SetLogger(logger *log.Logger) *Pea {
	p.logger = logger
	return p
}

// Internal logging method for convenience.
func (p *Pea) log(v ...interface{}) {
	if p.logger == nil {
		return
	}

	p.logger.Println(v...)
}

func (p *Pea) Run(data interface{}) {
	p.RunWithTag(data, "")
}

func (p *Pea) RunWithTag(data interface{}, tag string) {
	id, _ := GenShortID()
	job := Job{
		id:    id,
		tag:   tag,
		retry: 0,
		data:  data,
	}
	p._run(job)
}

func (p *Pea) RunJob(job Job) {
	// 处理job_id
	// 处理job_tag
	p._run(job)
}

func (p *Pea) _run(job Job) {
	p.ch <- job
}

//
func (p *Pea) _start(maxGoroutine int) {
	sem := make(chan int, maxGoroutine)
	for job := range p.ch {
		// 判断是否有标志位的逻辑
		tag := job.tag
		if tag != "" {
			_, ok := p.cache[tag]
			if ok {
				// 重新放入队列
				fmt.Println("因为存在同一个tag处理gorouting，故重新放入队列", job)
				p._run(job)
				continue
			}
			p.cache[tag] = tag
		}

		sem <- 1
		go func(job Job) {
			p._process_job(job)
			delete(p.cache, tag)
			<-sem
		}(job)
	}

}

func (p *Pea) _process_job(job Job) {
	fmt.Println("======读取到队列数据", job, time.Now().Format("2006-01-02 15:04:05"))
	handelr := job.handler
	if job.handler == nil {
		handelr = p.handler
	}
	err := p._execute(job.data, handelr)
	job.retry++
	fmt.Printf("执行了当前job[%s]的第%d次：\n", job.id, job.retry)
	if job.retry < p.retry && err != nil {
		// 重新放入队列
		fmt.Println("重新放入队列", job)
		p._run(job)
	} else if p.caller != nil {
		p.caller(err, job.data)
	}

}

// 执行当前数据
func (p *Pea) _execute(data interface{}, fun handler) error {
	var err error
	for i := 0; i < p.times; i++ {
		// text = fmt.Sprintf("worker的第 %d 次", i)
		// fmt.Println(text)
		err = fun(data)

		if err == nil {
			return err
		}
	}
	return err
}
