package workerpocessor

import (
	"context"
	"log"
	"time"
)

// BaseServerProcessor implement sẵn Start/Stop/Restart
// để các server embed lại
type BaseWorkerProcessor struct {
	processor WorkerProcessor
	cancel    context.CancelFunc
}

func (b *BaseWorkerProcessor) Init(p WorkerProcessor) {
	b.processor = p
}

func (b *BaseWorkerProcessor) Start() error {
	log.Println("Starting Worker...")

	// chạy task trong goroutine riêng
	go func() {
		if err := b.processor.RunningTask(); err != nil {
			log.Printf("Worker stopped with error: %v", err)
		}
	}()
	log.Println("Started Worker!!")
	return nil
}

func (b *BaseWorkerProcessor) Stop() error {
	log.Println("Stopping Worker...")
	// Ở đây base class không biết chi tiết stop,
	// có thể override trong HttpWorker nếu cần shutdown http.Worker
	return nil
}

func (b *BaseWorkerProcessor) Restart() error {
	log.Println("Restarting Worker...")
	if err := b.Stop(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return b.Start()
}
