package parking

import (
	"testing"
)

func TestPark_Create(t *testing.T) {
	// creat mock park
	var numslot uint16 = 4
	park := createMockPark(numslot)
	if len(park.slotCh) != int(numslot) {
		t.Error("should be ", numslot, " but have ", len(park.slotCh))
	}

	numslot = 6
	park = createMockPark(numslot)
	if len(park.slotCh) != int(numslot) {
		t.Error("should be ", numslot, " but have ", len(park.slotCh))
	}
}

func TestPark_CarParking(t *testing.T) {
	var numslot uint16 = 4
	park := createMockPark(numslot)
	slotNumber := make(chan uint16)
	park.CarParking("KA-01-HH-9999", "White", slotNumber)
	number := <-slotNumber

	if number <= 0 || number != 1 {
		t.Error("can't park got slot", number)
	}

	slotNumber2 := make(chan uint16)
	park.CarParking("KA-01-HH-8888", "Black", slotNumber2)
	number2 := <-slotNumber2

	if number2 <= 0 || number2 != 2 {
		t.Error("can't park got slot ", number2)
	}

	slotNumber3 := make(chan uint16)
	park.CarParking("KA-01-HH-7777", "White", slotNumber3)
	number3 := <-slotNumber3

	if number3 <= 0 || number3 != 3 {
		t.Error("can't park got slot ", number3)
	}

	slotNumber4 := make(chan uint16)
	park.CarParking("KA-01-HH-6666", "Black", slotNumber4)
	number4 := <-slotNumber4
	if number4 <= 0 || number4 != 4 {
		t.Error("can't park got slot ", number4)
	}

	slotNumber5 := make(chan uint16)
	park.CarParking("KA-01-HH-5555", "Blue", slotNumber5)
	number5 := <-slotNumber5
	if number5 > 0 {
		t.Error("car should not park got slot ", number5)
	}

	if _, avilable := park.isThereASlotToJoin(); avilable {
		t.Error(" should no slot left ")
	}
}

func TestPark_CarLeave(t *testing.T) {
	var numslot uint16 = 3
	park := createMockPark(numslot)
	slotNumber := make(chan uint16)
	park.CarParking("KA-01-HH-9999", "White", slotNumber)
	<-slotNumber

	slotNumber2 := make(chan uint16)
	park.CarParking("KA-01-HH-8888", "Black", slotNumber2)
	<-slotNumber2

	slotNumber3 := make(chan uint16)
	park.CarParking("KA-01-HH-7777", "White", slotNumber3)
	number3 := <-slotNumber3

	freeNumber := make(chan uint16)
	park.CarLeave(number3, freeNumber)
	num := <-freeNumber
	if num <= 0 {
		t.Error("car can't leave", num)
	}

	freeNumber = make(chan uint16)
	park.CarLeave(number3, freeNumber)
	num = <-freeNumber
	if num != 0 {
		t.Error("leave not remove car", num)
	}

	if len(park.parkedCar) != 2 {
		t.Error("leave not free slot ", len(park.parkedCar))
	}

	if slot, avilable := park.isThereASlotToJoin(); avilable {
		if slot.number != 3 {
			t.Error("return wrong slot")
		}
	} else {
		t.Error("no return car slot")
	}

}

func TestPark_CarJoinAndLeave(t *testing.T) {
	var numslot uint16 = 3
	park := createMockPark(numslot)
	slotNumber := make(chan uint16)
	park.CarParking("KA-01-HH-9999", "White", slotNumber)
	<-slotNumber

	slotNumber2 := make(chan uint16)
	park.CarParking("KA-01-HH-8888", "Black", slotNumber2)
	number2 := <-slotNumber2

	slotNumber3 := make(chan uint16)
	park.CarParking("KA-01-HH-7777", "White", slotNumber3)
	<-slotNumber3

	// leave
	freeNumber := make(chan uint16)
	park.CarLeave(number2, freeNumber)
	numFree := <-freeNumber
	if len(park.parkedCar) != 2 {
		t.Error("leave not free slot ", len(park.parkedCar))
	}
	// join
	slotNumber2 = make(chan uint16)
	park.CarParking("KA-02-HH-2222", "White", slotNumber2)
	number2 = <-slotNumber2

	if len(park.parkedCar) != 3 {
		t.Error("join not remove slot ", len(park.parkedCar))
	}

	if numFree != number2 || numFree == 0 || number2 == 0 {
		t.Error("not parking at free slot")
	}

}

