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

	guests := make(chan int, len(plist))
	finishStatus := make(chan int, len(plist))

	buffetStarttime = time.Now()
	fmt.Println("buffet Start time", buffetStarttime)

	for i := 0; i <= SPOON; i++ {
		go startEating(plist[i], guests, finishStatus)
	}

	for j := 1; j <= len(plist); j++ {
		guests <- j
	}
	close(guests)

	for a := 1; a <= len(plist); a++ {
		<-finishStatus
	}

	var totalwait time.Duration
	for _, p := range plist {
		totalwait = p.Waittime + time.Duration(totalwait)
	}

	fmt.Println(" ")
	fmt.Println("buffet End time", time.Now())
	fmt.Printf("All guests have eaten : %s ", time.Now().Sub(buffetStarttime))
}
func startEating(g Guest, guest <-chan int, finishStatus chan<- int) {
	for i := range guest {
		g.Starttime = time.Now()
		g.Waittime = g.Starttime.Sub(buffetStarttime)
		time.Sleep(time.Duration(g.Eattime) * time.Second)
		g.Finishtime = time.Now().Sub(buffetStarttime)

		fmt.Printf("user %d waited upto %s seconds and finished eating in %d", g.Id, g.Waittime, g.Eattime)
		fmt.Println(" ")
		finishStatus <- i
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
