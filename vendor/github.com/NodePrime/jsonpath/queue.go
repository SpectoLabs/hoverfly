package jsonpath

type Results struct {
	nodes []*Result
	head  int
	tail  int
	count int
}

func newResults() *Results {
	return &Results{
		nodes: make([]*Result, 3, 3),
	}
}

func (q *Results) push(n *Result) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*Result, len(q.nodes)*2, len(q.nodes)*2)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

func (q *Results) Pop() *Result {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

func (q *Results) peek() *Result {
	if q.count == 0 {
		return nil
	}
	return q.nodes[q.head]
}

func (q *Results) len() int {
	return q.count
}

func (q *Results) clear() {
	q.head = 0
	q.count = 0
	q.tail = 0
}