func TestPark_ParkStatus(t *testing.T) {
	var numslot uint16 = 3
	park := createMockPark(numslot)
	slotNumber := make(chan uint16)
	park.CarParking("KA-01-HH-9999", "White", slotNumber)
	<-slotNumber

	slotNumber2 := make(chan uint16)
	park.CarParking("KA-01-HH-8888", "Black", slotNumber2)
	number2 := <-slotNumber2

	slotNumber3 := make(chan uint16)
	park.CarParking("KA-01-HH-7777", "White", slotNumber3)
	<-slotNumber3

	// leave
	freeNumber := make(chan uint16)
	park.CarLeave(number2, freeNumber)
	<-freeNumber

	// join
	slotNumber2 = make(chan uint16)
	park.CarParking("KA-02-HH-2222", "White", slotNumber2)
	<-slotNumber2

	parkStatus := make(chan map[uint16]Slot)
	park.ParkStatus(parkStatus)
	status := <-parkStatus
	if detail, ok := status[1]; ok {
		if detail.regisNo != "KA-01-HH-9999" && detail.regisNo != "White" {
			t.Error("status not true")
		}
	} else {
		t.Error("car not found")
	}

	if detail, ok := status[2]; ok {
		if detail.regisNo != "KA-02-HH-2222" && detail.regisNo != "White" {
			t.Error("status not true")
		}
	} else {
		t.Error("car not found")
	}

	if detail, ok := status[3]; ok {
		if detail.regisNo != "KA-01-HH-7777" && detail.regisNo != "White" {
			t.Error("status not true")
		}
	} else {
		t.Error("car not found")
	}
}

func TestPark_FindRegistrationNumberCarWithColor(t *testing.T) {
	var numslot uint16 = 3
	park := createMockPark(numslot)
	slotNumber := make(chan uint16)
	park.CarParking("KA-01-HH-9999", "White", slotNumber)
	<-slotNumber

	slotNumber2 := make(chan uint16)
	park.CarParking("KA-01-HH-8888", "Black", slotNumber2)
	<-slotNumber2

	slotNumber3 := make(chan uint16)
	park.CarParking("KA-01-HH-7777", "White", slotNumber3)
	<-slotNumber3

	regisCh := make(chan []string)
	park.GetRegisNumberWithColor("White", regisCh)
	regisNumberList := <-regisCh
	var isPass bool = false
	for _, item := range regisNumberList {
		if item == "KA-01-HH-9999" || item == "KA-01-HH-7777" {
			isPass = true
		}
	}

	if !isPass {
		t.Error("regis number not match")
	}
}

func TestPark_FindSlotNumberWithCarColor(t *testing.T) {
	var numslot uint16 = 3
	park := createMockPark(numslot)
	slotNumber := make(chan uint16)
	park.CarParking("KA-01-HH-9999", "White", slotNumber)
	<-slotNumber

	slotNumber2 := make(chan uint16)
	park.CarParking("KA-01-HH-8888", "Black", slotNumber2)
	<-slotNumber2

	slotNumber3 := make(chan uint16)
	park.CarParking("KA-01-HH-7777", "White", slotNumber3)
	<-slotNumber3

	slotNumberCh := make(chan []uint16)
	park.GetSlotNumbersWithColor("Black", slotNumberCh)
	slotNumberList := <-slotNumberCh
	var isPass bool = false
	for _, item := range slotNumberList {
		if item == 2 {
			isPass = true
		}
	}

	if !isPass {
		t.Error("slot number not match")
	}
}

func TestPark_FindSlotNumberWithRegisNumber(t *testing.T) {
	var numslot uint16 = 3
	park := createMockPark(numslot)
	slotNumber := make(chan uint16)
	park.CarParking("KA-01-HH-9999", "White", slotNumber)
	<-slotNumber

	slotNumber2 := make(chan uint16)
	park.CarParking("KA-01-HH-8888", "Black", slotNumber2)
	<-slotNumber2

	slotNumber3 := make(chan uint16)
	park.CarParking("KA-01-HH-7777", "White", slotNumber3)
	<-slotNumber3

	slotNumberCh := make(chan uint16)
	park.GetSlotNumberWithRegisCarNumber("KA-01-HH-8888", slotNumberCh)
	number := <-slotNumberCh
	if number != 2 {
		t.Error("slot number not match")
	}

	slotNumberCh = make(chan uint16)
	park.GetSlotNumberWithRegisCarNumber("KA-01-HH-0000", slotNumberCh)
	number = <-slotNumberCh

	if number != 0 {
		t.Error("slot number not match")
	}
}

func TestPark_DestroyPark(t *testing.T) {
	var numslot uint16 = 3
	park := createMockPark(numslot)
	park.Destroy()
	if !IsClosed(park.eventCh) && park.canSendToEventCh {
		t.Error("channel not close")
	}

}
func createMockPark(numSlot uint16) *Park {
	return NewParking(numSlot)
}

func IsClosed(ch <-chan event) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
