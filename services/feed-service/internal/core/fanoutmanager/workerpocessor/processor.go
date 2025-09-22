package workerpocessor

type WorkerProcessor interface {
	Start() error
	Stop() error
	Restart() error
	RunningTask() error
}
