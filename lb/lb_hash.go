package lb

import (
	"hash/crc32"
)

var _ Balancer = new(HashBalancer)

// 哈希负载均衡算法
type HashBalancer struct {
}

func (r *HashBalancer) Next(key string, hosts []*ServerInstance) *ServerInstance {
	if len(hosts) == 0 {
		return nil
	}
	//自增
	count := crc32.ChecksumIEEE([]byte(key))
	index := int(count) % len(hosts)
	instance := hosts[index]
	return instance
}
