// Package auditor is used to store authorization audit data to the queue.
package auditor

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/vmihailenco/msgpack/v5"

	"iam-authz/internal/authzserver/queue"
)

const queueName = "iam-authz-audit"

// AuditRecord defines the details of a authorization request.
type AuditRecord struct {
	Timestamp  int64
	Username   string
	Effect     string
	Conclusion string
	Request    string
	Policies   string
	Deciders   string
}

// Auditor defines the structure of an auditor.
type Auditor struct {
	que         queue.Queue
	recordsChan chan *AuditRecord
	poolSize    int
	omitDetails bool
	shouldStop  uint32
	ctx         context.Context
	poolWg      sync.WaitGroup
	log         *logger.Logger
}

var auditor *Auditor

// InitAuditor initializes the global auditor and returns it.
func InitAuditor(ctx context.Context, opts *AuditorOptions, que queue.Queue) *Auditor {
	log := logger.L()
	log.Debugf("building auditor with options: %+v", opts)

	auditor = &Auditor{
		que:         que,
		recordsChan: make(chan *AuditRecord, opts.BufferSize),
		poolSize:    opts.PoolSize,
		omitDetails: opts.OmitDetails,
		ctx:         ctx,
		log:         log,
	}

	return auditor
}

// GetAuditor returns the global auditor.
func GetAuditor() *Auditor {
	return auditor
}

// Start starts the auditor.
func (a *Auditor) Start() {
	atomic.SwapUint32(&a.shouldStop, 0)
	for i := 0; i < a.poolSize; i++ {
		a.poolWg.Add(1)
		go a.startWorker()
	}
}

// Stop flushes the buffer and stop the auditor.
func (a *Auditor) Stop(ctx context.Context) error {
	// flag to stop sending records into channel
	atomic.SwapUint32(&a.shouldStop, 1)

	// close channel to stop workers
	close(a.recordsChan)

	// wait for all workers to be done
	a.poolWg.Wait()

	return nil
}

// RecordHit stores an AuditRecord in queue.
func (a *Auditor) RecordHit(r *AuditRecord) {
	// check if we should stop sending records
	if atomic.LoadUint32(&a.shouldStop) > 0 {
		return
	}

	a.recordsChan <- r
}

func (a *Auditor) startWorker() {
	defer a.poolWg.Done()

	for {
		select {
		case record, ok := <-a.recordsChan:
			// channel was closed
			if !ok {
				return
			}
			a.pushToQueue(record)

		case <-a.ctx.Done():
			return

		}
	}
}

func (a *Auditor) pushToQueue(record *AuditRecord) {
	if a.omitDetails {
		record.Policies = ""
		record.Deciders = ""
	}

	encoded, _ := msgpack.Marshal(record)

	if err := a.que.Push(a.ctx, queueName, encoded); err != nil {
		a.log.Errorw("record audit error", "error", err.Error())
	}
}
