package logging

// Logger interface for transport logging
type Logger interface {
	TransportNotConnected(endpointName string)
	TransportSendingMessage(endpointName, messageID string)
	TransportSendFailed(endpointName, messageID string, err error)
	TransportSentMessage(endpointName, messageID string)
	TransportEnteringReadMessagesLoop(endpointName string)
	TransportWaitingForMessage(endpointName string)
	TransportEndOfStream(endpointName string)
	TransportReceivedMessage(endpointName, message string)
	TransportMessageParseFailed(endpointName, message string, err error)
	TransportReceivedMessageParsed(endpointName, messageID string)
	TransportMessageWritten(endpointName, messageID string)
	TransportReadMessagesCancelled(endpointName string)
	TransportReadMessagesFailed(endpointName string, err error)
	TransportCleaningUp(endpointName string)
	TransportWaitingForReadTask(endpointName string)
	TransportReadTaskCleanedUp(endpointName string)
	TransportCleanupReadTaskTimeout(endpointName string)
	TransportCleanupReadTaskFailed(endpointName string, err error)
	TransportCleanedUp(endpointName string)
}
