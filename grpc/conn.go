package grpc

import (
	"context"
	"github.com/coreos/etcd/clientv3/naming"
	"go.etcd.io/etcd/clientv3"
	grpcNaming "google.golang.org/grpc/naming"
	"time"
)

func NewGrpcResolover(service string, addr string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	r := &naming.GRPCResolver{Client: cli}
	err = r.Update(context.TODO(), service, grpcNaming.Update{Op: grpcNaming.Add, Addr: addr})
	if err != nil {
		return err
	}
}
