package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Customer struct {
	customerId 	  int
	numberOfItems int
}

type Till struct {
	tillId       int
	scannerSpeed float64
	queue        chan Customer
	opened       bool
}

	
func (t *Till) checkLength() int {
    return len(t.queue)
}

//worker is a checkout, jobs are customers
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        //fmt.Println("worker", id, "started  job", j)
        time.Sleep(time.Second)
        //fmt.Println("worker", id, "finished job", j)
        results <- j * 2
    }
}

//Create customers at random intervals
func generateCustomers(customers *[]Customer, running *bool) {
		rand.Seed(time.Now().UnixNano())
		count := 0
		for *running {
			customer := Customer {
				customerId : count,
				numberOfItems: (rand.Intn(100-1)+1),
			}
			*customers = append(*customers, customer)
			count++
			fmt.Println(*customers)
			time.Sleep(500 * time.Millisecond) 
		}
}

func createTills(tills *[]Till) {
	rand.Seed(time.Now().UnixNano())

	for i:= 0; i < 8; i++ {
		till := Till {
			tillId: i+1,
			scannerSpeed: float64(rand.Intn(4-1)+1),
			queue: make (chan Customer, 6),
			opened: false,
		}
		*tills = append(*tills, till)
	}
	fmt.Println(*tills)

	tillsOpen := (rand.Intn(9-1)+1)
	fmt.Println("Tills open at start of day: ", tillsOpen)

	for i:= 0; i < tillsOpen; i++ {
		(*tills)[i].opened = true
	}
	fmt.Println(*tills)
}

func main() {
	//Variables
	running := true
	var customers []Customer
	var tills []Till

	createTills(&tills)

	//Go routines
	go generateCustomers(&customers, &running)

	time.Sleep(60 * time.Second) 
	running = false









	//Go routine to generate customers into a slice
	const numJobs = 5
	
	//Put customers into a channel
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= numJobs; a++ {
        <-results
    }
}