package grpc_test

import (
	"github.com/sundogrd/gopkg/grpc"
	"testing"
)

func TestNewGrpcResolover(t *testing.T) {
	r, err := grpc.NewGrpcResolover()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resolver: %#v", r)
}