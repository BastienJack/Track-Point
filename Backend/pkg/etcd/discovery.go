package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"commerce/pkg/zap"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	defaultWeight = 10
)

// EtcdResolver is a resolver using etcd.
type EtcdResolver struct {
	etcdClient *clientv3.Client
}

// NewEtcdResolver creates an etcd based resolver with empty username and password.
func NewEtcdResolver(endpoints []string) (discovery.Resolver, error) {
	return NewEtcdResolverWithAuth(endpoints, "", "")
}

// NewEtcdResolverWithAuth creates an etcd based resolver with given username and password.
func NewEtcdResolverWithAuth(endpoints []string, username, password string) (discovery.Resolver, error) {
	// create an etcd client using username and password
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}

	// return an etcd resolver
	return &EtcdResolver{
		etcdClient: etcdClient,
	}, nil
}

// Diff implements the Resolver interface.
func (e *EtcdResolver) Diff(cacheKey string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(cacheKey, prev, next)
}

// Name implements the Resolver interface.
func (e *EtcdResolver) Name() string {
	return "etcd"
}

// Target implements the Resolver interface.
func (e *EtcdResolver) Target(ctx context.Context, target rpcinfo.EndpointInfo) (description string) {
	return target.ServiceName()
}

// Resolve implements the Resolver interface.
func (e *EtcdResolver) Resolve(ctx context.Context, serviceName string) (discovery.Result, error) {
	logger := zap.InitLogger()

	// userservice discovery from etcd using get
	prefix := ServiceKeyPrefix(serviceName)
	resp, err := e.etcdClient.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return discovery.Result{}, err
	}

	// userservice instances info
	var (
		info      InstanceInfo
		instances []discovery.Instance
	)

	for _, kv := range resp.Kvs {
		// get userservice instances info
		err := json.Unmarshal(kv.Value, &info)
		if err != nil {
			logger.Warnf("fail to unmarshal with err: %v, ignore key: %v", err, string(kv.Key))
			continue
		}

		// instance weight assign
		weight := info.Weight
		if weight <= 0 {
			weight = defaultWeight
		}

		// concatenate userservice instances
		instances = append(instances, discovery.NewInstance(info.Network, info.Address, weight, info.Tags))
	}

	// no remain userservice instances
	if len(instances) == 0 {
		return discovery.Result{}, fmt.Errorf("no instance remains for %v", serviceName)
	}

	// return discovered userservice instances
	return discovery.Result{
		Cacheable: true,
		CacheKey:  serviceName,
		Instances: instances,
	}, nil
}
