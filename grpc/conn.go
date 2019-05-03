package grpc

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/naming"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/sirupsen/logrus"
	grpcNaming "google.golang.org/grpc/naming"
	"log"
	"time"
)

var stopSignal chan bool

// TODO: gRPC服务发现
// http://ralphbupt.github.io/2017/11/27/etcd%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0/
func NewGrpcResolover() (*naming.GRPCResolver, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	r := &naming.GRPCResolver{Client: cli}
	if err != nil {
		return nil, err
	}
	return r, nil
}

func ResgiterServer(r naming.GRPCResolver, service string, addr string, interval time.Duration, ttl int) error {
	go func() {
		client := r.Client
		serviceValue := string(service + ":" + addr)
		serviceKey := service + ":instance"
		ticker := time.NewTicker(interval)
		for {
			// minimum lease TTL is ttl-second
			resp, _ := client.Grant(context.TODO(), int64(ttl))
			// should get first, if not exist, set it
			_, err := client.Get(context.Background(), serviceKey)
			if err != nil {
				if err == rpctypes.ErrKeyNotFound {
					if _, err := client.Put(context.TODO(), serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
						log.Printf("grpclb: set service '%s' with ttl to etcd3 failed: %s", service, err.Error())
					}
				} else {
					log.Printf("grpclb: service '%s' connect to etcd3 failed: %s", service, err.Error())
				}
			} else {
				// refresh set to true for not notifying the watcher
				if _, err := client.Put(context.Background(), serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
					logrus.Infof("grpclb: refresh service '%s' with ttl to etcd3 failed: %s", service, err.Error())

				}
			}
			select {
			case <-stopSignal:
				return
			case <-ticker.C:
			}
		}
	}()
	err := r.Update(context.TODO(), service, grpcNaming.Update{Op: grpcNaming.Add, Addr: addr})

	return err
}

func UnRegisterServer(r naming.GRPCResolver, service string) error {
	stopSignal <- true
	stopSignal = make(chan bool, 1) // just a hack to avoid multi UnRegister deadlock
	serviceKey := service + ":instance"
	var err error
	if _, err := r.Client.Delete(context.Background(), serviceKey); err != nil {
		log.Printf("grpclb: deregister '%s' failed: %s", serviceKey, err.Error())
	} else {
		log.Printf("grpclb: deregister '%s' ok.", serviceKey)
	}
	return err
}
