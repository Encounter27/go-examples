package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type I interface {
	worker()
	start()
}

type Node struct {
	wg sync.WaitGroup
	ch chan int
}

var node Node

func CreateNode() *Node {
	node.ch = make(chan int, 1)

	return &node
}

const (
	Stopped = 0
	Paused  = 1
	Running = 2
)

func (n *Node) controller(workers []chan int) {
	n.wg.Add(1)
	defer n.wg.Done()

	n.ch <- Running
	n.ch <- Paused
	n.ch <- Running
	n.ch <- Stopped
}

func (n *Node) worker() {
	n.wg.Add(1)
	fmt.Println("started worker")
	defer fmt.Println("end worker")
	defer n.wg.Done()

	state := Paused
	for i := 0; i < 20; i++ {
		select {
		case state = <-n.ch:
			switch state {
			case Stopped:
				fmt.Printf("Worker: Stopped\n")
				return
			case Running:
				for i := 0; i < 5; i++ {
					fmt.Printf("Worker: Running\n")
				}
				n.ch <- Paused
			case Paused:
				fmt.Printf("Worker: Paused\n")
			}
		default:
			runtime.Gosched()
			fmt.Printf("Worker: Default\n")
			// if state == Paused {
			// 	break
			// }
		}
		time.Sleep(time.Second * 1)
	}
}

func (n *Node) start() {
	fmt.Println("started start")
	defer fmt.Println("end start")

	go n.worker()

	go func() {
		//for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 2)
		n.ch <- Running
		time.Sleep(time.Second * 2)
		n.ch <- Running
		time.Sleep(time.Second * 2)
		n.ch <- Stopped
		//}
	}()
}

func main() {
	n := CreateNode()

	go n.start()

	//n.wg.Wait()
	time.Sleep(time.Second * 20)
}
