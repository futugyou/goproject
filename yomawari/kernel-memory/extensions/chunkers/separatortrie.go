package chunkers

type SeparatorTrie struct {
	root *TrieNode
}
type TrieNode struct {
	Children  map[rune]*TrieNode
	Separator *string
}

func NewSeparatorTrie(separators []string) *SeparatorTrie {
	t := &SeparatorTrie{
		root: &TrieNode{
			Children: make(map[rune]*TrieNode),
		},
	}

	for _, v := range separators {
		t.Insert(v)
	}
	return t
}
func (s *SeparatorTrie) Length() int {
	if s == nil || s.root == nil {
		return 0
	}
	return len(s.root.Children)
}

func (s *SeparatorTrie) Insert(separator string) {
	var node = s.root

	for _, v := range separator {
		if _, ok := node.Children[v]; !ok {
			node.Children[v] = &TrieNode{
				Children: make(map[rune]*TrieNode),
			}
		}
		node = node.Children[v]
	}

	node.Separator = &separator
}

func (s *SeparatorTrie) MatchLongest(text string, startIndex int) *string {
	var node = s.root
	var longestMatch *string

	for i := startIndex; i < len(text); i++ {
		if v, ok := node.Children[rune(text[i])]; ok {
			node = v
		} else {
			if node.Separator != nil {
				longestMatch = node.Separator
			}
			break
		}
	}

	return longestMatch
}
