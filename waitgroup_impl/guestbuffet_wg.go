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
	// sem is channel of size 20 i.e spoons count
	// so 20 guest can eat together at any given time.
	sem := make(semaphore, SPOON)
	// waitgroup created of size 100 i.e count of guests
	wg.Add(len(plist))
	// spin go routines for each guest
	for i := 0; i <= len(plist)-1; i++ {
		go startEating(&plist[i], wg, sem)
	}

	// wait until waitgroup counter is 0, i.e wait for all guests to finish their food
	wg.Wait()
	// calculate total wait time.
	var totalwait time.Duration
	for _, p := range plist {
		totalwait = p.Waittime + time.Duration(totalwait)
	}
	fmt.Println(" ")
	fmt.Println("buffet End time", time.Now())
	fmt.Printf("All guests have eaten : %s ", time.Now().Sub(buffetStarttime))
	fmt.Println("total wait time is ", totalwait)
}
func startEating(g *Guest, wg *sync.WaitGroup, sem semaphore) {
	// write to sem channel. we can write 20 times to this channel as size of channel buffer is 20, until we from it
	sem <- struct{}{}
	g.Starttime = time.Now()
	g.Waittime = g.Starttime.Sub(buffetStarttime)
	time.Sleep(time.Duration(g.Eattime) * time.Second)
	g.Finishtime = time.Now().Sub(buffetStarttime)
	fmt.Printf("user %d waited upto %s seconds and finished eating in %d", g.Id, g.Waittime, g.Eattime)
	fmt.Println(" ")
	// when guest finish his food, read from sem channel, so next guest i.e 21th guest can start eating.
	<-sem
	// decreament waitgroup counter
	wg.Done()
}

func getRandomNo() int {
	return rand.Intn(5-4) + 4
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
