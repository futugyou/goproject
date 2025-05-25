package abstractions

type MemoryQueryResult struct {
	Metadata  MemoryRecordMetadata
	Relevance float64
	Embedding []float32
}

func FromMemoryRecord(record MemoryRecord, relevance float64) MemoryQueryResult {
	return MemoryQueryResult{
		Metadata:  record.Metadata,
		Relevance: relevance,
		Embedding: record.Embedding,
	}
}
