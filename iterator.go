package stream

import "github.com/youthlin/stream/types"

const unknownSize = -1

type iterator interface {
	GetSizeIfKnown() int64
	HasNext() bool
	Next() types.T
}

func it(elements ...types.T) iterator {
	return &sliceIterator{
		elements: elements,
		current:  0,
	}
}
func withSeed(seed types.T, f types.UnaryOperator) iterator {
	return &seedIt{
		element:  seed,
		operator: f,
		first:    true,
	}
}
func withSupplier(get types.Supplier) iterator {
	return &supplierIt{get: get}
}
func withRange(fromInclude, toExclude endpoint, step int) iterator {
	return &rangeIt{
		from: fromInclude,
		to:   toExclude,
		step: step,
		next: fromInclude,
	}
}

// region sliceIterator

type sliceIterator struct {
	elements []types.T
	current  int
}

func (s *sliceIterator) GetSizeIfKnown() int64 {
	return int64(len(s.elements))
}

func (s *sliceIterator) HasNext() bool {
	return s.current < len(s.elements)
}

func (s *sliceIterator) Next() types.T {
	e := s.elements[s.current]
	s.current++
	return e
}

// endregion sliceIterator

// region seedIt

type seedIt struct {
	element  types.T
	operator types.UnaryOperator
	first    bool
}

func (s *seedIt) GetSizeIfKnown() int64 {
	return unknownSize
}

func (s *seedIt) HasNext() bool {
	return true
}

func (s *seedIt) Next() types.T {
	if s.first {
		s.first = false
		return s.element
	}
	s.element = s.operator(s.element)
	return s.element
}

// endregion seedIt

// region supplierIt

type supplierIt struct {
	get types.Supplier
}

func (s *supplierIt) GetSizeIfKnown() int64 {
	return unknownSize
}

func (s *supplierIt) HasNext() bool {
	return true
}

func (s *supplierIt) Next() types.T {
	return s.get()
}

// endregion supplierIt

// region rangeIt

type rangeIt struct {
	from endpoint
	to   endpoint
	step int
	next endpoint
}

func (r *rangeIt) GetSizeIfKnown() int64 {
	return unknownSize
}

func (r *rangeIt) HasNext() bool {
	if r.step >= 0 {
		return r.next.CompareTo(r.to) < 0
	}
	return r.next.CompareTo(r.to) > 0
}

func (r *rangeIt) Next() types.T {
	curr := r.next
	r.next = curr.Add(r.step)
	return curr
}

// endregion rangeIt

// region Sortable

// Sortable use types.Comparator to sort []types.T 可以使用指定的 cmp 比较器对 list 进行排序
// see sort.Interface
type Sortable struct {
	List []types.T
	Cmp  types.Comparator
}

func (a *Sortable) Len() int {
	return len(a.List)
}

func (a *Sortable) Less(i, j int) bool {
	return a.Cmp(a.List[i], a.List[j]) < 0
}

func (a *Sortable) Swap(i, j int) {
	a.List[i], a.List[j] = a.List[j], a.List[i]
}

// endregion Sortable

// region endpoint

// endpoint used in rangeIt.
type endpoint interface {
	CompareTo(other endpoint) int
	Add(step int) endpoint
}

type epInt int

func (m epInt) CompareTo(other endpoint) int {
	return int(m - other.(epInt))
}

func (m epInt) Add(step int) endpoint {
	return m + epInt(step)
}

type epInt64 int64

func (m epInt64) CompareTo(other endpoint) int {
	return int(m - other.(epInt64))
}

func (m epInt64) Add(step int) endpoint {
	return m + epInt64(step)
}

// endregion endpoint
