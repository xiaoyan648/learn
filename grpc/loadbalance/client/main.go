package main

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"math/rand"
	"time"

	pb "github.com/xiaoyan648/learn/grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const (
	myScheme      = "exmaple"
	myServiceName = "resolver.my.com"
)

type weightAddr struct {
	addr string
	w    int
}

var backendAddr = []weightAddr{
	{addr: "localhost:50051", w: 2},
	{addr: "localhost:50052", w: 1},
}

func callUnaryEcho(c pb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := pb.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/name_resolving")
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	exampleConn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:///%s", myScheme, myServiceName),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // This sets the initial balancing policy.
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer exampleConn.Close()

	fmt.Printf("--- calling helloworld.Greeter/SayHello to \"%s:///%s\"\n", myScheme, myServiceName)
	makeRPCs(exampleConn, 2)
}

// Following is an example name resolver. It includes a
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder)
// and a Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
//
// A ResolverBuilder is registered for a scheme (in this example, "example" is
// the scheme). When a ClientConn is created for this scheme, the
// ResolverBuilder will be picked to build a Resolver. Note that a new Resolver
// is built for each ClientConn. The Resolver will watch the updates for the
// target, and send updates to the ClientConn.

// exampleResolverBuilder is a
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder).
type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]weightAddr{
			myServiceName: backendAddr,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return myScheme }

// exampleResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]weightAddr
}

func (r *exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {
	// 直接从map中取出对于的addrList
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{
			Addr:               s.addr,
			BalancerAttributes: attributes.New("weight", s.w),
		}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) Close() {}

const WEIFGT_LB_NAME = "weight_lb"

// pickerBuilder impl grpc base.PickerBuilder
type pickerBuilder struct{}

func (*pickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for subConn, addr := range info.ReadySCs {
		weight := addr.Address.BalancerAttributes.Value("weight").(int)
		if weight <= 0 {
			weight = 1
		}
		for i := 0; i < weight; i++ {
			scs = append(scs, subConn)
		}
	}
	return &weightPicker{
		conns: scs,
	}
}

// picker impl grpc balancer.Picker
type weightPicker struct {
	conns []balancer.SubConn
}

func (p *weightPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if v := info.Ctx.Value("shard_key"); v != nil {
		key := v.(string)
		crc32 := crc32.ChecksumIEEE([]byte(key))
		index := crc32 % uint32(len(p.conns))
		return balancer.PickResult{SubConn: p.conns[index]}, nil
	}
	return balancer.PickResult{SubConn: p.conns[rand.Intn(len(p.conns))]}, nil
}

func init() {
	// Register the example ResolverBuilder. This is usually done in a package's
	// init() function.
	resolver.Register(&exampleResolverBuilder{})
	balancer.Register(base.NewBalancerBuilder(WEIFGT_LB_NAME, &pickerBuilder{}, base.Config{HealthCheck: false}))
}
