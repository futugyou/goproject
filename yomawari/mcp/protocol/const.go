package protocol

const ErrorCodes_ParseError int = -32700

const ErrorCodes_InvalidRequest int = -32600

const ErrorCodes_MethodNotFound int = -32601

const ErrorCodes_InvalidParams int = -32602

const ErrorCodes_InternalError int = -32603

const NotificationMethods_ToolListChangedNotification string = "notifications/tools/list_changed"

const NotificationMethods_PromptListChangedNotification string = "notifications/prompts/list_changed"

const NotificationMethods_ResourceListChangedNotification string = "notifications/resources/list_changed"

const NotificationMethods_ResourceUpdatedNotification string = "notifications/resources/updated"

const NotificationMethods_RootsUpdatedNotification string = "notifications/roots/list_changed"

const NotificationMethods_LoggingMessageNotification string = "notifications/message"

const NotificationMethods_InitializedNotification string = "notifications/initialized"

const NotificationMethods_ProgressNotification string = "notifications/progress"

const NotificationMethods_CancelledNotification string = "notifications/cancelled"

const RequestMethods_ToolsList string = "tools/list"

const RequestMethods_ToolsCall string = "tools/call"

const RequestMethods_PromptsList string = "prompts/list"

const RequestMethods_PromptsGet string = "prompts/get"

const RequestMethods_ResourcesList string = "resources/list"

const RequestMethods_ResourcesRead string = "resources/read"

const RequestMethods_ResourcesTemplatesList string = "resources/templates/list"

const RequestMethods_ResourcesSubscribe string = "resources/subscribe"

const RequestMethods_ResourcesUnsubscribe string = "resources/unsubscribe"

const RequestMethods_RootsList string = "roots/list"

const RequestMethods_Ping string = "ping"

const RequestMethods_LoggingSetLevel string = "logging/setLevel"

const RequestMethods_CompletionComplete string = "completion/complete"

const RequestMethods_SamplingCreateMessage string = "sampling/createMessage"

const RequestMethods_Initialize string = "initialize"

const RequestMethods_ElicitationCreate string = "elicitation/create"
