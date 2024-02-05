package dag

// 有向无环图：邻接矩阵、邻接表
// 图结构
type DAG struct {
	Vertexs []*Vertex
}

// 顶点
type Vertex struct {
	Key   string
	Value interface{}
	// Parents  []*Vertex
	Children []*Vertex
}

// 添加顶点
func (dag *DAG) AddVertex(v *Vertex) {
	dag.Vertexs = append(dag.Vertexs, v)
}

// 添加边
func (dag *DAG) AddEdge(from, to *Vertex) {
	from.Children = append(from.Children, to)
	// to.Parents = append(to.Parents, from)
}

func NewDAGDemo() *DAG {
	dag := &DAG{}
	va := &Vertex{Key: "a", Value: "1"}
	vb := &Vertex{Key: "b", Value: "2"}
	vc := &Vertex{Key: "c", Value: "3"}
	vd := &Vertex{Key: "d", Value: "4"}
	ve := &Vertex{Key: "e", Value: "5"}
	vf := &Vertex{Key: "f", Value: "6"}
	vg := &Vertex{Key: "g", Value: "7"}
	vh := &Vertex{Key: "h", Value: "8"}
	vi := &Vertex{Key: "i", Value: "9"}
	// 添加顶点
	dag.AddVertex(va)
	dag.AddVertex(vb)
	dag.AddVertex(vc)
	dag.AddVertex(vd)
	dag.AddVertex(ve)
	dag.AddVertex(vf)
	dag.AddVertex(vg)
	dag.AddVertex(vh)
	dag.AddVertex(vi)

	// 添加边
	dag.AddEdge(va, vb)
	dag.AddEdge(va, vc)
	dag.AddEdge(va, vd)
	dag.AddEdge(vb, ve)
	dag.AddEdge(vb, vh)
	dag.AddEdge(vb, vf)
	dag.AddEdge(vc, vf)
	dag.AddEdge(vc, vg)
	dag.AddEdge(vd, vg)
	dag.AddEdge(vh, vi)
	dag.AddEdge(ve, vi)
	dag.AddEdge(vf, vi)
	// dag.AddEdge(vg, vi)
	return dag
}

// BFS 进行分层.
func BFS(root *Vertex) ([]*Vertex, bool) {
	// 定义结果数组
	var (
		result []*Vertex
		queue  []*Vertex
		isCrc  bool // 是否有环路
	)
	// 初始化访问标记
	visited := make(map[string]struct{})
	// 将根节点添加到队列中
	queue = append(queue, root)
	// 循环遍历队列
	for len(queue) > 0 {
		// 出队
		v := queue[0]
		queue = queue[1:]
		// 如果已访问过，则跳过
		if _, ok := visited[v.Key]; ok {
			isCrc = true
			continue
		}
		// 更新访问标记
		visited[v.Key] = struct{}{}
		// 将当前节点添加到结果数组中
		result = append([]*Vertex{v}, result...)
		// 入队
		queue = append(queue, v.Children...)
	}
	// 返回结果数组和是否有环路
	return result, isCrc
}

func GourpTasks(root *Vertex) [][]*Vertex {
	result := make([][]*Vertex, 0)
	queue := make([]*Vertex, 0)
	visited := make(map[string]struct{})

	queue = append(queue, root)
	for len(queue) > 0 {
		qSize := len(queue)
		tmp := make([]*Vertex, 0, qSize)
		for i := 0; i < qSize; i++ {
			// 出队
			v := queue[0]
			queue = queue[1:]
			if _, ok := visited[v.Key]; ok {
				continue
			}
			visited[v.Key] = struct{}{}
			tmp = append(tmp, v)
			// 入队
			queue = append(queue, v.Children...)
		}
		result = append([][]*Vertex{tmp}, result...)
	}
	return result
}
