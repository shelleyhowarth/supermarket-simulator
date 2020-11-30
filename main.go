package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
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


func (t *Till) processCustomers(processed *[]Customer) {
	processedCustomers := make (chan Customer)
	for customer := range t.queue {
		for i:= 0; i < customer.numberOfItems; i++ {
			time.Sleep(time.Duration((t.scannerSpeed*10)) * time.Millisecond) //scanning speed
			fmt.Println("Scanning speed: ", t.scannerSpeed*10 )
			//When processing last item
			if i == customer.numberOfItems - 1 {
				time.Sleep(time.Duration((t.scannerSpeed*10)) * time.Millisecond) //scanning speed

				//Remove customer from channel
				fmt.Println("processed ", customer)
				*processed = append(*processed, customer)
				processedCustomers <- customer
			}
		}
	}
}

//Create customers every 0.3 or 0.5 seconds
func generateCustomers(customers *[]Customer, running *bool, weather *int, allCust *[]Customer) {
		rand.Seed(time.Now().UnixNano())
		//good weather or bad weather
		//weather := (rand.Intn(2-1)+1)
		fmt.Println("Weather is: ", *weather)
		count := 0
		for *running {
			customer := Customer {
				customerId : count,
				numberOfItems: (rand.Intn(200-1)+1),
			}
			*customers = append(*customers, customer)
			*allCust = append(*allCust, customer)
			fmt.Println("Customers generated: ", *customers)
			count++
			if *weather == 1 {
				time.Sleep(200 * time.Millisecond) 
			} else if *weather == 2 {
				time.Sleep(500 * time.Millisecond) 
			}
		}
}

func shortestQueue(tills *[]Till) (t Till){
	//min := values[0]
	min := len((*tills)[0].queue)
	shortestTill := (*tills)[0]
		for i:= 0; i < len(*tills); i++ {
			if (len((*tills)[i].queue) < min && (*tills)[i].opened) {
				min = len((*tills)[i].queue)
				shortestTill = (*tills)[i]
			}

		}
	return shortestTill
}

//Assigning customers to queues every 0.5 seconds
func customersToQueues(customers *[]Customer, tills *[]Till, lostCustomers *[]Customer, running *bool) {
	//sleep for 1 second so there's always customers generated before they're assigned
	time.Sleep(1 * time.Second)
	count := 0
	for *running {
		//Adds customer to shortest queue
		shortestTill := shortestQueue(tills)
		fmt.Println("Shortest till: ", shortestTill.tillId, " Queue length: ", len(shortestTill.queue))
		if shortestTill.checkLength() < 6 && len(*customers)!= 0 {
				shortestTill.queue <- (*customers)[0]
				fmt.Println("Assigning customers to till ", shortestTill.tillId, ": ", shortestTill.queue)

				//After added to queue, delete customer from slice
				*customers = append((*customers)[:0], (*customers)[0+1:]...)

		} else if shortestTill.checkLength() == 6 && len(*customers)!= 0  {
				fmt.Println("Customer lost: ", (*customers)[0])

				//Add to lost customers slice
				*lostCustomers = append(*lostCustomers, (*customers)[0])

				//Remove from original customers slice
				*customers = append((*customers)[:0], (*customers)[0+1:]...)
		}
		count++
	}
}

//Creating the initial till slice and opening a few of them
func createTills(tills *[]Till) {
	rand.Seed(time.Now().UnixNano())
	fmt.Print("Length: ", len(*tills))
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


    fmt.Print("Weather? 1=Bad, 2=Good: ")
    var weather int
    _, err := fmt.Scanf("%d", &weather)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
	fmt.Println(weather)

	//Variables
	running := true
	var customers []Customer
	var tills []Till
	var lostCustomers []Customer
	var processed []Customer
	var allCust []Customer


	//Setting up tills
	createTills(&tills)

	//Go routines

	go generateCustomers(&customers, &running, &weather, &allCust)
	go customersToQueues(&customers, &tills, &lostCustomers, &running)
	for i:= 0; i < 8; i ++ {
		go tills[i].processCustomers(&processed)
	}

	
	time.Sleep(20 * time.Second) 
	running = false
	fmt.Println("TIMES UP!")
	fmt.Println("All customers: ", allCust)
	fmt.Println("Processed customers: ", processed)
	fmt.Println("Lost customers: ", lostCustomers)

}