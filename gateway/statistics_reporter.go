package ceratopoGon

import (
	"log"
	"sync/atomic"
	"time"
)

type StatisticsReporter struct {
	sendPublish  uint64
	recvPublish  uint64
	pubComp      uint64
	sessionCount uint64
	interval     time.Duration
}

func NewStatisticsReporter(interval time.Duration) *StatisticsReporter {
	r := &StatisticsReporter{0, 0, 0, 0, interval}
	return r
}

func (r *StatisticsReporter) loggingLoop() {
	t := time.NewTicker(r.interval * time.Second)
	for {
		select {
		case <-t.C:
			log.Println("[STATISTICS-REPORT] SESSION-COUNT : ",
				atomic.LoadUint64(&r.sessionCount),
				", PUBLISH-RECV-COUNT : ",
				atomic.LoadUint64(&r.recvPublish),
				", PUBLISH-SEND-COUNT : ",
				atomic.LoadUint64(&r.sendPublish),
				", PUBLISH-COMP-COUNT : ",
				atomic.LoadUint64(&r.pubComp))
			atomic.StoreUint64(&r.sendPublish, 0)
			atomic.StoreUint64(&r.recvPublish, 0)
			atomic.StoreUint64(&r.pubComp, 0)
		}
	}
}

func (r *StatisticsReporter) countUpSendPublish() {
	atomic.AddUint64(&r.sendPublish, 1)
}

func (r *StatisticsReporter) countUpRecvPublish() {
	atomic.AddUint64(&r.recvPublish, 1)
}

func (r *StatisticsReporter) countUpPubComp() {
	atomic.AddUint64(&r.pubComp, 1)
}

func (r *StatisticsReporter) storeSessionCount(c uint64) {
	atomic.StoreUint64(&r.sessionCount, c)
}
