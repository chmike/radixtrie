package radixtrie

import (
	"testing"
)

func TestInsert(t *testing.T) {
	tests := []struct {
		keyIn  string
		valIn  interface{}
		valOut interface{}
		okOut  bool
	}{
		{"", 1, nil, false},
		{"a", 2, nil, false},
		{"b", 3, nil, false},
		{"ab", 4, nil, false},
		{"long", 5, nil, false},
		{"lo", 6, nil, false},
		{"la", 7, nil, false},
		{"long", 8, 5, true},
		{"long", 9, 8, true},
	}
	var r RadixTrie
	for _, test := range tests {
		valOut, okOut := r.Insert(test.keyIn, test.valIn)
		if okOut != test.okOut {
			t.Errorf("expected ok %v, got %v for %+v", test.okOut, okOut, test)
		}
		if valOut != test.valOut {
			t.Errorf("expected val %v, got %v for %+v", test.valOut, valOut, test)
		}
	}
	t.Log(r)
}

func TestFind(t *testing.T) {
	var r RadixTrie
	if _, ok := r.Find("a"); ok {
		t.Error("expected ok false, got true")
	}
	r.Insert("", 1)
	r.Insert("a", 2)
	r.Insert("b", 3)
	r.Insert("ab", 4)
	r.Insert("long", 5)
	r.Insert("lo", 6)
	r.Insert("la", 7)
	tests := []struct {
		keyIn  string
		valOut interface{}
		okOut  bool
	}{
		{"", 1, true},
		{"a", 2, true},
		{"b", 3, true},
		{"ab", 4, true},
		{"long", 5, true},
		{"lo", 6, true},
		{"la", 7, true},
		{"long", 5, true},
		{"lon", nil, false},
		{"longs", nil, false},
		{"lonb", nil, false},
	}
	for _, test := range tests {
		valOut, okOut := r.Find(test.keyIn)
		if okOut != test.okOut {
			t.Errorf("expected ok %v, got %v for %+v", test.okOut, okOut, test)
		}
		if valOut != test.valOut {
			t.Errorf("expected val %v, got %v for %+v", test.valOut, valOut, test)
		}
	}
	t.Log(r)
}

func TestRemove(t *testing.T) {
	var r RadixTrie
	if _, ok := r.Remove("a"); ok {
		t.Error("expected ok false, got true")
	}
	r.Insert("", 1)
	r.Insert("a", 2)
	r.Insert("b", 3)
	r.Insert("ab", 4)
	r.Insert("long", 5)
	r.Insert("lo", 6)
	r.Insert("la", 7)
	r.Insert("low", 8)
	r.Insert("lowa", 9)
	r.Insert("lowb", 10)

	tests := []struct {
		keyIn  string
		valOut interface{}
		okOut  bool
	}{
		{"", 1, true},
		{"a", 2, true},
		{"b", 3, true},
		{"ab", 4, true},
		{"lon", nil, false},
		{"lonb", nil, false},
		{"l", nil, false},
		{"lz", nil, false},
		{"low", 8, true},
		{"long", 5, true},
		{"lowa", 9, true},
		{"lowb", 10, true},
		{"lo", 6, true},
		{"la", 7, true},
		{"long", nil, false},
		{"longs", nil, false},
	}
	for _, test := range tests {
		valOut, okOut := r.Remove(test.keyIn)
		if okOut != test.okOut {
			t.Errorf("expected ok %v, got %v for %+v", test.okOut, okOut, test)
		}
		if valOut != test.valOut {
			t.Errorf("expected val %v, got %v for %+v", test.valOut, valOut, test)
		}
	}
	t.Log(r)
}

var result bool

func BenchmarkRadixTrie(b *testing.B) {
	var r RadixTrie
	r.Insert("", 1)
	r.Insert("a", 2)
	r.Insert("b", 3)
	r.Insert("ab", 4)
	r.Insert("long", 5)
	r.Insert("lo", 6)
	r.Insert("la", 7)
	r.Insert("very long string", 8)
	//b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, result = r.Find("very long string")
	}
}

func BenchmarkMap(b *testing.B) {
	r := map[string]interface{}{}
	r[""] = 1
	r["a"] = 2
	r["b"] = 3
	r["ab"] = 4
	r["long"] = 5
	r["lo"] = 6
	r["la"] = 7
	r["very long string"] = 8
	//b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, result = r["very long string"]
	}
}
