package notification

type ChansStore map[string]chan struct{}

var NotificationChans ChansStore = make(ChansStore)

func (cs ChansStore) Add(uuid string, ch chan struct{}) {
	cs[uuid] = ch
}

func (cs ChansStore) Remove(uuid string) {
	close(cs[uuid])
	delete(cs, uuid)
}
