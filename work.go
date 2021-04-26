package peajob

import "fmt"

// 执行策略，处理 N次重试机制
type worker struct {
	times int
	name  string
}

func newWorker(times int, name string) *worker {
	return &worker{
		times: times,
		name:  name,
	}
}

// 执行当前数据
func (w *worker) execute(data interface{}, fun handler) error {
	var err error
	var text string
	for i := 0; i < w.times; i++ {
		text = fmt.Sprintf("worker: %s的第 %d 次", w.name, i)
		fmt.Println(text)
		err = fun(data)

		if err == nil {
			return err
		}
	}
	return err
}

// 实际执行 方法
func (w *worker) process() error {
	return nil
}
