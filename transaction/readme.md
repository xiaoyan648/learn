## 本地消息表
- 实现原理
  1.利用本地事务，处理业务的同时格外插入一条状态记录，调用其他服务成功后状态更新；
  2.定时任务检查是否状态变更，如果没有变更则重试（下游接口注意幂等）；
- 特点
  无需外部依赖，直接好理解（开发量大）
- 适用场景
  对于大部分最终一致的场景都适用
  
## saga
- 实现原理
  执行+回滚的方案，执行失败会进行对应的回滚逻辑，执行和回滚操作都可重试；
  是一种最终一致的实现，执行过程中并不是完全一致的，可以查询到中间状态；
- 特点
  开发量少，也不会锁资源，异步执行;
- 适用场景
  适用于跨行转帐，下订单等场景;