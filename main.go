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
	startTime     time.Time
	endTime       time.Time
}

var totalWaitTime = 0.0
var tillsClosed = 0
var tillsOpened = 0
var totalTime = 60000.0

type Till struct {
	tillId          int
	scannerSpeed    float64
	queue           chan Customer
	opened          bool
	productsScanned int
	maxItems		int
}

func openTill(tillId int, tills *[]Till) {
	(*tills)[tillId].opened = true
	tillsOpened++
}

func closeTill(tillId int, tills *[]Till) {
	(*tills)[tillId].opened = false
	tillsClosed++
}

func (t *Till) checkLength() int {
	return len(t.queue)
}

func (c *Customer) startWaitTime() time.Time {
	startTime := time.Now()
	return startTime
}

func (c *Customer) endWaitTime(startTime time.Time) time.Duration {
	endTime := time.Now()
	waitTime := endTime.Sub(startTime)

	return waitTime
}

func (t *Till) processCustomers(running *bool, processed *[]Customer) {
	for *running {
		oneSecond := (((totalTime / 12) / 60) / 60)
		for customer := range t.queue {
			for i := 0; i < customer.numberOfItems; i++ {
				time.Sleep((time.Duration(t.scannerSpeed) * 10) * time.Millisecond)
				t.productsScanned++
			}
			customer.endTime = time.Now()
			*processed = append(*processed, customer)
			waitTime := customer.endTime.Sub(customer.startTime)
			fmt.Println("Customer ", customer.customerId, " left the store. Their wait time was", ((float64(waitTime) * oneSecond) / 60), "mins")
			totalWaitTime = totalWaitTime + float64(waitTime)
		}
	}
}

func generateCustomers(customers *[]Customer, genCustomers *bool, weather *int, allCustomers *[]Customer, result *int) {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Weather is: ", *weather)
	count := 0
	for *genCustomers {
		customer := Customer{
			customerId:    count,
			numberOfItems: (rand.Intn(200-1) + 1),
		}
		*customers = append(*customers, customer)
		*allCustomers = append(*allCustomers, customer)

		// records the number of products processed
		*result += customer.numberOfItems

		if *weather == 1 {
			time.Sleep(150 * time.Millisecond) //Every 2.4 mins
		} else if *weather == 2 {
			time.Sleep(200 * time.Millisecond) //Every 1.8 mins
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

		var queuesFull = 0

		for j := 0; j < 8; j++ {
			if (*tills)[j].checkLength() == 6 {
				queuesFull++
			}
		}

		if (*tills)[tillNumber].checkLength() <= 6 && (*tills)[tillNumber].opened {
			if len(*customers) != 0 {
				if (*customers)[0].numberOfItems < (*tills)[tillNumber].maxItems {
					(*customers)[0].startTime = time.Now()
					(*tills)[tillNumber].queue <- (*customers)[0]
					fmt.Println("Customer ", (*customers)[0].customerId, " entered the queue for till ", (*tills)[tillNumber].tillId)
					//After added to queue, delete customer from slice
					*customers = append((*customers)[:0], (*customers)[0+1:]...)
	
					time.Sleep(50 * time.Millisecond)
				}
			}
		} else if queuesFull == 7 && len(*customers) != 0 {
			fmt.Println("Customer ", (*customers)[0].customerId, " left the store as all the queues were full.")
			*lostCustomers = append(*lostCustomers, (*customers)[0])
			*customers = append((*customers)[:0], (*customers)[0+1:]...)
		}
		count++
	}
}

//Creating the initial till slice and opening a few of them
func createTills(tills *[]Till) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 8; i++ {
		if i == 0 {
			till := Till{
				tillId:       i + 1,
				scannerSpeed: float64(rand.Intn(4-1) + 1),
				queue:        make(chan Customer, 6),
				opened:       false,
				maxItems:	  20,
			}
			*tills = append(*tills, till)
		}
		till := Till{
			tillId:       i + 1,
			scannerSpeed: float64(rand.Intn(4-1) + 1),
			queue:        make(chan Customer, 6),
			opened:       false,
			maxItems:	  200,
		}
		*tills = append(*tills, till)
	}

	tillsOpen := (rand.Intn(9-1) + 1)
	fmt.Println("Tills open at start of day: ", tillsOpen)

	for i := 0; i < tillsOpen; i++ {
		(*tills)[i].opened = true
	}
}

func findShortestQueue(tills *[]Till, q1 chan int) {
	var shortest = 0
	var length = 6

	for i := 0; i < 8; i++ {
		if (*tills)[i].opened {
			if i == 0 {
				length = (*tills)[i].checkLength()
				shortest = i
			}

			if (*tills)[i].checkLength() < length {
				length = (*tills)[i].checkLength()
				shortest = i
			}
		}
	}
	q1 <- shortest
}

func calcTillsNeeded(tills *[]Till, running *bool) {
	time.Sleep(900 * time.Millisecond)
	var length = 0
	var openedTills = 0

	for *running {
		time.Sleep(900 * time.Millisecond)

		for i := 0; i < 8; i++ {
			if (*tills)[i].opened {

				openedTills++

			}
		}

		for i := 0; i < 8; i++ {
			if (*tills)[i].opened {

				length = (*tills)[i].checkLength()

				if length > 4 {
					for z := 0; z < 8; z++ {
						if (*tills)[z].opened == false {
							openTill(z, tills)
							fmt.Println("Opened till:", (*tills)[z].tillId)
							break
						}
					}
				}

			}
		}

		for i := 0; i < 8; i++ {
			if (*tills)[i].opened {
				length = (*tills)[i].checkLength()

				if length <= 2 {
					if openedTills >= 4 {
						closeTill(i, tills)
						fmt.Println("Closed till:", (*tills)[i].tillId)
						break
					}
				}
			}
		}

	}
}

func checkCustomerEmpty(customers *[]Customer, running *bool) {
	if (len(*customers)) == 0 {
		fmt.Println("The store has closed. The tills will now process the remaining customers.")
		*running = false
	}
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

	fmt.Println("Store is now open.")


	//Variables
	running := true
	genCustomers := true
	var customers []Customer
	var allCustomers []Customer
	var tills []Till
	var lostCustomers []Customer
	var processedCustomers []Customer
	var totalProducts int

	//Setting up tills
	createTills(&tills)

	//Go routines
	go generateCustomers(&customers, &genCustomers, &weather, &allCustomers, &totalProducts)
	go customersToQueues(&customers, &tills, &lostCustomers, &running)
	go calcTillsNeeded(&tills, &running)

	for i := 0; i < 8; i++ {
		go tills[i].processCustomers(&running, &processedCustomers)
	}

	time.Sleep(60 * time.Second)
	genCustomers = false

	go checkCustomerEmpty(&customers, &running)

	time.Sleep(20 * time.Second)

	fmt.Println("Store is now closed.")
	fmt.Println("Statistics:")
	fmt.Println("")
	fmt.Println("Total Number of customers generated: ", len(allCustomers))
	fmt.Println("Average wait time per customer: ", totalWaitTime/float64(len(allCustomers)), " min")
	fmt.Println("Total Number of processed customers: ", len(processedCustomers))
	fmt.Println("Total Number of products in all trolleys: ", totalProducts)
	fmt.Println("Average products per trolley: ", totalProducts/len(allCustomers))
	fmt.Println("Average till utilisation: ", totalProducts/len(tills))
	fmt.Println("Number of times tills opened: ", tillsOpened)
	fmt.Println("Number of times tills closed: ", tillsClosed)
	fmt.Println("Number of lost customers ", len(lostCustomers))
}
