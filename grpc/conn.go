package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/naming"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Addr     string `json:"Addr"`
	Metadate string `json:"Metadate"`
}

var stopSignal chan bool

// TODO: gRPC服务发现
// http://ralphbupt.github.io/2017/11/27/etcd%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0/
func NewGrpcResolover(endpoints ...string) (*naming.GRPCResolver, error) {
	var config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if len(endpoints) != 0 {
		config.Endpoints = endpoints
	}

	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	r := &naming.GRPCResolver{Client: cli}
	if err != nil {
		return nil, err
	}
	return r, nil
}

func ResgiterServer(r naming.GRPCResolver, serviceName string, addr string, interval time.Duration, ttl int) error {
	service := Service{
		Addr:     addr,
		Metadate: "...",
	}
	bts, err := json.Marshal(service)
	if err != nil {
		return err
	}
	serviceValue := string(bts)
	serviceKey := fmt.Sprintf("%s/%s", serviceName, serviceValue)

	go func() {
		client := r.Client
		ticker := time.NewTicker(interval)
		for {
			// minimum lease TTL is ttl-second
			resp, err := client.Grant(context.TODO(), int64(ttl))
			// should get first, if not exist, set it
			getResp, err := client.Get(context.Background(), serviceKey)
			logrus.Print(getResp)
			if err != nil {
				if err == rpctypes.ErrKeyNotFound {
					if _, err := client.Put(context.TODO(), serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
						log.Printf("grpclb: set service '%s' with ttl to etcd3 failed: %s", service, err.Error())
					}
				} else {
					log.Printf("grpclb: service '%s' connect to etcd3 failed: %s", service, err.Error())
				}
			} else {
				if getResp.Count == 0 {
					if _, err := client.Put(context.TODO(), serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
						logrus.Printf("grpclb: set service '%s' with ttl to etcd3 failed: %s", service, err.Error())
					}
				}
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

	return nil
}

//func UnRegisterServer(r naming.GRPCResolver, serviceName string) error {
//	stopSignal <- true
//	stopSignal = make(chan bool, 1) // just a hack to avoid multi UnRegister deadlock
//	serviceKey := fmt.Sprintf("%s/%s", serviceName, serviceValue)
//	var err error
//	if _, err := r.Client.Delete(context.Background(), serviceKey); err != nil {
//		log.Printf("grpclb: deregister '%s' failed: %s", serviceKey, err.Error())
//	} else {
//		log.Printf("grpclb: deregister '%s' ok.", serviceKey)
//	}
//	return err
//}
