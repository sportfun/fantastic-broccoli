package channel

//TODO: Do not close output chans
func Dispatch(inputChannel <-chan interface{}, outputChannels ...chan<- interface{}) error {
	if inputChannel == nil || len(outputChannels) == 0 {
		return nil
	}

	go func(rcvChn <-chan interface{}, sndChns []chan<- interface{}) {
		defer func() { recover() }() //Fix: avoid closed channel panic

		for v := range rcvChn {
			for _, chn := range sndChns {
				go func(c chan<- interface{}) { c <- v }(chn) //TODO: May create a nil chan if the chan is closed
			}
		}
	}(inputChannel, outputChannels)

	return nil
}
