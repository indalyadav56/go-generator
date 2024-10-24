package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type AsynqServer struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewAsynqServer(redisAddr, redisUsername, redisPassword string, concurrency int) *AsynqServer {
	redisConnection := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
	}

	server := asynq.NewServer(
		redisConnection,
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				string(CriticalQueue): QueuePriorities[CriticalQueue],
				string(DefaultQueue):  QueuePriorities[DefaultQueue],
				string(LowQueue):      QueuePriorities[LowQueue],
			},
		},
	)

	mux := asynq.NewServeMux()

	return &AsynqServer{
		server: server,
		mux:    mux,
	}
}

func (s *AsynqServer) HandleFunc(taskType string, handler func(ctx context.Context, t *asynq.Task) error) {
	s.mux.HandleFunc(taskType, handler)
}

func (s *AsynqServer) Run() error {
	if err := s.server.Run(s.mux); err != nil {
		return err
	}
	return nil
}

func (s *AsynqServer) Stop() {
	s.server.Stop()
}
