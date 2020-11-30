//12 hours -> 60 seconds
//1 hour -> 5 seconds
//30 mins -> 2.5 seconds
//6 mins -> 0.5 seconds
//3 mins -> 0.25 seconds
//1.5 min ->  ->0.125 seconds
//1 second -> 0.000084 seconds
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Customer struct {
	customerId    int
	numberOfItems int
	startTime 	  time.Time
	endTime 	  time.Time
}

var totalWaitTime = 0.0
var totalTime = 60000.0

type Till struct {
	tillId       int
	scannerSpeed float64
	queue        chan Customer
	opened       bool
	productsScanned int
}


func (t *Till) checkLength() int {
	return len(t.queue)
}

func (c *Customer) startWaitTime() (time.Time) {
	startTime := time.Now()
	return startTime
}

func (c *Customer) endWaitTime(startTime time.Time) (time.Duration) {
	endTime := time.Now()
	waitTime := endTime.Sub(startTime)

	return waitTime	
}

func (t *Till) processCustomers(running *bool, processed *[]Customer) {
	for *running {
		oneSecond := (((totalTime/12)/60)/60)
		for customer := range t.queue {
			for i:=0; i < customer.numberOfItems; i++ {
				time.Sleep((time.Duration(t.scannerSpeed)) * time.Millisecond)
				t.productsScanned++
			}
			customer.endTime = time.Now()
			*processed = append(*processed, customer)
			waitTime := customer.endTime.Sub(customer.startTime)
			fmt.Println("Customer ID: ", customer.customerId, " Wait time: ", ((float64(waitTime)*oneSecond)/60), " min")
			totalWaitTime = totalWaitTime + float64(waitTime)
		}

	}

	// for *running {
	// 	if t.checkLength() > 0 {
	// 		fmt.Println("OVER 0 CUSTOMERS PLEASE PROCESS")
	// 	}
	// 	for customer := range t.queue {
	// 		for i := 0; i < customer.numberOfItems; i++ {
	// 			time.Sleep(50 * time.Millisecond) //change this to scanning speed
	// 			//When processing last item
	// 			if i == customer.numberOfItems-1 {
	// 				time.Sleep(50 * time.Millisecond) //change this to scanning speed
	// 				//Remove customer from channel
	// 				fmt.Println("processed ", customer)
	// 				processed = append(processed, customer)
	// 				processedCustomers <- customer
	// 			}
	// 		}
	// 	}
	// }
}

//Create customers every 0.3 or 0.5 seconds
func generateCustomers(customers *[]Customer, running *bool, weather *int, allCustomers *[]Customer, result *int) {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Weather is: ", *weather)
	count := 0
	for *running {
		customer := Customer{
			customerId:    count,
			numberOfItems: (rand.Intn(200-1) + 1),
		}
		*customers = append(*customers, customer)
		*allCustomers = append(*allCustomers, customer)

		// records the number of products processed
		*result += customer.numberOfItems

		if *weather == 1 {
			time.Sleep(80 * time.Millisecond) //Every 60 secs
		} else if *weather == 2 {
			time.Sleep(120 * time.Millisecond) //Every 120 secs
		}
		count++
	}
}

