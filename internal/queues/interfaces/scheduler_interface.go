package interfaces

type SchedulerInterface interface {
	Init()
	Restore() error
	Persist() error
}
