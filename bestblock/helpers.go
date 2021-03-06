package bestblock

import "time"

// getBlockNumber calls the eth_blockNumber function on the node and saves the result in the
// dedicated variable on the Tracker struct
func (b *Tracker) getBlockNumber() {
	b.logger.Trace("getting best block")
	start := time.Now()

	block, err := b.conn.GetBlockNumber()
	d := time.Since(start)

	if err != nil {
		b.logger.Error(err)
		b.errChan <- err
	} else {
		b.logger.WithField("block", block).WithField("duration", d).Trace("got best block")
		b.mu.Lock()
		b.bestBlockNumber = block
		b.mu.Unlock()

		if !b.started {
			b.started = true
			b.errChan <- nil
		}
	}
}
