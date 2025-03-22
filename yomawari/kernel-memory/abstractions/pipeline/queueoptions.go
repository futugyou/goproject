package pipeline

type QueueOptions struct {
	DequeueEnabled bool
}

var (
	PubSub      = QueueOptions{DequeueEnabled: true}
	PublishOnly = QueueOptions{DequeueEnabled: false}
)

func (q QueueOptions) Equals(other QueueOptions) bool {
	return q.DequeueEnabled == other.DequeueEnabled
}

func (q QueueOptions) HashCode() int {
	if q.DequeueEnabled {
		return 1
	}
	return 2
}
