package parking

// NewParking for init park with maximum slot
func NewParking(maxSlot uint16) *Park {
	park := Park{
		canSendToEventCh: true,
		eventCh:          make(chan event, 32),
		slotCh:           make(chan *slot, maxSlot),
		parkedCar:        make(map[uint16]*slot),
	}

	for i := 1; i <= int(maxSlot); i++ {
		park.slotCh <- &slot{
			number: uint16(i),
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
			p.carParking(ev.slotNumber)
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
func (p *Park) CarParking(regisNumber string, color string, slotNumberCb chan uint16) {
	p.eventCh <- carParkingEvent{
		slotNumber:  slotNumberCb,
		regisNumber: regisNumber,
		color:       color,
	}
}

func (p *Park) carParking(slotNumber chan uint16) {
	defer close(slotNumber)
	slotNumber <- 0
}
