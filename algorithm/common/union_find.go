package common

type UnionFind struct {
	size   []int
	parent []int
	count  int
}

func NewUnionFind(n int) *UnionFind {
	uf := UnionFind{}
	uf.count = n
	uf.size = make([]int, n)
	uf.parent = make([]int, n)
	for i := 0; i < n; i++ {
		uf.size[i] = 1
		uf.parent[i] = i
	}
	return &uf
}

func (u *UnionFind) Count() int {
	return u.count
}

func (u *UnionFind) findRoot(x int) int {
	for {
		if x == u.parent[x] {
			return x
		}
		u.parent[x] = u.parent[u.parent[x]]
		x = u.parent[x]
	}
}

func (u *UnionFind) Union(p, q int) {
	pRoot := u.findRoot(p)
	qRoot := u.findRoot(q)
	if pRoot == qRoot {
		return
	}
	// 小树接到大树下，较平衡
	if u.size[pRoot] > u.size[qRoot] {
		u.parent[qRoot] = pRoot
		u.size[pRoot] += u.size[qRoot]
	} else {
		u.parent[pRoot] = qRoot
		u.size[qRoot] += u.size[pRoot]
	}
	u.count--
}

func (u *UnionFind) Connected(p, q int) bool {
	pRoot := u.findRoot(p)
	qRoot := u.findRoot(q)
	return pRoot == qRoot
}
