package queues

import (
	"github.com/roylee0704/gron"
	"queue-manager/internal/providers"
	"queue-manager/internal/queues/interfaces"
	"queue-manager/internal/structures"
	"time"
)

type Scheduler struct {
	io     interfaces.QueueIOManagerInterface
	qm     interfaces.QueueManagerInterface
	config *structures.Config
	logger providers.Logger
}

func (s *Scheduler) Init() {
	c := gron.New()
	interval := s.config.Persistence.SaveInterval
	c.AddFunc(gron.Every(interval*time.Second), func() {
		err := s.io.SaveQueuesToFile(s.qm, s.config.Persistence.FilePath)
		if err != nil {
			s.logger.Errorf(providers.TypeApp, "Error while persisting data: %s", err)
			return
		}
		s.logger.Infof(providers.TypeApp, "Persisted data to file %s", s.config.Persistence.FilePath)
	})
	c.Start()
}

func (s *Scheduler) Restore() error {
	err := s.io.LoadQueuesFromFile(s.qm, s.config.Persistence.FilePath)
	if err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) Persist() error {
	s.logger.Infof(providers.TypeApp, "Persisting queues to file...")
	err := s.io.SaveQueuesToFile(s.qm, s.config.Persistence.FilePath)
	if err != nil {
		s.logger.Errorf(providers.TypeApp, "Error while persisting data: %s", err)
		return err
	}
	return nil
}

func NewScheduler(io interfaces.QueueIOManagerInterface, qm interfaces.QueueManagerInterface, config *structures.Config, logger providers.Logger) interfaces.SchedulerInterface {
	return &Scheduler{
		io:     io,
		qm:     qm,
		config: config,
		logger: logger,
	}
}
