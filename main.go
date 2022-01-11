package main

import "fmt"

func main() {
	arrStr := IteratorOut{assemblyItem{Value: "a"}, assemblyItem{Value: "b"}, assemblyItem{Value: "c"}, assemblyItem{Value: "d"}}
	arrInt := IteratorOut{assemblyItem{Value: "1"}, assemblyItem{Value: "2"}, assemblyItem{Value: "3"}}
	arrChar := IteratorOut{assemblyItem{Value: "一"}, assemblyItem{Value: "二"}, assemblyItem{Value: "三"}}
	iter := NewSliceIterator(arrStr).ToAssemblyIterator().BuildNewAssemblyIterator(NewSliceIterator(arrInt)).BuildNewAssemblyIterator(NewSliceIterator(arrChar))
	for {
		if value, ok := iter.Next(); ok {
			fmt.Println(value)
		} else {
			break
		}
	}
}

type IteratorOut []assemblyItem //唯一需要变的就是这里，出了泛型就不用变了

type assemblyItem struct {
	Value  string
	Type   string
	Anyhow interface{}
}

type SliceIterator struct {
	origin IteratorOut
	len    int
	index  int
}

func NewSliceIterator(origin IteratorOut) *SliceIterator {
	return &SliceIterator{
		origin: origin,
		index:  0,
		len:    len(origin),
	}
}

func (n *SliceIterator) Next() (IteratorOut, bool) {
	if n.index < n.len {
		n.index++
		return IteratorOut{n.origin[n.index-1]}, true
	}
	return nil, false
}

func (n *SliceIterator) Reset() {
	n.index = 0
}

func (n *SliceIterator) ToAssemblyIterator() *AssemblyIterator {
	iter := new(AssemblyIterator)
	iter.isInit = true
	iter.IteratorA = n
	iter.valueB, iter.statusB = IteratorOut{}, false
	iter.valueA, iter.statusA = n.Next()
	return iter
}

type Iterator interface {
	Next() (IteratorOut, bool)
	Reset()
}

type AssemblyIterator struct {
	IteratorA Iterator
	IteratorB Iterator
	valueA    IteratorOut
	valueB    IteratorOut
	statusA   bool
	statusB   bool
	isInit    bool
}

func (u *AssemblyIterator) Reset() {
	if u.IteratorA != nil {
		u.IteratorA.Reset()
	}
	if u.IteratorB != nil {
		u.IteratorB.Reset()
	}
}

func (u *AssemblyIterator) BuildNewAssemblyIterator(next Iterator) *AssemblyIterator {
	item := &AssemblyIterator{
		IteratorA: u,
		IteratorB: next,
		valueA:    nil,
		valueB:    nil,
		statusA:   true,
		statusB:   true,
		isInit:    true,
	}
	item.valueA, item.statusA = item.IteratorA.Next()
	item.valueB, item.statusB = item.IteratorB.Next()
	return item
}

func (u *AssemblyIterator) Next() (IteratorOut, bool) {
	if (!u.statusA && !u.statusB) || (u.valueA == nil || u.valueB == nil) {
		return nil, false
	}
	if u.isInit {
		u.isInit = false
		if u.valueA == nil || u.valueB == nil {
			return nil, false
		}
		return append(u.valueA, u.valueB...), true
	} else {
		for {
			var tmp IteratorOut
			if u.statusB {
				if tmp, u.statusB = u.IteratorB.Next(); u.statusB {
					u.valueB = tmp
					return append(u.valueA, u.valueB...), true
				}
			}
			if u.IteratorB != nil {
				u.IteratorB.Reset()
				u.statusB = true
			}
			if u.statusA {
				if tmp, u.statusA = u.IteratorA.Next(); u.statusA {
					u.valueA = tmp
					if u.IteratorB != nil {
						continue
					}
					return append(u.valueA, u.valueB...), true
				}
			}
			return nil, false
		}
	}
}
