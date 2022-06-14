// Package auditor is used to store authorization audit data to the queue.
package auditor

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

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
	que              queue.Queue
	recordsChan      chan *AuditRecord
	poolSize         int
	workerBufferSize int
	flushInterval    time.Duration
	omitDetails      bool
	shouldStop       uint32
	ctx              context.Context
	poolWg           sync.WaitGroup
	log              *logger.Logger
}

var auditor *Auditor

// InitAuditor initializes the global auditor and returns it.
func InitAuditor(ctx context.Context, opts *AuditorOptions, que queue.Queue) *Auditor {
	log := logger.L()
	log.Debugf("building auditor with options: %+v", opts)

	workerBufferSize := opts.BufferSize / opts.PoolSize
	recordsChan := make(chan *AuditRecord, opts.BufferSize)

	auditor = &Auditor{
		que:              que,
		recordsChan:      recordsChan,
		poolSize:         opts.PoolSize,
		workerBufferSize: workerBufferSize,
		flushInterval:    opts.FlushInterval,
		omitDetails:      opts.OmitDetails,
		ctx:              ctx,
		log:              log,
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

	buffer := make([][]byte, 0, a.workerBufferSize)
	ticker := time.NewTicker(a.flushInterval)

	for {
		select {
		case record, ok := <-a.recordsChan:
			// channel was closed
			if !ok {
				a.flushBuffer(buffer)
				return
			}
			buffer = a.appendBuffer(buffer, record)

		case <-ticker.C:
			buffer = a.flushBuffer(buffer)

		case <-a.ctx.Done():
			a.flushBuffer(buffer)
			return

		}
	}
}

func (a *Auditor) flushBuffer(buffer [][]byte) [][]byte {
	if len(buffer) == 0 {
		return buffer
	}

	if err := a.que.PushMany(a.ctx, queueName, buffer); err != nil {
		a.log.Errorw("record audit error", "error", err.Error())
	}

	return buffer[:0]
}

func (a *Auditor) appendBuffer(buffer [][]byte, record *AuditRecord) [][]byte {
	if a.omitDetails {
		record.Policies = ""
		record.Deciders = ""
	}

	encoded, _ := msgpack.Marshal(record)
	buffer = append(buffer, encoded)

	if len(buffer) == a.workerBufferSize {
		buffer = a.flushBuffer(buffer)
	}

	return buffer
}
