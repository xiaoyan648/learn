package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestFlow(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go startServer(&wg)
	time.Sleep(time.Second)
	conn := startClient(&wg)

	t1 := time.Now()
	wg.Wait()
	elapsed := time.Since(t1)

	conn.Close()
	fmt.Println("耗时: ", elapsed)

	// os.Exit(m.Run())
}

func TestBatch(t *testing.T) {
	batch := &Batch{
		sendBuf:  make(chan []*RequestResponse),
		dealData: []*RequestResponse{},
	}
	reqChan := make(chan *RequestResponse)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		batch.Deal(reqChan, 30, 100*time.Millisecond, func(data []*RequestResponse) error {
			for _, v := range data {
				fmt.Println("Deal:", v.Payload)
			}
			// raw, _ := json.Marshal(data)
			// n, err := conn.Write(raw)
			return nil
		})
	}()

	go func() {
		for i := 0; i < 200000; i++ {
			reqChan <- &RequestResponse{Serial: 1, Payload: "hello" + fmt.Sprint(i)}
		}
		close(reqChan)
	}()

	wg.Wait()
}
