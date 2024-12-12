package scheduler

import (
	"github.com/hibiken/asynq"
)

type AsynqTaskScheduler struct {
	scheduler *asynq.Scheduler
}

func NewAsynqScheduler(redisAddr, redisPassword string) *AsynqTaskScheduler {
	redisConnection := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
	}
	scheduler := asynq.NewScheduler(redisConnection, &asynq.SchedulerOpts{})
	return &AsynqTaskScheduler{scheduler: scheduler}
}

func (s *AsynqTaskScheduler) Register(schedule string, task *asynq.Task) (string, error) {
	entryID, err := s.scheduler.Register(schedule, task)
	if err != nil {
		return "", err
	}

	return entryID, nil
}

func (s *AsynqTaskScheduler) Run() error {
	return s.scheduler.Run()
}

func (s *AsynqTaskScheduler) Shutdown() {
	s.scheduler.Shutdown()
}
