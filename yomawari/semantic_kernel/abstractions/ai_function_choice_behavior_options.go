package abstractions

type FunctionChoiceBehaviorOptions struct {
	AllowParallelCalls         bool `json:"allow_parallel_calls"`
	AllowConcurrentInvocation  bool `json:"allow_concurrent_invocation"`
	AllowStrictSchemaAdherence bool `json:"allow_strict_schema_adherence"`
	RetainArgumentTypes        bool `json:"-"`
}
