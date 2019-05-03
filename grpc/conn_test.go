package grpc_test

import (
	grpcUtils "github.com/sundogrd/gopkg/grpc"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestNewGrpcResolover(t *testing.T) {
	r, err := grpcUtils.NewGrpcResolover()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resolver: %#v", r)
}

func TestResgiterServer(t *testing.T) {
	r, err := grpcUtils.NewGrpcResolover()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resolver: %#v", r)
	err = grpcUtils.ResgiterServer(*r, "worinima", "localhost:58911", 5*time.Second, 5)
	if err != nil {
		t.Fatalf("RegisterServer err: %s", err.Error())
		panic(err)
	}
	time.Sleep(5 * time.Minute)
}

func TestResgiterServer_client(t *testing.T) {
	r, err := grpcUtils.NewGrpcResolover()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resolver: %#v", r)
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial("sundog.comment", grpc.WithBalancer(b), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conn.GetState())
}