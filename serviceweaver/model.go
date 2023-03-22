package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"

	openai "github.com/futugyousuzu/go-openai"
)

var (
	addCount = metrics.NewCounter(
		"add_count",
		"The number of times IModel.ListModel has been called",
	)
	addConcurrent = metrics.NewGauge(
		"add_concurrent",
		"The number of concurrent IModel.ListModel calls",
	)
	addSum = metrics.NewHistogram(
		"add_sum",
		"The sums returned by IModel.ListModel count",
		[]float64{1, 10, 100, 1000, 10000},
	)
)

// IModel component.
type IModel interface {
	ListModel(context.Context) (ListModel, error)
}

// Implementation of the Setting component.
type model struct {
	weaver.Implements[IModel]
	weaver.WithConfig[openAIOptions]
}

func (r *model) ListModel(_ context.Context) (ListModel, error) {
	addCount.Add(1.0)
	addConcurrent.Add(1.0)
	defer addConcurrent.Sub(1.0)

	option := r.Config()
	client := openai.NewClient(option.OpenAIKey)
	response := client.ListModels()

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
	addSum.Put(float64(len(items)))
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
