package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"
)

// 减少内存分配、文件读写等系统调用次数

var (
	zRecvCount = uint32(0)      // 张大爷听到了多少句话
	lRecvCount = uint32(0)      // 李大爷听到了多少句话
	total      = uint32(100000) // 总共需要遇见多少次
)

var (
	z0 = "吃了没，您吶?"
	z3 = "嗨！吃饱了溜溜弯儿。"
	z5 = "回头去给老太太请安！"
	l1 = "刚吃。"
	l2 = "您这，嘛去？"
	l4 = "有空家里坐坐啊。"
)

type RequestResponse struct {
	Serial  uint32 // 序号
	Payload string // 内容
}

// go run -race  rpc/duplexcommunication/v2/reduce_syscall_times.go
// 2.681874687s 加入读buffer 预分配需要分配一个合适的大小，否则在请求量小的时候可能不如直接小量分配
// 1.251762426s 加入写buffer，批量写
// 序列化RequestResponse，并发送
// 序列化后的结构如下：
//
//	长度  4字节
//	Serial 4字节
//	PayLoad 变长

func batchWriter(rchan chan *RequestResponse, conn *net.TCPConn) int32 {
	batcher := NewBatch()
	sendCount := int32(0)
	// 监听发送是否结束
	go func() {
		for atomic.LoadInt32(&sendCount) < int32(total*3) {
			// fmt.Println(sendCount)
			runtime.Gosched()
		}
		close(rchan)
	}()

	batcher.Deal(rchan, 30, 50*time.Millisecond, func(rs []*RequestResponse) error {
		if len(rs) == 0 {
			return nil
		}
		totalLen := 0
		for _, r := range rs {
			totalLen += len(r.Payload) + 8
		}
		sendBuf := make([]byte, totalLen)
		index := 0
		for _, r := range rs {
			l := uint32(len(r.Payload) + 4)
			binary.BigEndian.PutUint32(sendBuf[index:index+4], l)
			binary.BigEndian.PutUint32(sendBuf[index+4:index+8], r.Serial)
			copy(sendBuf[index+8:], []byte(r.Payload))
			index += len(r.Payload) + 8
		}
		_, _ = conn.Write(sendBuf) // conn 写入会系统调用并加锁，减少次数

		atomic.AddInt32(&sendCount, int32(len(rs))) // 计数

		// raw, _ := json.Marshal(rs)
		// fmt.Printf("发送了: %d, content: %v \n", len(rs), string(raw))
		return nil
	})

	return sendCount
}

// 接收数据，反序列化成RequestResponse
func readFrom(conn *net.TCPConn, recvBuf []byte, recvIndex int) ([]*RequestResponse, int, error) {
	rets := make([]*RequestResponse, 0)
	n, err := conn.Read(recvBuf[recvIndex:])
	if err != nil {
		return nil, 0, err
	}
	readNum := n + recvIndex

	index := 0
	for {
		if (readNum - index) < 8 { // 说明不够存放payload
			break
		}

		contentLen := binary.BigEndian.Uint32(recvBuf[index : index+4])
		serial := binary.BigEndian.Uint32(recvBuf[index+4 : index+8])
		fullLen := int(contentLen) + 4
		if (readNum - index) < fullLen { // 检查buf是否能放下剩下的内容
			if fullLen > len(recvBuf) { // 如果超过了buffer
				completeBuf := make([]byte, fullLen-len(recvBuf))
				_, err := io.ReadFull(conn, completeBuf)
				if err != nil {
					return nil, 0, err
				}
				recvBuf = append(recvBuf, completeBuf...)
				payload := recvBuf[8:fullLen]
				rets = append(rets, &RequestResponse{Payload: string(payload), Serial: serial})
				return rets, 0, nil
			}
			break
		}

		rets = append(rets, &RequestResponse{Payload: string(recvBuf[index+8 : index+fullLen]), Serial: serial})

		index = index + fullLen
	}

	if readNum-index > 0 {
		copy(recvBuf[:readNum-index], recvBuf[index:readNum])
	}

	// data, _ := json.Marshal(rets)
	// fmt.Printf("rets: %s,readNum:%d, index:%d, recindex: %d\n", data, readNum, index, readNum-index)
	return rets, readNum - index, nil
}

// 张大爷的耳朵
func zhangDaYeListen(conn *net.TCPConn, wg *sync.WaitGroup, reqChan chan *RequestResponse) {
	defer wg.Done()
	recvbuf := make([]byte, 1024)
	recvIndex := 0

	for zRecvCount < total*3 {
		rs, nextRecvIndex, err := readFrom(conn, recvbuf, recvIndex)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		recvIndex = nextRecvIndex

		for _, r := range rs {
			tmp := r
			if r.Payload == l2 { // 如果收到：您这，嘛去？
				// 起携程是因为这里李大爷和张大爷的沟通是个环路，不写没法读不读没法写；必须把这里写异步出去；
				go func() { reqChan <- &RequestResponse{tmp.Serial, z3} }() // 回复：嗨！吃饱了溜溜弯儿
			} else if r.Payload == l4 { // 如果收到：有空家里坐坐啊。
				go func() { reqChan <- &RequestResponse{tmp.Serial, z5} }() // 回复：回头去给老太太请安！
			} else if r.Payload == l1 { // 如果收到：刚吃。
				// 不用回复
			} else {
				fmt.Println("张大爷听不懂：" + r.Payload)
				break
			}
			zRecvCount++
		}
	}
}

