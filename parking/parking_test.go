package parking

import (
	"testing"
)

func TestPark_Create(t *testing.T) {
	// creat mock park
	var numslot uint32 = 4
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
	park := createMockPark(6)
	isParked := make(chan bool)
	park.CarParking(isParked)
	canPark := <-isParked
	if canPark {

	} else {
		t.Error("cant park ", canPark)
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

func createMockPark(numSlot uint32) *Park {
	return NewParking(numSlot)
}
