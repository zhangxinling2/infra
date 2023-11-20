package lb

import (
	"math/rand"
	"sync/atomic"
)

var _ Balancer = new(RoundRobinBalancer)

// 简单轮训算法
type RoundRobinBalancer struct {
	ct uint32
}

func (r *RoundRobinBalancer) Next(key string, hosts []*ServerInstance) *ServerInstance {
	if len(hosts) == 0 {
		return nil
	}
	//自增
	count := atomic.AddUint32(&r.ct, 1)
	index := int(count) % len(hosts)
	instance := hosts[index]
	return instance
}

var _ Balancer = new(RandomBalancer)

// 随机负载均衡算法
type RandomBalancer struct {
}

func (r *RandomBalancer) Next(key string, hosts []*ServerInstance) *ServerInstance {
	if len(hosts) == 0 {
		return nil
	}
	//自增
	count := rand.Int31()
	index := int(count) % len(hosts)
	instance := hosts[index]
	return instance
}