// 张大爷的嘴
func zhangDaYeSay(conn *net.TCPConn, reqChan chan *RequestResponse) {
	nextSerial := uint32(0)

	for i := uint32(0); i < total; i++ {
		reqChan <- &RequestResponse{nextSerial, z0}
		nextSerial++
	}
}

// 李大爷的耳朵，实现是和张大爷类似的
func liDaYeListen(conn *net.TCPConn, wg *sync.WaitGroup, reqChan chan *RequestResponse) {
	defer wg.Done()
	recvbuf := make([]byte, 1024)
	recvIndex := 0

	for lRecvCount < total*3 {
		rs, nextRecvIndex, err := readFrom(conn, recvbuf, recvIndex)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		recvIndex = nextRecvIndex

		for _, r := range rs {
			temp := r
			// fmt.Println("李大爷听到：" + r.Payload)
			if r.Payload == z0 { // 如果收到：吃了没，您吶?
				go func() { reqChan <- &RequestResponse{temp.Serial, l1} }() // 回复：刚吃。
			} else if r.Payload == z3 {
				// do nothing
			} else if r.Payload == z5 {
				// do nothing
			} else {
				fmt.Println("李大爷听不懂：" + temp.Payload)
				break
			}
			lRecvCount++
		}
	}
}

// 李大爷的嘴
func liDaYeSay(conn *net.TCPConn, reqChan chan *RequestResponse) {
	nextSerial := uint32(0)

	for i := uint32(0); i < total; i++ {
		reqChan <- &RequestResponse{nextSerial, l2}
		nextSerial++
		reqChan <- &RequestResponse{nextSerial, l4}
		nextSerial++
	}
}

func startServer(wg *sync.WaitGroup) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	fmt.Println("张大爷在胡同口等着 ...")
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("碰见一个李大爷:" + conn.RemoteAddr().String())
		reqChan := make(chan *RequestResponse, 30)

		go func() {
			sd := batchWriter(reqChan, conn)
			fmt.Println("zhangDaYe 发送了:", int(sd))
		}()
		go zhangDaYeListen(conn, wg, reqChan)
		go zhangDaYeSay(conn, reqChan)
	}
}

func startClient(wg *sync.WaitGroup) *net.TCPConn {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	reqChan := make(chan *RequestResponse, 30)
	go liDaYeListen(conn, wg, reqChan)
	go liDaYeSay(conn, reqChan)
	sd := batchWriter(reqChan, conn)
	fmt.Println("liDaYe 发送了:", int(sd))
	return conn
}

func main() {
	cpufile, _ := os.Create("./cpu.pprof")
	memfile, _ := os.Create("./mem.pprof")
	mutexfile, _ := os.Create("./mutex.pprof")
	// 程序运行时开启统计
	_ = pprof.StartCPUProfile(cpufile)

	var wg sync.WaitGroup
	wg.Add(2)

	go startServer(&wg)
	time.Sleep(1 * time.Second)
	t1 := time.Now()
	conn := startClient(&wg)
	wg.Wait()
	elapsed := time.Since(t1)
	fmt.Println("耗时: ", elapsed)

	conn.Close()
	// 程序结束时关闭
	pprof.StopCPUProfile()
	_ = pprof.Lookup("allocs").WriteTo(memfile, 0)
	_ = pprof.Lookup("mutex").WriteTo(mutexfile, 0)
}

type Batch struct {
	sendBuf  chan []*RequestResponse
	dealData []*RequestResponse
}

func NewBatch() *Batch {
	sendBuf := make(chan []*RequestResponse, 10)
	var data []*RequestResponse
	return &Batch{
		sendBuf:  sendBuf,
		dealData: data,
	}
}

// Deal 流程处理
func (bf *Batch) Deal(data chan *RequestResponse, limit int, sequence time.Duration, dealFunc func(data []*RequestResponse) error) {
	go func() {
		timeTicker := time.NewTicker(sequence)
		for {
			select {
			case <-timeTicker.C:
				// fmt.Println("timeTicker.C ! len:", len(bf.dealData))
				bf.sendBuf <- bf.dealData
				bf.dealData = []*RequestResponse{}
				timeTicker.Reset(sequence)
			case v, ok := <-data:
				if !ok {
					if len(bf.dealData) > 0 {
						bf.sendBuf <- bf.dealData
						bf.dealData = []*RequestResponse{}
					}
					close(bf.sendBuf)
					return
				}
				bf.dealData = append(bf.dealData, v)
				if len(bf.dealData) >= limit {
					// fmt.Println("limit.reach ! len:", len(bf.dealData))
					bf.sendBuf <- bf.dealData
					bf.dealData = []*RequestResponse{}
					timeTicker.Stop()
					timeTicker.Reset(sequence)
				}
			}
		}
	}()

	for v := range bf.sendBuf {
		// fmt.Println("0---------- sendBuf ! len:", len(v))
		err := dealFunc(v)
		if err != nil {
			fmt.Printf("%+v", err)
		}
	}
}
