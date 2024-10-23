package worker

type QueueName string

const (
	CriticalQueue QueueName = "critical"
	DefaultQueue  QueueName = "default"
	LowQueue      QueueName = "low"
)

// QueuePriorities maps queue names to their priorities.
var QueuePriorities = map[QueueName]int{
	CriticalQueue: 6,
	DefaultQueue:  3,
	LowQueue:      1,
}
