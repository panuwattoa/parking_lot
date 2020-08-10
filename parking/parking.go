package parking

// NewParking for init park with maximum slot
func NewParking(maxSlot uint32) *Park {
	park := Park{
		canSendToEventCh: true,
		eventCh:          make(chan event, 32),
		slotCh:           make(chan *slot, maxSlot),
		parkedCar:        make(map[uint32]*slot),
	}

	for i := 1; i <= int(maxSlot); i++ {
		park.slotCh <- &slot{
			number: uint32(i),
		}
	}
	go park.eventLoop()

	return &park
}

func (p *Park) safeCloseEventChannel() {
	if p.canSendToEventCh {
		close(p.eventCh)
	}
	p.canSendToEventCh = false
}

// eventLoop for protect race condition
func (p *Park) eventLoop() {
	for event := range p.eventCh {
		switch ev := event.(type) {
		case carParkingEvent:
			p.carParking(ev.isParked)
		}
		if p.eventDoneCh != nil {
			p.eventDoneCh <- struct{}{}
		}
	}
	if p.eventDoneCh != nil {
		close(p.eventDoneCh)
	}
	close(p.slotCh)
}

// CarParking is
func (p *Park) CarParking(isParkedCb chan bool) {
	p.eventCh <- carParkingEvent{
		isParked: isParkedCb,
	}
}

func (p *Park) carParking(isParked chan bool) {
	defer close(isParked)
	isParked <- false
}
