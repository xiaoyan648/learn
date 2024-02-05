package dag

import (
	"fmt"
)

func ExampleBFS() {
	dag := NewDAGDemo()
	ret, iscrc := BFS(dag.Vertexs[0])
	keys := make([]string, 0, len(ret))
	for _, v := range ret {
		keys = append(keys, v.Key)
	}
	fmt.Printf("ret %v, iscrc: %v", keys, iscrc)

	// Output:
	// ret [i g f h e d c b a], iscrc: true
}

// 分层并发，没有达到最大并发度
func ExampleGorup() {
	dag1 := NewDAGDemo()
	ret1 := GourpTasks(dag1.Vertexs[0])
	keys := make([][]string, 0, len(ret1))
	for _, vs := range ret1 {
		tmp := make([]string, 0, len(vs))
		for _, v := range vs {
			tmp = append(tmp, v.Key)
		}
		keys = append(keys, tmp)
	}
	fmt.Printf("ret %v", keys)

	// Output:
	// ret [[i] [e h f g] [b c d] [a]]
}
