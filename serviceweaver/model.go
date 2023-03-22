package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	openai "github.com/futugyousuzu/go-openai"
)

// Setting component.
type IModel interface {
	ListModel(context.Context) (ListModel, error)
}

// Implementation of the Setting component.
type model struct {
	weaver.Implements[IModel]
}

func (r *model) ListModel(_ context.Context) (ListModel, error) {
	response := &openai.ListModelResponse{}
	items := make([]item, 0)
	for _, data := range response.Datas {
		i := item{
			ID:      data.ID,
			Object:  data.Object,
			Created: data.Created,
			OwnedBy: data.OwnedBy,
			Root:    data.Root,
		}
		items = append(items, i)
	}

	result := ListModel{Datas: items}
	return result, nil
}

type ListModel struct {
	weaver.AutoMarshal
	Datas []item
}
type item struct {
	weaver.AutoMarshal
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int32  `json:"created"`
	OwnedBy string `json:"owned_by"`
	Root    string `json:"root"`
}
