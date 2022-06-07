// Package auditor is used to store authorization audit data to the queue (redis list).
package auditor

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/che-kwas/iam-kit/logger"
	"gopkg.in/vmihailenco/msgpack.v2"

	"iam-authz/internal/pkg/redis"
)

const (
	queueName       = "iam-authz-audit"
	queueExpiration = time.Duration(24 * time.Hour)
)

// AuditRecord defines the details of a authorization request.
type AuditRecord struct {
	TimeStamp  int64     `json:"timestamp"`
	Username   string    `json:"username"`
	Effect     string    `json:"effect"`
	Conclusion string    `json:"conclusion"`
	Request    string    `json:"request"`
	Policies   string    `json:"policies"`
	Deciders   string    `json:"deciders"`
	ExpireAt   time.Time `json:"expireAt"`
}

// Auditor defines the structure of an auditor.
type Auditor struct {
	recordsChan      chan *AuditRecord
	poolSize         int
	workerBufferSize int
	flushInterval    time.Duration
	shouldStop       uint32
	ctx              context.Context
	poolWg           sync.WaitGroup
	log              *logger.Logger
}

var auditor *Auditor

// InitAuditor initializes the global auditor and returns it.
func InitAuditor(ctx context.Context, opts *AuditorOptions) *Auditor {
	log := logger.L()
	log.Debugf("building auditor with options: %+v", opts)

	workerBufferSize := opts.BufferSize / opts.PoolSize
	recordsChan := make(chan *AuditRecord, opts.BufferSize)

	auditor = &Auditor{
		recordsChan:      recordsChan,
		poolSize:         opts.PoolSize,
		workerBufferSize: workerBufferSize,
		flushInterval:    opts.FlushInterval,
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

// RecordHit stores an AuditRecord in redis.
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

			encoded, _ := msgpack.Marshal(record)
			buffer = append(buffer, encoded)

			if len(buffer) == a.workerBufferSize {
				buffer = a.flushBuffer(buffer)
			}

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

	pipe := redis.Client().Pipeline()
	for _, record := range buffer {
		pipe.RPush(a.ctx, queueName, record)
	}
	pipe.Expire(a.ctx, queueName, queueExpiration)

	if _, err := pipe.Exec(a.ctx); err != nil {
		a.log.Errorw("record audit error", "error", err.Error())
	}

	return buffer[:0]
}