//Assigning customers to queues every 0.5 seconds
func customersToQueues(customers *[]Customer, tills *[]Till, lostCustomers *[]Customer, running *bool) {

	//sleep for 1 second so there's always customers generated before they're assigned
	time.Sleep(1000 * time.Millisecond)
	count := 0
	q1 := make(chan int)

	for *running {

		go findShortestQueue(tills, q1)
		tillNumber := <-q1

		//fmt.Println(tillNumber)

		if (*tills)[tillNumber].checkLength() <= 6 && (*tills)[tillNumber].opened {
			if len(*customers) != 0 {
				(*customers)[0].startTime = time.Now()
				(*tills)[tillNumber].queue <- (*customers)[0]
				fmt.Println("Assigning customers to till ", tillNumber+1, ": ", (*tills)[1].queue)
				//After added to queue, delete customer from slice
				*customers = append((*customers)[:0], (*customers)[0+1:]...)
				//time.Sleep(200 * time.Millisecond)
			}
		}

		// for (*tills)[i].checkLength() < 6 && (*tills)[i].opened {
		// 	//Adds customer to queue
		// 	if len(*customers)!= 0 {
		// 		(*tills)[i].queue <- (*customers)[0]
		// 		fmt.Println("Assigning customers to till ", i+1, ": ", (*tills)[1].queue)
		// 		//After added to queue, delete customer from slice
		// 		*customers = append((*customers)[:0], (*customers)[0+1:]...)
		// 		fmt.Println("Slice after assignment", *customers)
		// 		time.Sleep(500 * time.Millisecond)
		// 	}
		// }
		// if (*tills)[i].checkLength() == 6 {
		// 	//check other till lengths
		// 	fmt.Println("Customer lost: ", (*customers)[0])
		// 	//Add to lost customers slice
		// 	*lostCustomers = append(*lostCustomers, (*customers)[0])
		// 	//Remove from original customers slice

		// 	*customers = append((*customers)[:0], (*customers)[0+1:]...)
		// }
		// }
		count++
	}
}

//Creating the initial till slice and opening a few of them
func createTills(tills *[]Till) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 8; i++ {
		till := Till{
			tillId:       i + 1,
			scannerSpeed: float64(rand.Intn(4-1) + 1),
			queue:        make(chan Customer, 6),
			opened:       false,
		}
		*tills = append(*tills, till)
	}

	tillsOpen := (rand.Intn(9-1) + 1)
	fmt.Println("Tills open at start of day: ", tillsOpen)

	for i := 0; i < tillsOpen; i++ {
		(*tills)[i].opened = true
	}
	fmt.Println("Tills at start of day: ", *tills)
}

func findShortestQueue(tills *[]Till, q1 chan int) {
	var shortest = 0
	var length = 6

	for i := 0; i < 8; i++ {
		if (*tills)[i].opened {
			fmt.Println("Till opened: ", (*tills)[i].tillId, "Queue length: ", (*tills)[i].checkLength())
			if i == 0 {
				length = (*tills)[i].checkLength()
				shortest = i
			}

			if (*tills)[i].checkLength() < length {
				fmt.Println("Shortest till length", (*tills)[i].checkLength())
				length = (*tills)[i].checkLength()
				shortest = i
			}
		}
	}
	q1 <- shortest
}


func main() {

	fmt.Print("Weather? 1 = Bad, 2 = Good: ")
	var weather int
	_, err := fmt.Scanf("%d", &weather)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if weather != 1 && weather != 2 {
		fmt.Println("You didn't input 1 or 2")
		os.Exit(3)
	}

	//Variables
	running := true
	var customers []Customer
	var allCustomers []Customer
	var tills []Till
	var lostCustomers []Customer
	var processedCustomers []Customer
	var totalProducts int

	//Setting up tills
	createTills(&tills)

	//Go routines
	go generateCustomers(&customers, &running, &weather, &allCustomers, &totalProducts)
	go customersToQueues(&customers, &tills, &lostCustomers, &running)
	for i := 0; i < 8; i++ {
		go tills[i].processCustomers(&running, &processedCustomers)
	}

	time.Sleep(20 * time.Second)
	running = false
	fmt.Println("TIMES UP!")
	fmt.Println("Total Number of customers generated: ", len(allCustomers))
	fmt.Println("Average wait time per customer: " , totalWaitTime/float64(len(allCustomers)), " min")
	fmt.Print("Processed customers: ")
	for i:= 0; i<len(processedCustomers); i++ {
		fmt.Print("{", processedCustomers[i].customerId, ", ", processedCustomers[i].numberOfItems, "}")
	}
	fmt.Println("\nTotal Number of processed customers: ", len(processedCustomers))
	fmt.Println("\nTotal Number of Products: ", totalProducts)
	fmt.Println("Average Products per person: ", totalProducts/len(allCustomers))
	fmt.Println("Lost customers: ", lostCustomers)
	for i:= 0; i < len(tills); i++ {
		fmt.Println("Till ID: ", tills[i], " Total products scanned: ", tills[i].productsScanned)
	}
	fmt.Println("Average till utilisation: ", totalProducts/len(tills))

}
