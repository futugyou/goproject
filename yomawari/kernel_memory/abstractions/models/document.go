package models

import (
	"io"
	"time"

	"github.com/google/uuid"
)

type Document struct {
	Id    string
	Files *FileCollection
	Tags  *TagCollection
}

func NewDocument(id *string, tags *TagCollection, filePaths []string) *Document {
	var docID string
	if id != nil {
		docID = *id
	} else {
		docID = RandomId()
	}
	if tags == nil {
		tags = &TagCollection{}
	}

	f := NewFileCollection()
	for _, v := range filePaths {
		f.AddFile(v)
	}

	return &Document{
		Id:    docID,
		Files: f,
		Tags:  tags,
	}
}

func (doc *Document) AddTag(name string, value string) {
	if doc == nil {
		return
	}
	doc.Tags.AddOrAppend(name, value)
}

func (doc *Document) AddFile(filePath string) {
	if doc == nil {
		return
	}
	doc.Files.AddFile(filePath)
}

func (doc *Document) AddFiles(filePaths []string) {
	if doc == nil {
		return
	}
	for _, filePath := range filePaths {
		doc.Files.AddFile(filePath)
	}
}

func (doc *Document) AddStream(fileName string, content io.ReadCloser) {
	if doc == nil {
		return
	}
	doc.Files.AddStream(fileName, content)
}

func RandomId() string {
	id := uuid.New().String()
	currentTime := time.Now().Format("20060102150405.9999999")
	return id + currentTime
}
