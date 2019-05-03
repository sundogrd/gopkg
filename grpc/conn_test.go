package grpc_test

import (
	"github.com/sundogrd/gopkg/grpc"
	"testing"
	"time"
)

func TestNewGrpcResolover(t *testing.T) {
	r, err := grpc.NewGrpcResolover()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resolver: %#v", r)
}

func TestResgiterServer(t *testing.T) {
	r, err := grpc.NewGrpcResolover()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resolver: %#v", r)
	err = grpc.ResgiterServer(*r, "worinima", "localhost:58911", 5*time.Second, 5)
	if err != nil {
		t.Fatalf("RegisterServer err: %s", err.Error())
		panic(err)
	}
	time.Sleep(5 * time.Minute)
}