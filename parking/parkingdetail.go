package parking

// for event loop queuing
type event interface{}

type carParkingEvent struct {
	isParked    chan bool
	regisNumber string
	color       string
}

// Park is struct park control role and car data.
type Park struct {
	canSendToEventCh bool // for check channel available
	eventCh          chan event
	slotCh           chan *slot // slot car channel
	parkedCar        map[uint32]*slot
	eventDoneCh      chan struct{} // for unit test
}

type slot struct {
	number  uint32
	regisNo string
	color   string
}
