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

func (t *Till) processCustomers() {

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
			fmt.Println("Customers: ", *customers)
			time.Sleep(500 * time.Millisecond) 
		}
}

//Assigning customers to queues
func customersToQueues(customers *[]Customer, tills *[]Till, running *bool) {
	time.Sleep(1 * time.Second)
	count := 0
	for *running {
			for i:= 0; i < 8; i++ {
				for (*tills)[i].checkLength() < 6 {
					(*tills)[i].queue <- (*customers)[0]
					fmt.Println("Assigning customers to till ", i, ": ", (*tills)[1].queue)
					*customers = append((*customers)[:0], (*customers)[0+1:]...)
					fmt.Println("Slice after assignment", *customers)
					time.Sleep(500 * time.Millisecond)
				}
			}
		count++
	}
}


//Creating the initial till slice and opening a few of them
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

	tillsOpen := (rand.Intn(9-1)+1)
	fmt.Println("Tills open at start of day: ", tillsOpen)

	for i:= 0; i < tillsOpen; i++ {
		(*tills)[i].opened = true
	}
	fmt.Println("Tills at start of day: ", *tills)
}

func main() {
	//Variables
	running := true
	var customers []Customer
	var tills []Till

	createTills(&tills)

	//Go routines
	go generateCustomers(&customers, &running)
	go customersToQueues(&customers, &tills, &running)

	time.Sleep(60 * time.Second) 
	running = false
}
