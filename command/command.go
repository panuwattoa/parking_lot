package command

import (
	"fmt"
	"os"
	"parkinglot/parking"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type cmdFunc func(cmd string, param string)

var park *parking.Park

var CommandMapper = map[string]cmdFunc{
	"create_parking_lot": commandCreatePark,
	"park":               commandPark,
	"leave":              commandLeave,
	"status":             commandStatus,
	"registration_numbers_for_cars_with_colour": commandCheckRegisNumberWithColor,
	"slot_numbers_for_cars_with_colour":         commandCheckSlotNumberWithClor,
	"slot_number_for_registration_number":       commandCheckSlotNumberWithRegisNumber,
	"exit":                                      exit,
}

func commandCreatePark(cmd string, param string) {
	var parkslot uint16
	if count, err := fmt.Sscan(param, &parkslot); err == nil && count == 1 {
		park = parking.NewParking(parkslot)
		fmt.Printf("Created a parking lot with %d slots", parkslot)
	} else {
		fmt.Printf("commannd[%s] need only 1 parameter.", cmd)
	}
}

func commandPark(cmd string, param string) {
	if checkParkAready() {
		var regisNumber string
		var color string
		if count, err := fmt.Sscan(param, &regisNumber, &color); err == nil && count == 2 {
			slotCh := make(chan uint16)
			park.CarParking(regisNumber, color, slotCh)
			slotNumber := <-slotCh
			if slotNumber > 0 {
				fmt.Printf("Allocated slot number: %d", slotNumber)
			} else {
				fmt.Println("Sorry, parking lot is full")
			}
		} else {
			fmt.Printf("commannd[%s] need 2 parameter 'Registration No' and 'Colour'", cmd)
		}

	}
}

func commandLeave(cmd string, param string) {
	if checkParkAready() {
		var slotnumber uint16
		if count, err := fmt.Sscan(param, &slotnumber); err == nil && count == 1 {
			slotCh := make(chan uint16)
			park.CarLeave(slotnumber, slotCh)
			slotNumber := <-slotCh
			if slotNumber > 0 {
				fmt.Printf("Slot number %d is free", slotNumber)
			} else {
				fmt.Println("Sorry, slot already free")
			}
		} else {
			fmt.Printf("commannd[%s] need only 1 parameter.", cmd)
		}
	}
}

func commandStatus(cmd string, param string) {
	if checkParkAready() {
		parkCarsCh := make(chan map[uint16]parking.Slot)
		park.ParkStatus(parkCarsCh)
		parkedCars := <-parkCarsCh

		// initialize tabwriter
		w := new(tabwriter.Writer)

		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)

		defer w.Flush()

		fmt.Fprintf(w, "\n %s\t%s\t%s\t", "Slot No.", "Registration No", "Colour")
		sortKey := sortMap(parkedCars)
		for _, value := range sortKey {
			fmt.Fprintf(w, "\n %d\t%s\t%s\t", value, parkedCars[uint16(value)].RegisNo, parkedCars[uint16(value)].Color)
		}
	}
}

func commandCheckRegisNumberWithColor(cmd string, param string) {
	if checkParkAready() {
		var color string
		if count, err := fmt.Sscan(param, &color); err == nil && count == 1 {
			regisNumbersCh := make(chan []string)
			park.GetRegisNumberWithColor(color, regisNumbersCh)
			regisNumberList := <-regisNumbersCh
			if len(regisNumberList) > 0 {
				fmt.Printf("%s", strings.Join(regisNumberList, ", "))
			} else {
				fmt.Println("Not found")
			}
		} else {
			fmt.Printf("commannd[%s] need only 1 parameter.", cmd)
		}
	}
}

func commandCheckSlotNumberWithClor(cmd string, param string) {
	if checkParkAready() {
		var color string
		if count, err := fmt.Sscan(param, &color); err == nil && count == 1 {
			slotNumberCh := make(chan []uint16)
			park.GetSlotNumbersWithColor(color, slotNumberCh)
			slotNumberList := <-slotNumberCh
			if len(slotNumberList) > 0 {
				slotNumberString := make([]string, 0)
				for _, value := range slotNumberList {
					slotNumberString = append(slotNumberString, strconv.Itoa(int(value)))
				}
				fmt.Printf("%s", strings.Join(slotNumberString, ", "))
			} else {
				fmt.Println("Not found")
			}
		} else {
			fmt.Printf("commannd[%s] need only 1 parameter.", cmd)
		}
	}
}

func commandCheckSlotNumberWithRegisNumber(cmd string, param string) {
	if checkParkAready() {
		var regisNumber string
		if count, err := fmt.Sscan(param, &regisNumber); err == nil && count == 1 {
			slotNumberCh := make(chan uint16)
			park.GetSlotNumberWithRegisCarNumber(regisNumber, slotNumberCh)
			slotNumber := <-slotNumberCh
			if slotNumber > 0 {
				fmt.Println(slotNumber)
			} else {
				fmt.Println("Not found")
			}
		} else {
			fmt.Printf("commannd[%s] need only 1 parameter.", cmd)
		}
	}
}

func exit(cmd string, param string) {
	if checkParkAready() {
		park.Destroy()
		for i := 0; i < 5; i++ {
			fmt.Print("#")
			time.Sleep(time.Millisecond * 200)
		}
		fmt.Println()
	}
	os.Exit(1)

}

func checkParkAready() bool {
	if park == nil {
		fmt.Println("Parking not ready please 'create_parking_lot' first")
		return false
	}
	return true
}

func sortMap(m map[uint16]parking.Slot) (index []int) {
	keys := make([]int, 0)
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return keys
}
