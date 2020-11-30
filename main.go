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

var processed []Customer

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
	processedCustomers := make (chan Customer)
	for customer := range t.queue {
		for i:= 0; i < customer.numberOfItems; i++ {
			time.Sleep(50 * time.Millisecond) //change this to scanning speed
			//When processing last item
			if i == customer.numberOfItems - 1 {
				time.Sleep(50 * time.Millisecond) //change this to scanning speed
				//Remove customer from channel
				fmt.Println("processed ", customer)
				processed = append(processed, customer)
				processedCustomers <- customer
			}
		}
	}
}

//Create customers every 0.3 or 0.5 seconds
func generateCustomers(customers *[]Customer, running *bool) {
		rand.Seed(time.Now().UnixNano())
		//good weather or bad weather
		weather := (rand.Intn(2-1)+1)
		fmt.Println("Weather is: ", weather)
		count := 0
		for *running {
			customer := Customer {
				customerId : count,
				numberOfItems: (rand.Intn(200-1)+1),
			}
			*customers = append(*customers, customer)
			fmt.Println("Customers generated: ", *customers)
			count++
			if weather == 1 {
				time.Sleep(500 * time.Millisecond) 
			} else if weather == 2 {
				time.Sleep(500 * time.Millisecond) 
			}
		}
}

//Assigning customers to queues every 0.5 seconds
func customersToQueues(customers *[]Customer, tills *[]Till, lostCustomers *[]Customer, running *bool) {
	//sleep for 1 second so there's always customers generated before they're assigned
	time.Sleep(1 * time.Second)
	count := 0
	for *running {
			for i:= 0; i < 8; i++ {
				go (*tills)[i].processCustomers()
				for (*tills)[i].checkLength() < 6 && (*tills)[i].opened {
					//Adds customer to queue
					if len(*customers)!= 0 {
						(*tills)[i].queue <- (*customers)[0]
						fmt.Println("Assigning customers to till ", i+1, ": ", (*tills)[1].queue)
						//After added to queue, delete customer from slice
						*customers = append((*customers)[:0], (*customers)[0+1:]...)
						fmt.Println("Slice after assignment", *customers)
						time.Sleep(500 * time.Millisecond)
					}
				}
				if (*tills)[i].checkLength() == 6 {
					//check other till lengths
					fmt.Println("Customer lost: ", (*customers)[0])
					//Add to lost customers slice
					*lostCustomers = append(*lostCustomers, (*customers)[0])
					//Remove from original customers slice

					*customers = append((*customers)[:0], (*customers)[0+1:]...)
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
	var lostCustomers []Customer


	//Setting up tills
	createTills(&tills)

	//Go routines
	go generateCustomers(&customers, &running)
	go customersToQueues(&customers, &tills, &lostCustomers, &running)
	
	time.Sleep(60 * time.Second) 
	fmt.Println("TIMES UP!")
	fmt.Println("Lost customers: ", lostCustomers)
	fmt.Println("Processed customers: ", processed)
	running = false
}