package abstractions

type FunctionChoiceBehaviorConfigurationContext struct {
	ChatHistory          ChatHistory
	Kernel               *Kernel
	RequestSequenceIndex int
}

type FunctionChoiceBehaviorConfiguration struct {
	Choice     FunctionChoice
	Functions  []KernelFunction
	AutoInvoke bool
	Options    FunctionChoiceBehaviorOptions
}
