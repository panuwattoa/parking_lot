package parking

// for event loop queuing
type event interface{}

type carParkingEvent struct {
	slotNumber  chan uint16
	regisNumber string
	color       string
}

type carLeaveEvent struct {
	leaveSlotNumber uint16
	slotNumberFree  chan uint16
}

type parkStatusEvent struct {
	parkedCar chan map[uint16]Slot
}

type regisNumbersWithColor struct {
	color        string
	regisNumbers chan []string
}

type slotNumbersWithColor struct {
	color       string
	slotNumbers chan []uint16
}

type destroyEvent struct{}

type slotNumberWithRegisNumber struct {
	regisNumber string
	slotNumbers chan uint16
}

// Park is struct park control role and car data.
type Park struct {
	canSendToEventCh bool // for check channel available
	eventCh          chan event
	slotCh           chan *Slot // slot car channel
	parkedCar        map[uint16]*Slot
	eventDoneCh      chan struct{} // for unit test
}

type Slot struct {
	number  uint16
	regisNo string
	color   string
}
