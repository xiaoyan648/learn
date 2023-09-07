package pipe

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/reactivex/rxgo/v2"
)

type Customer struct {
	ID        int
	Name      string
	Age       int
	TaxNumber string
}

func PipeDemo() error {
	// Create the input channel
	ch := make(chan rxgo.Item)
	// Data producer
	go producer(ch)

	// Create an Observable
	ob := rxgo.FromChannel(ch)

	ob = ob.Filter(func(i interface{}) bool {
		// Filter operation
		customer := i.(Customer)
		return customer.Age > 18
	}).Map(
		func(_ context.Context, item interface{}) (interface{}, error) {
			// Enrich operation
			customer := item.(Customer)
			taxNumber, err := getTaxNumber()
			if err != nil {
				return nil, err
			}
			customer.TaxNumber = taxNumber
			return customer, nil
		},
		// Create multiple instances of the map operator
		rxgo.WithPool(10),
		// Serialize the items emitted by their Customer.ID
		rxgo.Serialize(func(item interface{}) int {
			customer := item.(Customer)
			return customer.ID
		}),
		rxgo.WithBufferedChannel(1),
	)

	for customer := range ob.Observe() {
		if customer.Error() {
			return customer.E
		}
		fmt.Println(customer)
	}
	return nil
}

func DeferDemo() {
	observable := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		next <- rxgo.Of(1)
		next <- rxgo.Of(2)
		next <- rxgo.Of(3)
	}})
	for customer := range observable.Observe() {
		fmt.Println(customer)
	}
	for customer := range observable.Observe() {
		fmt.Println(customer)
	}
}

func GroupDemo() {
	count := 3
	observable := rxgo.Range(0, 10).
		Filter(
			func(i interface{}) bool {
				// Filter operation
				v := i.(int)
				return v < 9
			},
			// rxgo.WithPool(4),
			// rxgo.WithBufferedChannel(4),
		).Map(
		func(_ context.Context, item interface{}) (interface{}, error) {
			// Enrich operation
			v := item.(int)
			v *= 2
			return v, nil
		},
		// Create multiple instances of the map operator
		// rxgo.WithPool(4),
		// rxgo.WithBufferedChannel(4),
	).GroupBy(
		count,
		func(item rxgo.Item) int {
			return item.V.(int) % count
		},
		rxgo.WithBufferedChannel(10),
		// rxgo.WithPool(2),
	)

	num := 0
	for i := range observable.Observe() {
		num++
		go func(it rxgo.Item, num int) {
			for n := range it.V.(rxgo.Observable).Observe() {
				fmt.Printf("item:%d %v\n", num, n.V)
			}
		}(i, num)
	}
	time.Sleep(1 * time.Second)
}

func IntervalDemo() {
	observable := rxgo.Interval(rxgo.WithDuration(5 * time.Second))
	for customer := range observable.Observe() {
		fmt.Println(customer)
	}
}

func getTaxNumber() (string, error) {
	time.Sleep(50 * time.Millisecond)
	return fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(1000)), nil
}

func producer(ch chan<- rxgo.Item) {
	// Batch Load a customer
	for i := 0; i < 10; i++ {
		customer := Customer{
			ID:   1 + i,
			Name: "John",
			Age:  15 + i,
		}
		// Emit the customer
		ch <- rxgo.Of(customer)
	}
	// Close the channel
	close(ch)
}
