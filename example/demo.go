package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haobird/peajob"
)

func print() error {
	fmt.Println("执行队列任务")
	return nil
}

func info(v interface{}) error {
	fmt.Printf("大家好,我叫今年%d岁\n", v)
	text := fmt.Sprintf("测试的%d", v)
	return errors.New(text)
}

func main() {
	// 处理包
	pea := peajob.New(2).SetHandler(info)
	pea.SetCallback(func(err error, data interface{}) {
		fmt.Println("回调的执行", err, data)
	})
	// pea.Process(print)
	// pea.Run(1)
	// pea.Run(2)
	// pea.Run(3)
	// pea.Run(4)
	i := 1
	j := 2

	go func() {
		for {
			pea.RunWithTag(i, "test")
			i = i + 2
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		pea.RunWithTag(j, "test")
		j = j + 2
		time.Sleep(3 * time.Second)
	}

	// forever := make(chan bool)
	// fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever
	keepAlive()
}

func keepAlive() {
	//合建chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("启动")
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("退出信号", s)
}
