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
	var numslot uint16 = 6
	park := createMockPark(numslot)
	slotNumber := make(chan uint16)
	park.CarParking("KA-01-HH-9999", "White", slotNumber)
	number := <-slotNumber

	if number <= 0 || number > numslot {
		t.Error("cant park ", numslot)
	}
}

func TestPark_CarLeave(t *testing.T) {

}

func TestPark_ParkStatus(t *testing.T) {

}

func TestPark_FindRegistrationNumberCarWithColor(t *testing.T) {

}

func TestPark_FindSlotNumberWithCarClor(t *testing.T) {

}

func TestPark_FindSlotNumberWithRegisNumber(t *testing.T) {

}

func createMockPark(numSlot uint16) *Park {
	return NewParking(numSlot)
}
