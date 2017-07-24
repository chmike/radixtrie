package radixtrie

import (
	"bytes"
	"fmt"
	"sync"
)

type node struct {
	str   string      // a substring of a key
	val   interface{} // nil when ok is false, otherwise the associated value
	nodes nodes       // sub-nodes holding next key segments
	ok    bool        // true if str is the last substring of a key
}

type nodes []node

// RadixTrie is a radix trie.
// The followinf code shows two ways to insantiate a RadixTrie
//
//   var r RadixTrie
//   var rp = new(RadixTrie)
type RadixTrie struct {
	val   interface{} // nil when ok is false, otherwise the associated value
	nodes nodes       // sub-nodes holding next key segments
	ok    bool        // true if str is the last substring of a key
}

// Find return the value associated to key and true, or nil and false when not found.
func (r *RadixTrie) Find(key string) (interface{}, bool) {
	if len(key) == 0 {
		return r.val, r.ok
	}
	if r.nodes == nil {
		return nil, false
	}
	ns := r.nodes
search:
	for {
		for i := 0; i < len(ns); i++ {
			if key[0] == ns[i].str[0] {
				if len(key) < len(ns[i].str) {
					return nil, false
				}
				for j := 1; j < len(ns[i].str); j++ {
					if key[j] != ns[i].str[j] {
						return nil, false
					}
				}
				if len(key) == len(ns[i].str) {
					return ns[i].val, ns[i].ok
				}
				key = key[len(ns[i].str):]
				ns = ns[i].nodes
				continue search
			}
		}
		return nil, false
	}
}

// Insert add key and value pair to the RadixTrie.
// If a same key already exist in the RadixTrie, replace its value and return
// the previous value and true, otherwize return nil and false.
func (r *RadixTrie) Insert(key string, val interface{}) (interface{}, bool) {
	if len(key) == 0 {
		oldVal, oldOk := r.val, r.ok
		r.val, r.ok = val, true
		return oldVal, oldOk
	}
	if r.nodes == nil {
		r.nodes = make([]node, 0, 4)
	}
	nsp := &r.nodes
insert:
	for {
		ns := *nsp
		// find a node with a first matching char
		for i := 0; i < len(ns); i++ {
			if key[0] == ns[i].str[0] {
				// find common radix length l
				minLen := len(ns[i].str)
				if len(key) < minLen {
					minLen = len(key)
				}
				var l int
				for l = 1; l < minLen; l++ {
					if ns[i].str[l] != key[l] {
						break
					}
				}
				// if ns[i].str is prefix of key
				if l == len(ns[i].str) {
					// if ns[i].str == key
					if l == len(key) {
						oldVal, oldOk := ns[i].val, ns[i].ok
						ns[i].val, ns[i].ok = val, true
						return oldVal, oldOk
					}
					// key is longer than ns[i].str, insert in nodes
					if ns[i].nodes == nil {
						ns[i].nodes = make([]node, 0, 4)
					}
					nsp = &ns[i].nodes
					key = key[l:]
					continue insert
				}
				if l == len(key) {
					// key is radix of ns[i].str
					nodes := make([]node, 1, 4)
					nodes[0] = ns[i]
					nodes[0].str = ns[i].str[l:]
					ns[i] = node{str: key, val: val, ok: true, nodes: nodes}
					return nil, false
				}
				// partial match, insert split node
				nodes := make([]node, 2, 4)
				nodes[0] = ns[i]
				nodes[0].str = ns[i].str[l:]
				nodes[1] = node{str: key[l:], val: val, ok: true}
				ns[i] = node{str: ns[i].str[:l], nodes: nodes}
				return nil, false
			}
		}
		*nsp = append(ns, node{str: key, val: val, ok: true})
		return nil, false
	}
}

type nodeEntry struct {
	nsp *nodes
	i   int
}

var nodeStackPool = sync.Pool{
	New: func() interface{} { return make([]nodeEntry, 0, 50) },
}

// Remove return the value associated to key and true, or nil and false when not found.
func (r *RadixTrie) Remove(key string) (interface{}, bool) {
	if len(key) == 0 {
		oldVal, oldOk := r.val, r.ok
		r.val, r.ok = nil, false
		return oldVal, oldOk
	}
	if r.nodes == nil {
		return nil, false
	}
	nodeStack := nodeStackPool.Get().([]nodeEntry)[:0]
	defer nodeStackPool.Put(nodeStack)
	nsp := &r.nodes
remove:
	for {
		ns := *nsp
		for i := 0; i < len(ns); i++ {
			if key[0] == ns[i].str[0] {
				if len(key) < len(ns[i].str) {
					return nil, false
				}
				for j := 1; j < len(ns[i].str); j++ {
					if key[j] != ns[i].str[j] {
						return nil, false
					}
				}
				if len(key) == len(ns[i].str) {
					if !ns[i].ok {
						return nil, false
					}
					// we found the value to delete
					oldVal := ns[i].val
					ns[i].val, ns[i].ok = nil, false
					for {
						nodesLen := len(ns[i].nodes)
						if nodesLen > 1 {
							return oldVal, true
						}
						if nodesLen == 1 {
							ns[i].nodes[0].str = ns[i].str + ns[i].nodes[0].str
							ns[i] = ns[i].nodes[0]
							return oldVal, true
						}
						// remove node i from ns
						if len(ns) == 1 {
							*nsp = nil
						} else {
							last := len(ns) - 1
							if i < last {
								ns[i] = ns[last]
							}
							*nsp = (*nsp)[:last]
						}
						if len(nodeStack) == 0 {
							return oldVal, true
						}
						last := len(nodeStack) - 1
						nsp, i = nodeStack[last].nsp, nodeStack[last].i
						nodeStack = nodeStack[:last]
						ns = *nsp
						if ns[i].ok {
							return oldVal, true
						}
					}
				}
				nodeStack = append(nodeStack, nodeEntry{nsp, i})
				key = key[len(ns[i].str):]
				nsp = &ns[i].nodes
				continue remove
			}
		}
		return nil, false
	}
}

// Convert the RadixTree into a multiline string.
func (r *RadixTrie) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, "RadixTrie:\n")
	if r.ok {
		fmt.Fprintf(buf, "%q val: %v\n", "", r.val)
	}
	if r.nodes != nil {
		r.nodes.sdump(buf, 0)
	}
	return buf.String()
}

func (ns nodes) sdump(buf *bytes.Buffer, depth int) {
	for i := range ns {
		for j := 0; j < depth; j++ {
			fmt.Fprint(buf, "|   ")
		}
		fmt.Fprintf(buf, "[%d] %q", i, ns[i].str)
		if ns[i].ok {
			fmt.Fprintf(buf, " val: %v", ns[i].val)
		}
		fmt.Fprint(buf, "\n")
		if ns[i].nodes != nil {
			ns[i].nodes.sdump(buf, depth+1)
		}
	}
}
