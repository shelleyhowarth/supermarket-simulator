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
}

var processed []Customer

type Till struct {
	tillId       int
	scannerSpeed float64
	queue        chan Customer
	opened       bool
}

func openTill(tillId int, tills *[]Till) {
	(*tills)[tillId].opened = true
}

func closeTill(tillId int, tills *[]Till) {
	(*tills)[tillId].opened = false
}

func (t *Till) checkLength() int {
	return len(t.queue)
}

func (t *Till) processCustomers(running *bool) {
	// processedCustomers := make(chan Customer)
	for *running {
		for customer := range t.queue {
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(customer)
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
	//good weather or bad weather
	//weather := (rand.Intn(2-1)+1)
	fmt.Println("Weather is: ", *weather)
	count := 0
	for *running {
		customer := Customer{
			customerId:    count,
			numberOfItems: (rand.Intn(200-1) + 1),
		}
		*customers = append(*customers, customer)
		*allCustomers = append(*allCustomers, customer)
		fmt.Println("Customers generated: ", *customers)

		// records the number of products processed
		*result += customer.numberOfItems
		fmt.Println("Customers generated: ", *customers)
		count++
		if *weather == 1 {
			time.Sleep(150 * time.Millisecond)
		} else if *weather == 2 {
			time.Sleep(100 * time.Millisecond)
		}
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

		fmt.Println(tillNumber)

		if (*tills)[tillNumber].checkLength() <= 6 && (*tills)[tillNumber].opened {
			if len(*customers) != 0 {
				(*tills)[tillNumber].queue <- (*customers)[0]
				fmt.Println("Assigning customers to till ", tillNumber+1, ": ", (*tills)[1].queue)
				//After added to queue, delete customer from slice
				*customers = append((*customers)[:0], (*customers)[0+1:]...)
				fmt.Println("Slice after assignment", *customers)
				time.Sleep(80 * time.Millisecond)
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
			fmt.Println("till opened: ", (*tills)[i].tillId, "que lenght: ", (*tills)[i].checkLength())

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

func startTillProcess(customers *[]Customer, tills *[]Till, running *bool) {
	for i := 0; i < 8; i++ {
		go (*tills)[i].processCustomers(running)
	}
}

func calcTillsNeeded(tills *[]Till, running *bool) {
	// // var under3 = 0
	// var over3 = 0
	// var zeroCust = 0
	// var openedTills = 0
	// var length = 0
	// // tooMany := false
	// // tooLittle := false
	var length = 0
	// var over3 = 0
	// var zeroCust = 0
	var openedTills = 0

	for *running {
		time.Sleep(1000 * time.Millisecond)
		for i := 0; i < 8; i++ {
			if (*tills)[i].opened {

				openedTills++

				length = (*tills)[i].checkLength()

				if length > 3 {
					for z := 0; z < 8; z++ {
						if (*tills)[z].opened == false {
							openTill(z, tills)
							fmt.Println("OPENED TILL:", z)
							break

						}
					}
				}

			}
		}

		for i := 0; i < 8; i++ {
			if (*tills)[i].opened {
				openedTills++

				length = (*tills)[i].checkLength()

				if length <= 2 {
					if openedTills >= 4 {
						closeTill(i, tills)
						break
					}
				}
			}
		}

		// if zeroCust >= 2 {
		// 	for i := 0; i < 8; i++ {
		// 		if (*tills)[i].opened && openedTills > 5 {
		// 			if (*tills)[i].checkLength() <= 1 {
		// 				closeTill(i, tills)
		// 				fmt.Println("CLOSING TILL:", i)
		// 				break
		// 			}
		// 		}
		// 	}

		// 	if over3 >= 2 {
		// 		for i := 0; i < 8; i++ {
		// 			if !(*tills)[i].opened {
		// 				openTill(i, tills)
		// 				fmt.Println("OPENED TILL:", i)
		// 				break
		// 			}

		// 		}
		// 	}

		// if zeroCust >= 2 {
		// 	if openedTills >= 4 {
		// 		tooMany = true

		// 	}
		// } else if over3 > 1 {
		// 	tooLittle = true

		// }

		// 	if tooLittle {
		// 		for i := 0; i < 8; i++ {
		// 			if !(*tills)[i].opened {
		// 				openTill(i, tills)
		// 				fmt.Println("OPENED TILL:", i)
		// 				break
		// 			}

		// 		}
		// 	}

		// 	if tooMany {
		// 		for i := 0; i < 8; i++ {
		// 			if (*tills)[i].opened {
		// 				if (*tills)[i].checkLength() == 0 {
		// 					closeTill(i, tills)
		// 					fmt.Println("CLOSING TILL:", i)
		// 					break
		// 				}
		// 			}
		// 		}

		// 	}
	}
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
	var allCustomers []Customer
	var tills []Till
	var lostCustomers []Customer
	var result int

	//Setting up tills
	createTills(&tills)

	//Go routines
	go generateCustomers(&customers, &running, &weather, &allCustomers, &result)
	go customersToQueues(&customers, &tills, &lostCustomers, &running)
	go startTillProcess(&customers, &tills, &running)
	go calcTillsNeeded(&tills, &running)

	//totalProductsProccessed(&customers)

	time.Sleep(20 * time.Second)
	fmt.Println("TIMES UP!")
	fmt.Println("Lost customers: ", lostCustomers)
	fmt.Println("Processed customers: ", processed)
	fmt.Println("Total Number of Products: ", result)
	fmt.Println("Average Products per person: ", result/len(allCustomers))

	running = false
}
