package dag

import "sync"

type DAGTaskSimple struct {
	nodes     map[string][]string
	methods   map[string]func()
	indegrees map[string]int
	queue     []string
	lock      sync.Mutex
}

func NewDAGTaskSimple() *DAGTaskSimple {
	return &DAGTaskSimple{
		nodes:     make(map[string][]string),
		methods:   make(map[string]func()),
		indegrees: make(map[string]int),
		queue:     []string{},
	}
}

func (d *DAGTaskSimple) AddNode(node string, method func()) {
	_, ok := d.methods[node]
	if !ok {
		d.indegrees[node] = 0
	}

	d.methods[node] = method
}

func (d *DAGTaskSimple) AddEdge(from, to string) {
	d.nodes[from] = append(d.nodes[from], to)
	d.indegrees[to]++
}

func (d *DAGTaskSimple) TopologicalSort() {
	for node, indegree := range d.indegrees {
		if indegree == 0 {
			d.queue = append(d.queue, node)
		}
	}

	for len(d.queue) > 0 {
		wg := sync.WaitGroup{}

		queue := d.queue
		d.queue = d.queue[:0]

		for _, node := range queue {
			wg.Add(1)
			go func(node string) {
				defer wg.Done()
				d.methods[node]() // 并发执行节点对应的方法

				d.lock.Lock()
				for _, dependent := range d.nodes[node] {
					d.indegrees[dependent]--
					if d.indegrees[dependent] == 0 {
						d.queue = append(d.queue, dependent)
					}
				}
				d.lock.Unlock()
			}(node)
		}

		wg.Wait()
	}
}
