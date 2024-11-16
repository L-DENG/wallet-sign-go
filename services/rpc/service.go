package rpc

import (
	"context"
	"fmt"
	"github.com/L-DENG/wallet-sign-go/leveldb"
	"github.com/L-DENG/wallet-sign-go/protobuf/wallet"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync/atomic"
)

const MaxRecvMessageSize = 1024 * 1024 * 300

type RpcServerConfig struct {
	GrpcHostName string
	GrpcPort     int
}

type RpcServer struct {
	*RpcServerConfig
	db *leveldb.Keys

	*wallet.UnimplementedWalletServiceServer
	stopped atomic.Bool
}

func (r *RpcServer) Stop(ctx context.Context) error {
	r.stopped.Store(true)
	return nil
}

func (r *RpcServer) Stopped() bool {
	return r.stopped.Load()
}

func NewRpcServer(db *leveldb.Keys, config *RpcServerConfig) (*RpcServer, error) {
	return &RpcServer{
		db:              db,
		RpcServerConfig: config,
	}, nil
}

func (rs *RpcServer) Start(ctx context.Context) error {
	go func(s *RpcServer) {
		addr := fmt.Sprintf("%s:%d", rs.GrpcHostName, rs.GrpcPort)
		log.Info("Start rpc service ...")
		log.Info("addr :", addr)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("Could not start tcp listener. ")
		}

		opt := grpc.MaxRecvMsgSize(MaxRecvMessageSize)
		gs := grpc.NewServer(opt, grpc.ChainUnaryInterceptor(nil))

		reflection.Register(gs)

		wallet.RegisterWalletServiceServer(gs, s)

		log.Info("rpc service started: address:", listener.Addr().String(), "port:", s.GrpcPort)

		if err := gs.Serve(listener); err != nil {
			log.Error("rpc serve error:", err)
		}
	}(rs)
	panic("implement me")
}
