package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/core/filesystem"
	"github.com/google/uuid"
)

// ConnectToQueue implements pipeline.IQueue.
func (q *SimpleQueues) ConnectToQueue(ctx context.Context, queueName string, options pipeline.QueueOptions) (pipeline.IQueue, error) {
	if queueName == "" {
		return nil, errors.New("queue name cannot be empty")
	}

	if q.queueName == queueName {
		return q, nil
	}

	if q.queueName != "" {
		return nil, fmt.Errorf("queue is already connected to %s", q.queueName)
	}

	q.queueName = queueName
	q.poisonQueueName = queueName + q.config.PoisonQueueSuffix

	if err := q.createDirectories(ctx); err != nil {
		return nil, err
	}
	if options.DequeueEnabled {
		q.populateTicker = time.NewTicker(time.Duration(q.config.PollDelayMsecs) * time.Second)
		q.dispatchTicker = time.NewTicker(time.Duration(q.config.DispatchFrequencyMsecs) * time.Second)

		q.wg.Add(2)
		go q.populateQueue(ctx)
		go q.dispatchMessages(ctx)
	}

	return q, nil
}

// Enqueue implements pipeline.IQueue.
func (q *SimpleQueues) Enqueue(ctx context.Context, message string) error {
	messageID := time.Now().Format("20060102.150405.0000000") + "." + uuid.New().String()

	msg := &Message{
		Id:           messageID,
		Content:      message,
		DequeueCount: 0,
		Schedule:     time.Now(),
	}

	return q.storeMessage(ctx, q.queueName, msg)
}

// OnDequeue implements pipeline.IQueue.
func (q *SimpleQueues) OnDequeue(ctx context.Context, processMessageAction func(ctx context.Context, message string) pipeline.ReturnType) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.receivedHandlers = append(q.receivedHandlers, processMessageAction)
	return nil
}

const (
	fileExt             = ".sqm.json"
	defaultPollDelay    = 100 * time.Millisecond
	defaultDispatchFreq = 100 * time.Millisecond
	defaultLockSeconds  = 30
	defaultMaxAttempts  = 5
)

// SimpleQueues implements a file-based queue
type SimpleQueues struct {
	config           SimpleQueuesConfig
	queueName        string
	poisonQueueName  string
	maxAttempts      int
	queue            chan Message
	receivedHandlers []func(ctx context.Context, message string) pipeline.ReturnType
	filesystem       filesystem.IFileSystem
	populateTicker   *time.Ticker
	dispatchTicker   *time.Ticker
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	mu               sync.Mutex
}

// NewSimpleQueues creates a new file-based queue
func NewSimpleQueues(config SimpleQueuesConfig) (*SimpleQueues, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	var fs filesystem.IFileSystem
	switch config.StorageType {
	case "Disk":
		fs = filesystem.NewDiskFileSystem(config.Directory, nil)
	case "Volatile":
		fs = filesystem.GetVolatileFileSystemInstance(config.Directory, nil)
	default:
		return nil, fmt.Errorf("unknown storage type %s", config.StorageType)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &SimpleQueues{
		config:      config,
		maxAttempts: config.MaxRetriesBeforePoisonQueue + 1,
		queue:       make(chan Message, config.FetchBatchSize),
		filesystem:  fs,
		ctx:         ctx,
		cancel:      cancel,
	}, nil
}

// Close stops the queue and cleans up resources
func (q *SimpleQueues) Close() error {
	q.cancel()
	if q.populateTicker != nil {
		q.populateTicker.Stop()
	}
	if q.dispatchTicker != nil {
		q.dispatchTicker.Stop()
	}
	q.wg.Wait()
	close(q.queue)
	return nil
}

func (q *SimpleQueues) populateQueue(ctx context.Context) {
	defer q.wg.Done()

	for {
		select {
		case <-q.ctx.Done():
			return
		case <-q.populateTicker.C:
			if len(q.queue) >= q.config.FetchBatchSize {
				continue
			}

			files, err := q.filesystem.GetAllFileNames(q.ctx, q.queueName, "")
			if err != nil {
				if os.IsNotExist(err) {
					if err := q.createDirectories(q.ctx); err != nil {
						continue
					}
				}
				continue
			}

			for _, fileName := range files {
				if len(q.queue) >= q.config.FetchBatchSize {
					break
				}

				if !strings.HasSuffix(fileName, fileExt) {
					continue
				}

				messageID := strings.TrimSuffix(fileName, fileExt)
				msg, err := q.readMessage(ctx, messageID)
				if err != nil {
					continue
				}

				if msg.IsTimeToRun() && !msg.IsLocked() {
					msg.Lock(q.config.FetchLockSeconds)
					msg.DequeueCount++
					if err := q.storeMessage(ctx, q.queueName, msg); err != nil {
						continue
					}
					q.queue <- *msg
				}
			}
		}
	}
}

func (q *SimpleQueues) dispatchMessages(ctx context.Context) {
	defer q.wg.Done()

	for {
		select {
		case <-q.ctx.Done():
			return
		case <-q.dispatchTicker.C:
			if len(q.queue) == 0 {
				continue
			}

			msg := <-q.queue
			q.processMessage(ctx, &msg)
		}
	}
}

func (q *SimpleQueues) storeMessage(ctx context.Context, queueName string, msg *Message) error {
	d, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return q.filesystem.WriteFile(ctx, queueName, "", msg.Id+fileExt, io.NopCloser(bytes.NewReader(d)))
}

func (q *SimpleQueues) readMessage(ctx context.Context, id string) (*Message, error) {
	serializedMsg, err := q.filesystem.ReadFileAsText(ctx, q.queueName, "", id+fileExt)
	if err != nil {
		return nil, err
	}
	var msg Message
	err = json.Unmarshal([]byte(*serializedMsg), &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (q *SimpleQueues) createDirectories(ctx context.Context) error {
	err := q.filesystem.CreateVolume(ctx, q.queueName)
	if err != nil {
		return err
	}
	return q.filesystem.CreateVolume(ctx, q.poisonQueueName)
}

func (q *SimpleQueues) deleteMessage(ctx context.Context, id string) error {
	var fileName = id + fileExt
	err := q.filesystem.DeleteFile(ctx, q.queueName, "", fileName)
	if err != nil {
		return q.createDirectories(ctx)
	}
	return nil
}

func (q *SimpleQueues) processMessage(ctx context.Context, msg *Message) {
	for _, received := range q.receivedHandlers {
		var retry = false
		var poison = false
		returnType := received(ctx, msg.Content)
		switch returnType {
		case pipeline.ReturnTypeSuccess:
			q.deleteMessage(ctx, msg.Id)
		case pipeline.ReturnTypeTransientError:
			msg.LastError = "Message handler returned false"
			if msg.DequeueCount == q.maxAttempts {
				poison = true
			} else {
				retry = true
			}
		case pipeline.ReturnTypeFatalError:
			poison = true
		}

		msg.Unlock()
		if retry {
			var backoffDelay = time.Duration(1 * msg.DequeueCount)
			msg.RunIn(backoffDelay)
			q.storeMessage(q.ctx, q.queueName, msg)
		} else if poison {
			q.storeMessage(q.ctx, q.poisonQueueName, msg)
			q.deleteMessage(q.ctx, msg.Id)
		}
	}
}
