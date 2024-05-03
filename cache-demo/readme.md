# 第五章： 进程内缓存

## 简介

本章手写了fifo/lfu/lru，BigCache库的基本使用，基准测试及一些优化方法，比如以分片加速并行(减少锁的使用)，
和避免GC(在栈而不是堆上分配资源)

## fifo

<details>
<summary> 使用了标准库container/list双向链表，set时如果数据存在则调用MoveToBack()移到最后，
不存在则调用PushBack()添加到最后。 删除最老数据时先调用Front()再remove </summary>

```golang
func (f *fifo) Set(key string, value interface{}) {
 if e, ok := f.cache[key]; ok {
  f.ll.MoveToBack(e)
  en := e.Value.(*entry)
  f.usedBytes = f.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
  en.value = value
  return
 }
 en := &entry{key, value}
 e := f.ll.PushBack(en)
 f.cache[key] = e

 f.usedBytes += en.Len()
 if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
  f.DelOldest()
 }
}

func (f *fifo) DelOldest() {
 f.removeElement(f.ll.Front())
}
```

</details>

## lfu

lfu实现中需要用到最小堆算法（权重小的在前，删除时直接pop），所以选择了标准库container/heap。

<details>
<summary>
heap需要一个叫Interface的interface。
sort.Interface也是个叫Interface的interface。
这名字起的，也就能骗骗大小写不敏感的我了，golang标准库里面应该有很多叫这名字。 </summary>

```golang
type Interface interface {
 sort.Interface
 Push(x interface{}) // add x as element Len()
 Pop() interface{}   // remove and return element Len() - 1.
}

type Interface interface {
 // Len is the number of elements in the collection.
 Len() int
 // Less reports whether the element with
 // index i should sort before the element with index j.
 Less(i, j int) bool
 // Swap swaps the elements with indexes i and j.
 Swap(i, j int)
}

```

</details>

<details>
<summary> 来实现这个叫Interface的interface </summary>

```golang
type entry struct {
 keystring
 value  interface{}
 weight int
 index  int
}

func (e *entry) Len() int {
 return cache.CalcLen(e.value) + 4 + 4
}

type queue []*entry

func (q queue) Len() int {
 return len(q)
}
// '<'换成'>'就是最大堆
func (q queue) Less(i, j int) bool {
 return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
 q[i], q[j] = q[j], q[i]
 q[i].index = i
 q[j].index = j
}

func (q *queue) Push(x interface{}) {
 n := len(*q)
 en := x.(*entry)
 en.index = n
 *q = append(*q, en)
}

func (q *queue) Pop() interface{} {
 old := *q
 n := len(old)
 en := old[n-1]
 old[n-1] = nil
 en.index = -1
 *q = old[0 : n-1]
 return en
}
```

</details>

<details>
<summary> 最后来实现lfu，增删改查会影响权重weight所以都会涉及到重排，需要调用heap.Fix </summary>

```golang
// 重排
func (q *queue) update(en *entry, value interface{}, weight int) {
 en.value = value
 en.weight = weight
 heap.Fix(q, en.index)
}

func (l *lfu) Set(key string, value interface{}) {
 if e, ok := l.cache[key]; ok {
  l.usedBytes = l.usedBytes - cache.CalcLen(e.value) + cache.CalcLen(value)
  l.queue.update(e, value, e.weight+1)
  return
 }
 en := &entry{key: key, value: value}
 heap.Push(l.queue, en)
 l.cache[key] = en

 l.usedBytes += en.Len()
 if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
  l.removeElement(heap.Pop(l.queue))
 }
}

func (l *lfu) Get(key string) interface{} {
 if e, ok := l.cache[key]; ok {
  l.queue.update(e, e.value, e.weight+1)
  return e.value
 }
 return nil
}
```

</details>

## lru

<details>
<summary> lru和fifo一样实现，只用Get时多一步MoveToBack()的操作  </summary>

```golang
func (f *lru) Get(key string) interface{} {
 if e, ok := f.cache[key]; ok {
  f.ll.MoveToBack(e)
  return e.Value.(*entry).value
 }
 return nil
}

```

</details>

## 支持并发读写

重新包装了cache结构，加入了sync.RWMutex，这个就不细说了。

## bigcache

```golang
go get -u github.com/allegro/bigcache/v2
```

<details>
<summary> 很简单的code  </summary>

```golang
package main

import (
 "log"
 "time"

 "github.com/allegro/bigcache/v2"
)

func main() {
 cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(20 * time.Minute))
 if err != nil {
  log.Printf("one %v", err)
  return
 }

 // if key not exsits ,there will be an Entry not found error
 entry, err := cache.Get("my-unique-key")
 if err != nil && entry != nil {
  log.Printf("two %v", err)
  return
 }

 if entry == nil {
  entry = []byte("value")
  cache.Set("my-unique-key", entry)
 }
 log.Println(string(entry))
}

```

</details>

## 按照bigcache的分片思想写一个cache

将原先支持并发的cache封装一下作为一个cacheShard，实现一个hash算法，创建一个对外调用的Cache，内部的增删改查会先通过hash得到一个cacheShard，最终调用cacheShard响应的方法

<details>
<summary> 简单的hash算法  </summary>

```golang
package fast

func newDefaultHasher() fnv64a {
 return fnv64a{}
}

type fnv64a struct{}

const (
 offest64 = 14695981039346656037
 prime64  = 1099544628211
)

func (f fnv64a) Sum64(key string) uint64 {
 var hash uint64 = offest64
 for i := 0; i < len(key); i++ {
  hash ^= uint64(key[i])
  hash *= prime64
 }
 return hash
}

```

</details>

<details>
<summary> 对外Cache的实现  </summary>

```golang
package fast

type fastCache struct {
 shards[]*cacheShard
 shardMask uint64
 hash  fnv64a
}

func NewFastCahe(maxEntries int, shardsNum int, onEvicted func(key string, value interface{})) *fastCache {
 fastCache := &fastCache{
  hash:  newDefaultHasher(),
  shards:make([]*cacheShard, shardsNum),
  shardMask: uint64(shardsNum - 1),
 }

 for i := 0; i < shardsNum; i++ {
  fastCache.shards[i] = newCacheShard(maxEntries, onEvicted)
 }
 return fastCache
}

func (c *fastCache) getShard(key string) *cacheShard {
 hashkey := c.hash.Sum64(key)
 return c.shards[hashkey&c.shardMask]
}

func (c *fastCache) Set(key string, value interface{}) {
 c.getShard(key).set(key, value)
}

func (c *fastCache) Get(key string) interface{} {
 return c.getShard(key).get(key)
}

func (c *fastCache) Del(key string) {
 c.getShard(key).del(key)
}

func (c *fastCache) Len() int {
 l := 0
 for _, r := range c.shards {
  l += r.len()
 }
 return l
}

```

</details>
