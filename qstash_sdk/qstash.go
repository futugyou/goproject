package qstash

import "fmt"

type QstashClient struct {
	http             *httpClient
	common           service
	Message          *MessageService
	URLGroups        *URLGroupsService
	Queues           *QueuesService
	Schedules        *SchedulesService
	SigningKeys      *SigningKeysService
	Events           *EventsService
	DeadLetterQueues *DeadLetterQueuesService
	LLM              *LLMService
	Verify           *VerifyService
	Workflow         *WorkflowService
}

type service struct {
	client *QstashClient
}

const qstash_url string = "https://qstash.upstash.io"

func NewClient(token string) *QstashClient {
	c := &QstashClient{
		http: NewHttpClient(token, qstash_url),
	}
	c.initialize()
	return c
}

func (c *QstashClient) initialize() {
	c.common.client = c
	c.Message = (*MessageService)(&c.common)
	c.URLGroups = (*URLGroupsService)(&c.common)
	c.Queues = (*QueuesService)(&c.common)
	c.Schedules = (*SchedulesService)(&c.common)
	c.SigningKeys = (*SigningKeysService)(&c.common)
	c.Events = (*EventsService)(&c.common)
	c.DeadLetterQueues = (*DeadLetterQueuesService)(&c.common)
	c.LLM = (*LLMService)(&c.common)
	c.Verify = (*VerifyService)(&c.common)
	c.Workflow = (*WorkflowService)(&c.common)
}

type BaseResponse struct {
	Message *string `json:"message,omitempty"`
}

type QstashRequest interface {
	BuilderHeader() map[string]string
	GetPayload() string
}

type QstashHeader struct {
	Method                    *string           `json:"-"`
	Timeout                   *string           `json:"-"`
	Retries                   *int              `json:"-"`
	Forward                   map[string]string `json:"-"`
	Delay                     *string           `json:"-"`
	NotBefore                 *int              `json:"-"`
	DeduplicationId           *string           `json:"-"`
	ContentBasedDeduplication *bool             `json:"-"`
	Callback                  *string           `json:"-"`
	CallbackTimeout           *string           `json:"-"`
	CallbackRetries           *int              `json:"-"`
	CallbackDelay             *string           `json:"-"`
	CallbackMethod            *string           `json:"-"`
	CallbackForward           map[string]string `json:"-"`
	FailureCallback           *string           `json:"-"`
	FailureCallbackTimeout    *string           `json:"-"`
	FailureCallbackRetries    *int              `json:"-"`
	FailureCallbackDelay      *string           `json:"-"`
	FailureCallbackMethod     *string           `json:"-"`
	FailureCallbackForward    map[string]string `json:"-"`
	Cron                      *string           `json:"-"`
	ScheduleId                *string           `json:"-"`
}

func (q *QstashHeader) BuilderHeader() map[string]string {
	header := map[string]string{
		"Upstash-Method":  "POST",
		"Upstash-Timeout": "15m",
		"Upstash-Retries": "3",
	}

	if q == nil {
		return header
	}

	if q.Method != nil {
		header["Upstash-Method"] = *q.Method
	}
	if q.Timeout != nil {
		header["Upstash-Timeout"] = *q.Timeout
	}
	if q.Retries != nil {
		header["Upstash-Retries"] = fmt.Sprintf("%d", *q.Retries)
	}
	if q.Delay != nil {
		header["Upstash-Delay"] = *q.Delay
	}
	if q.NotBefore != nil {
		header["Upstash-Not-Before"] = fmt.Sprintf("%d", *q.NotBefore)
	}
	if q.DeduplicationId != nil {
		header["Upstash-Deduplication-Id"] = *q.DeduplicationId
	}
	if q.ContentBasedDeduplication != nil {
		header["Upstash-Content-Based-Deduplication"] = fmt.Sprintf("%t", *q.ContentBasedDeduplication)
	}
	for key, value := range q.Forward {
		header[fmt.Sprintf("Upstash-Forward-%s", key)] = value
	}
	// Callback
	if q.Callback != nil {
		header["Upstash-Callback"] = *q.Callback
	}
	if q.CallbackTimeout != nil {
		header["Upstash-Callback-Timeout"] = *q.CallbackTimeout
	}
	if q.CallbackRetries != nil {
		header["Upstash-Callback-Retries"] = fmt.Sprintf("%d", *q.CallbackRetries)
	}
	if q.CallbackDelay != nil {
		header["Upstash-Callback-Delay"] = *q.CallbackDelay
	}
	if q.CallbackMethod != nil {
		header["Upstash-Callback-Method"] = *q.CallbackMethod
	}
	for key, value := range q.CallbackForward {
		header[fmt.Sprintf("Upstash-Callback-Forward-%s", key)] = value
	}
	// Failure Callback
	if q.FailureCallback != nil {
		header["Upstash-Failure-Callback"] = *q.FailureCallback
	}
	if q.FailureCallbackTimeout != nil {
		header["Upstash-Failure-Callback-Timeout"] = *q.FailureCallbackTimeout
	}
	if q.FailureCallbackRetries != nil {
		header["Upstash-Failure-Callback-Retries"] = fmt.Sprintf("%d", *q.FailureCallbackRetries)
	}
	if q.FailureCallbackDelay != nil {
		header["Upstash-Failure-Callback-Delay"] = *q.FailureCallbackDelay
	}
	if q.FailureCallbackMethod != nil {
		header["Upstash-Failure-Callback-Method"] = *q.FailureCallbackMethod
	}
	for key, value := range q.FailureCallbackForward {
		header[fmt.Sprintf("Upstash-Failure-Callback-Forward-%s", key)] = value
	}
	//schedules
	if q.Cron != nil {
		header["Upstash-Cron"] = *q.Cron
	}
	if q.ScheduleId != nil {
		header["Upstash-Schedule-Id"] = *q.ScheduleId
	}

	return header
}
