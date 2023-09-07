第一版，总耗时 5.6s

第一版，cup占用主要开销
```shell
0      9.19s (flat, cum) 41.81% of Total
         .          .     45:func writeTo(r *RequestResponse, conn *net.TCPConn, lock *sync.Mutex) {
         .       40ms     46:   lock.Lock()
         .          .     47:   defer lock.Unlock()
         .          .     48:   payloadBytes := []byte(r.Payload)
         .          .     49:   serialBytes := make([]byte, 4)
         .          .     50:   binary.BigEndian.PutUint32(serialBytes, r.Serial)
         .          .     51:   length := uint32(len(payloadBytes) + len(serialBytes))
         .          .     52:   lengthByte := make([]byte, 4)
         .          .     53:   binary.BigEndian.PutUint32(lengthByte, length)
         .          .     54:
         .      3.08s     55:   conn.Write(lengthByte)
         .      3.14s     56:   conn.Write(serialBytes)
         .      2.93s     57:   conn.Write(payloadBytes)
         .          .     58:   // fmt.Println("发送: " + r.Payload)
         .          .     59:}
         .          .     60:
         .          .     61:// 接收数据，反序列化成RequestResponse

```
第一版，内存占用主要开销
```shell
         .          .     62:func readFrom(conn *net.TCPConn) (*RequestResponse, error) {
    8.50MB     8.50MB     63:   ret := &RequestResponse{}
       1MB        1MB     64:   buf := make([]byte, 4)
         .          .     65:   if _, err := io.ReadFull(conn, buf); err != nil {
         .          .     66:           return nil, fmt.Errorf("读长度故障：%s", err.Error())
         .          .     67:   }
         .          .     68:   length := binary.BigEndian.Uint32(buf)
         .          .     69:   if _, err := io.ReadFull(conn, buf); err != nil {
         .          .     70:           return nil, fmt.Errorf("读Serial故障：%s", err.Error())
         .          .     71:   }
         .          .     72:   ret.Serial = binary.BigEndian.Uint32(buf)
      14MB       14MB     73:   payloadBytes := make([]byte, length-4)
         .          .     74:   if _, err := io.ReadFull(conn, payloadBytes); err != nil {
         .          .     75:           return nil, fmt.Errorf("读Payload故障：%s", err.Error())
         .          .     76:   }
    9.50MB     9.50MB     77:   ret.Payload = string(payloadBytes)
         .          .     78:   return ret, nil
         .          .     79:}

```

第二版
主要瓶颈在系统调用上，需要消息发送和接收通过缓冲区减少系统调用。
- recvbuf，一次读出尽量多的内容处理，减少读文件的次数
- sendbuf，当发送的数量达到阈值或者达到一定时间在正真地进行网络调用


第三版
主要还可以对锁进行优化，减少锁的粒度，减少锁的冲突。