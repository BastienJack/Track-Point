package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"commerce/pkg/zap"
	"github.com/cloudwego/kitex/pkg/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ttlKey     = "KITEX_ETCD_REGISTRY_LEASE_TTL"
	defaultTTL = 60
)

func getTTL() int64 {
	var ttl int64 = defaultTTL
	if str, ok := os.LookupEnv(ttlKey); ok {
		if t, err := strconv.Atoi(str); err == nil {
			ttl = int64(t)
		}
	}
	return ttl
}

type EtcdRegistry struct {
	etcdClient *clientv3.Client
	leaseTTL   int64
	meta       *RegisterMeta
}

type RegisterMeta struct {
	ctx     context.Context
	leaseID clientv3.LeaseID
	cancel  context.CancelFunc
}

// NewEtcdRegistry creates a etcd based registry with empty username and password.
func NewEtcdRegistry(endpoints []string) (registry.Registry, error) {
	return NewEtcdRegistryWithAuth(endpoints, "", "")
}

// NewEtcdRegistryWithAuth creates a etcd based registry with given username and password.
func NewEtcdRegistryWithAuth(endpoints []string, username, password string) (registry.Registry, error) {
	// create an etcd client using username and password
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}

	// return a etcd registry object
	return &EtcdRegistry{
		etcdClient: etcdClient,
		leaseTTL:   getTTL(),
	}, nil
}

// grantLease request a lease from etcd userservice.
func (e *EtcdRegistry) grantLease() (clientv3.LeaseID, error) {
	// cancel if time out
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// request a lease from etcd userservice
	resp, err := e.etcdClient.Grant(ctx, e.leaseTTL)
	if err != nil {
		return clientv3.NoLease, err
	}

	return resp.ID, nil
}

// keepalive create a keepalive lease.
func (e *EtcdRegistry) keepAlive(meta *RegisterMeta) error {
	logger := zap.InitLogger()

	// request a keepalive lease
	aliveChannel, err := e.etcdClient.KeepAlive(meta.ctx, meta.leaseID)
	if err != nil {
		return err
	}

	//
	go func() {
		logger.Infof("Start keepalive lease %x for etcd registry.", meta.leaseID)

		// keepAlive is a channel receiving data continuously to keep lease alive.
		for range aliveChannel {
			select {
			case <-meta.ctx.Done():
				break
			default:
			}
		}

		logger.Infof("Stop keepalive lease %x for etcd registry.", meta.leaseID)
	}()

	return nil
}

// validateRegistryInfo validate a registry info.
func validateRegistryInfo(info *registry.Info) error {
	// check userservice name
	if info.ServiceName == "" {
		return fmt.Errorf("missing userservice name in Register")
	}

	// check ip address
	if info.Addr == nil {
		return fmt.Errorf("missing addr in Register")
	}

	return nil
}

func (e *EtcdRegistry) register(info *registry.Info, leaseID clientv3.LeaseID) error {
	// create json format registry info
	val, err := json.Marshal(&InstanceInfo{
		Network: info.Addr.Network(),
		Address: info.Addr.String(),
		Weight:  info.Weight,
		Tags:    info.Tags,
	})
	if err != nil {
		return err
	}

	// cancel if time out
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// client request etcd userservice to put userservice key to etcd
	_, err = e.etcdClient.Put(ctx, ServiceKey(info.ServiceName, info.Addr.String()),
		string(val), clientv3.WithLease(leaseID))

	return err
}

// Register registers a server with given registry info.
func (e *EtcdRegistry) Register(info *registry.Info) error {
	// validate registry info
	if err := validateRegistryInfo(info); err != nil {
		return err
	}

	// request a lease id in etcd
	leaseID, err := e.grantLease()
	if err != nil {
		return err
	}

	// register userservice using lease
	if err := e.register(info, leaseID); err != nil {
		return err
	}

	// create registry meta data
	registerMeta := RegisterMeta{
		leaseID: leaseID,
	}
	registerMeta.ctx, registerMeta.cancel = context.WithCancel(context.Background())

	// start a keep alive userservice
	if err := e.keepAlive(&registerMeta); err != nil {
		return err
	}
	e.meta = &registerMeta

	return nil
}

func (e *EtcdRegistry) deregister(info *registry.Info) error {
	// cancel if time out
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// delete userservice
	_, err := e.etcdClient.Delete(ctx, ServiceKey(info.ServiceName, info.Addr.String()))

	return err
}

// Deregister deregister a server with given registry info.
func (e *EtcdRegistry) Deregister(info *registry.Info) error {
	// validate registry info
	if err := validateRegistryInfo(info); err != nil {
		return err
	}

	// deregister
	if err := e.deregister(info); err != nil {
		return err
	}
	e.meta.cancel()

	return nil
}
