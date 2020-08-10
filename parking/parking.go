package parking

// NewParking for init park with maximum slot
func NewParking(maxSlot uint16) *Park {
	park := Park{
		canSendToEventCh: true,
		eventCh:          make(chan event, 32),
		slotCh:           make(chan *Slot, maxSlot),
		parkedCar:        make(map[uint16]*Slot),
	}

	for i := 1; i <= int(maxSlot); i++ {
		park.slotCh <- &Slot{
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
			p.parking(ev.regisNumber, ev.color, ev.slotNumber)
		case carLeaveEvent:
			p.leave(ev.leaveSlotNumber, ev.slotNumberFree)
		case parkStatusEvent:
			p.status(ev.parkedCar)
		case regisNumbersWithColor:
			p.findRegisNumberListWithColor(ev.color, ev.regisNumbers)
		case slotNumbersWithColor:
			p.findSlotNumberListWithColor(ev.color, ev.slotNumbers)
		case slotNumberWithRegisNumber:
			p.findSlotNumberWithCarRegisNumber(ev.regisNumber, ev.slotNumbers)
		case destroyEvent:
			p.safeCloseEventChannel()
		}
	}
	close(p.slotCh)
}

// CarParking for parking get car slot
func (p *Park) CarParking(regisNumber string, color string, slotNumberCb chan uint16) {
	if !p.canSendToEventCh {
		return
	}
	p.eventCh <- carParkingEvent{
		slotNumber:  slotNumberCb,
		regisNumber: regisNumber,
		color:       color,
	}
}

// CarLeave for leave park return slot free
func (p *Park) CarLeave(leaveSlotNumber uint16, slotNumberCb chan uint16) {
	if !p.canSendToEventCh {
		return
	}
	p.eventCh <- carLeaveEvent{
		leaveSlotNumber: leaveSlotNumber,
		slotNumberFree:  slotNumberCb,
	}
}

// ParkStatus check parking status
func (p *Park) ParkStatus(parkedCar chan map[uint16]Slot) {
	if !p.canSendToEventCh {
		return
	}
	p.eventCh <- parkStatusEvent{
		parkedCar: parkedCar,
	}
}

// GetRegisNumberWithColor for get regis number list of car in parking
func (p *Park) GetRegisNumberWithColor(color string, regisNumbers chan []string) {
	if !p.canSendToEventCh {
		return
	}
	p.eventCh <- regisNumbersWithColor{
		color:        color,
		regisNumbers: regisNumbers,
	}
}

// GetSlotNumbersWithColor for get slot number list with car color
func (p *Park) GetSlotNumbersWithColor(color string, slotNumbers chan []uint16) {
	if !p.canSendToEventCh {
		return
	}
	p.eventCh <- slotNumbersWithColor{
		color:       color,
		slotNumbers: slotNumbers,
	}
}

// GetSlotNumberWithRegisCarNumber for get slot number with car regis number
func (p *Park) GetSlotNumberWithRegisCarNumber(regisNumber string, slotNumbers chan uint16) {
	if !p.canSendToEventCh {
		return
	}
	p.eventCh <- slotNumberWithRegisNumber{
		regisNumber: regisNumber,
		slotNumbers: slotNumbers,
	}
}

// Destroy close channel
func (p *Park) Destroy() {
	p.eventCh <- destroyEvent{}
}

func (p *Park) parking(regisNumber string, color string, slotNumber chan uint16) {
	defer close(slotNumber)
	if slot, available := p.isThereASlotToJoin(); available {
		slot.color = color
		slot.regisNo = regisNumber
		p.parkedCar[slot.number] = slot
		slotNumber <- slot.number // return
	} else {
		slotNumber <- 0 // return
	}
}

func (p *Park) leave(leaveSlotNumber uint16, slotNumberFree chan uint16) {
	defer close(slotNumberFree)
	if slot, ok := p.parkedCar[leaveSlotNumber]; ok {
		var slotPosition = slot.number
		p.slotCh <- slot // return slot free
		delete(p.parkedCar, leaveSlotNumber)
		slotNumberFree <- slotPosition // return
	} else {
		slotNumberFree <- 0 // return
	}
}

func (p *Park) status(parkedCar chan map[uint16]Slot) {
	defer close(parkedCar)
	copyMap := make(map[uint16]Slot)
	for key, value := range p.parkedCar {
		copyMap[key] = *value
	}
	parkedCar <- copyMap
}

func (p *Park) findRegisNumberListWithColor(color string, regisNumbers chan []string) {
	defer close(regisNumbers)
	regisList := make([]string, len(p.parkedCar))
	for _, car := range p.parkedCar {
		if car.color == color {
			regisList = append(regisList, car.regisNo)
		}
	}
	regisNumbers <- regisList
}

func (p *Park) findSlotNumberListWithColor(color string, slotNumber chan []uint16) {
	defer close(slotNumber)
	slotList := make([]uint16, len(p.parkedCar))
	for _, car := range p.parkedCar {
		if car.color == color {
			slotList = append(slotList, car.number)
		}
	}
	slotNumber <- slotList
}

func (p *Park) findSlotNumberWithCarRegisNumber(regisNumber string, slotNumber chan uint16) {
	defer close(slotNumber)
	for _, car := range p.parkedCar {
		if car.regisNo == regisNumber {
			slotNumber <- car.number
			return
		}
	}
	slotNumber <- 0
}

func (p *Park) isThereASlotToJoin() (*Slot, bool) {
	if len(p.slotCh) == 0 {
		return nil, false
	}

	for slot := range p.slotCh {
		return slot, true
	}

	return nil, false
}
