package dag

import (
	"fmt"
	"testing"
	"time"
)

func TestTopologicalSort(t *testing.T) {
	d := NewDAGTaskSimple()
	d.AddNode("A", func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("A finish %d\n", time.Now().UnixMilli())
	})
	d.AddNode("B", func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("B finish %d\n", time.Now().UnixMilli())
	})
	d.AddNode("C", func() {
		time.Sleep(150 * time.Millisecond)
		fmt.Printf("C finish %d\n", time.Now().UnixMilli())
	})
	d.AddNode("D", func() {
		time.Sleep(300 * time.Millisecond)
		fmt.Printf("D finish %d\n", time.Now().UnixMilli())
	})
	d.AddNode("E", func() {
		time.Sleep(300 * time.Millisecond)
		fmt.Printf("E finish %d\n", time.Now().UnixMilli())
	})

	d.AddEdge("A", "B")
	d.AddEdge("A", "C")
	d.AddEdge("B", "D")
	d.AddEdge("C", "D")
	d.AddEdge("D", "E")

	d.TopologicalSort()

	// Assert the order of execution of the methods
	// Add assertions for the expected behavior of your code
}
