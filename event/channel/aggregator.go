package channel

//TODO: Do not close output chan
func Aggregate(outputChannel chan<- interface{}, inputChannels ...<-chan interface{}) error {
	if outputChannel == nil || len(inputChannels) == 0 {
		return nil
	}

	for _, ch := range inputChannels {
		go func(sndChn chan<- interface{}, rcvChn <-chan interface{}) {
			defer func() { recover() }() //Fix: avoid closed channel panic

			for v := range rcvChn {
				sndChn <- v
			}
		}(outputChannel, ch)
	}

	return nil
}
