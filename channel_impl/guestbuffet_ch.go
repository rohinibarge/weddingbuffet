package main

import (
	"fmt"
	"math/rand"
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
	// chanel of type pointer to guest objects of size 100
	eat := make(chan *Guest, len(plist))
	// chanel of type int of size 100
	guestIds := make(chan int, len(plist))

	buffetStarttime = time.Now()
	fmt.Println("buffet Start time", buffetStarttime)

	// we have 20 spoons, so only 20 guest can have dinner together at any time.
	// so 20 goroutines will be running at any given time until we all guest have their food.
	// startEating func called using goroutine, will wait for recive guest on channel eat
	for i := 0; i <= SPOON; i++ {
		go startEating(eat, guestIds)
	}

	// write on eat channel, looping over guest list and send user to eat one by one
	for j := 0; j < len(plist); j++ {
		eat <- &plist[j]
	}
	close(eat)

	// waiting for guests to finish their food.
	// reading channel guestids, where we will recive guest ids who have finished their food.
	for a := 1; a <= len(plist); a++ {
		<-guestIds
	}

	// calculate total time taken
	var totalwait time.Duration
	for _, p := range plist {
		totalwait = p.Waittime + time.Duration(totalwait)
	}

	fmt.Println(" ")
	fmt.Println("buffet End time", time.Now())
	fmt.Printf("All guests have eaten : %s ", time.Now().Sub(buffetStarttime))
}

// in this func which is called using goroutine, we will keep reading channel guest
// we will recieve pointer to object of guest who were sent to have food.
// func calculate wait time and completion time
// and finally write guest id on guestIds channel who have finished their food.
func startEating(guest <-chan *Guest, guestIds chan<- int) {
	for g := range guest {
		g.Starttime = time.Now()
		g.Waittime = g.Starttime.Sub(buffetStarttime)
		time.Sleep(time.Duration(g.Eattime) * time.Second)
		g.Finishtime = time.Now().Sub(buffetStarttime)

		fmt.Printf("user %d waited upto %s seconds and finished eating in %d", g.Id, g.Waittime, g.Eattime)
		fmt.Println(" ")
		guestIds <- g.Id
	}

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
