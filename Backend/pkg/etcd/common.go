package etcd

// general etcd prefix
const etcdPrefix = "kitex/registry-etcd"

// ServiceKeyPrefix get userservice key prefix with etcd prefix.
func ServiceKeyPrefix(serviceName string) string {
	return etcdPrefix + "/" + serviceName
}

// ServiceKey generates the userservice key stored in etcd.
func ServiceKey(serviceName, addr string) string {
	return ServiceKeyPrefix(serviceName) + "/" + addr
}

// instanceInfo store userservice basic info in etcd.
type InstanceInfo struct {
	Network string            `json:"network"`
	Address string            `json:"address"`
	Weight  int               `json:"weight"`
	Tags    map[string]string `json:"tags"`
}
