package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var buffetStarttime time.Time

const SPOON = 20
const DISH = 37

type Guest struct {
	Id         int
	Starttime  time.Time
	Eattime    int
	Waittime   time.Duration
	Finishtime time.Duration
}
type semaphore chan struct{}

func BuffetStart(plist []Guest) {
	buffetStarttime = time.Now()
	fmt.Println("buffet Start time", buffetStarttime)

	wg := new(sync.WaitGroup)
	sem := make(semaphore, SPOON)
	wg.Add(len(plist))
	for i := 0; i <= len(plist)-1; i++ {
		go startEating(plist[i], wg, sem)
	}

	var totalwait time.Duration
	for _, p := range plist {
		totalwait = p.Waittime + time.Duration(totalwait)
	}
	wg.Wait()
	fmt.Println(" ")
	fmt.Println("buffet End time", time.Now())
	fmt.Printf("All guests have eaten : %s ", time.Now().Sub(buffetStarttime))
}
func startEating(g Guest, wg *sync.WaitGroup, sem semaphore) {
	sem <- struct{}{}
	g.Starttime = time.Now()
	g.Waittime = g.Starttime.Sub(buffetStarttime)
	time.Sleep(time.Duration(g.Eattime) * time.Second)
	g.Finishtime = time.Now().Sub(buffetStarttime)
	//fmt.Println("---------")
	fmt.Printf("user %d waited upto %s seconds and finished eating in %d", g.Id, g.Waittime, g.Eattime)
	fmt.Println(" ")
	<-sem
	wg.Done()
}

func getRandomNo() int {
	return rand.Intn(50-40) + 40
}

func main() {
	lst := []Guest{}
	for i := 1; i <= 100; i++ {
		g := Guest{
			Id:      i,
			Eattime: getRandomNo(),
		}
		lst = append(lst, g)
	}
	BuffetStart(lst)
}
