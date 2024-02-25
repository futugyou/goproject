package mongo2struct

import "sync"

type Worker[E any, K any] struct {
	Items []E
	Work  func(E) *K
}

func (w *Worker[E, K]) Process() []K {
	result := make([]K, 0)
	ch := make(chan *K)
	var wg sync.WaitGroup
	for _, c := range w.Items {
		wg.Add(1)
		go w.singleWork(c, &wg, ch)
	}
	go func() {
		for v := range ch {
			if v != nil {
				result = append(result, *v)
			}
		}
	}()

	wg.Wait()
	close(ch)
	return result
}

func (w *Worker[E, K]) singleWork(c E, wg *sync.WaitGroup, ch chan *K) {
	defer wg.Done()
	ch <- w.Work(c)
}
