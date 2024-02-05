package dag

// 基于 dag_task_simple, 做如下优化
// 1. 预生成执行节点，检查并打印执行顺序
// 2. dag对象封装
// 3. 性能提升

type Dagger struct{}
