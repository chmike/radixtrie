# radixtrie
A radix trie package in Go

This is a simple and efficient radix trie in Go. The code is written to be concise and efficient, not the most readable.

This implementation use iteration instead of recursivity. Only Remove needed a stack. The library uses sync.Pool to minimize allocation for the stack.

To use this package

    go get "github/chmike/radixtrie"
    

This simple code show how to insert a key and value, find the value associated to a key, and remove a key and its associated value. 
    
    import "github/chmike/radixtrie"
    
    func main() {
        var r RadixSort
        
        r.Insert("key", "value")
        if val, ok := r.Find("key"); ok {
           // val.(string) == "value"
        }
        r.Remove("key")
    }
    
## A RadixTrie has four methods. 

### Find

The method `Find` lookup a key and return the associated value and a boolean set to true if the key was found. A nil value is thus a possible valid value.

    func (r *RadixTrie) Find(key string) (interface{}, bool)
    
### Insert 

The method `Insert` adds the key and value pair to the RadixTrie, and return the previous value and true if the key already exist in the trie.

    func (r *RadixTrie) Insert(key string, val interface{}) (interface{}, bool)
    
### Remove

The method `Remove`removes a key and its associated value if it exist in the trie. It returns the value and true if the key was found.

    func (r *RadixTrie) Remove(key string) (interface{}, bool)
    
### String

The method `String` returns a multiline representation of the radix trie content. It exposes the trie structure and makes it handy for educational purpose.

    func (r *RadixTrie) String() string
