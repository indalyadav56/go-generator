package task

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type Task struct {
	Type    string
	Payload []byte
}

type TaskEnqueuer interface {
	Enqueue(taskType string, payload interface{}, queueName string) (*asynq.Task, string, error)
	EnqueueAt(taskType string, payload interface{}, queueName string, processAt time.Time) (*asynq.Task, string, error)
	EnqueueIn(taskType string, payload interface{}, queueName string, delay time.Duration) (*asynq.Task, string, error)
	Close()
}

type AsynqTaskEnqueuer struct {
	client *asynq.Client
}

func NewAsynqTaskEnqueuer(redisAddr, redisUsernmae, redisPassword string) *AsynqTaskEnqueuer {
	redisConnection := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Username: redisUsernmae,
		Password: redisPassword,
	}
	client := asynq.NewClient(redisConnection)
	return &AsynqTaskEnqueuer{client: client}
}

func (e *AsynqTaskEnqueuer) Enqueue(taskType string, payload interface{}, queueName string) (*asynq.Task, string, error) {
	task, err := parseData(taskType, payload)
	if err != nil {
		return nil, "", err
	}
	return e.enqueue(task, queueName, nil)
}

// EnqueueAt creates and enqueues a task to run at a specific time.
func (e *AsynqTaskEnqueuer) EnqueueAt(taskType string, payload interface{}, queueName string, processAt time.Time) (*asynq.Task, string, error) {
	task, err := parseData(taskType, payload)
	if err != nil {
		return nil, "", err
	}
	opts := []asynq.Option{asynq.ProcessAt(processAt)}
	return e.enqueue(task, queueName, opts)
}

// EnqueueIn creates and enqueues a task to run after a delay.
func (e *AsynqTaskEnqueuer) EnqueueIn(taskType string, payload interface{}, queueName string, delay time.Duration) (*asynq.Task, string, error) {
	task, err := parseData(taskType, payload)
	if err != nil {
		return nil, "", err
	}
	opts := []asynq.Option{asynq.ProcessIn(delay)}
	return e.enqueue(task, queueName, opts)
}

// enqueue is a helper function to enqueue a task with optional options.
func (e *AsynqTaskEnqueuer) enqueue(task *Task, queueName string, opts []asynq.Option) (*asynq.Task, string, error) {
	if opts == nil {
		opts = []asynq.Option{}
	}
	opts = append(opts, asynq.Queue(queueName))

	asynqTask := asynq.NewTask(task.Type, task.Payload)
	info, err := e.client.Enqueue(asynqTask, opts...)
	if err != nil {
		return nil, "", err
	}

	return asynqTask, info.ID, nil

}

// parseData is a helper function to create a Task with a JSON-encoded payload.
func parseData(taskType string, payload interface{}) (*Task, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Task{
		Type:    taskType,
		Payload: payloadBytes,
	}, nil
}

func (e *AsynqTaskEnqueuer) Close() {
	e.client.Close()
}
