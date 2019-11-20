package server

import (
	"context"
	"environment/dump"
	"environment/logger"
	"net"
	"singledb/proto/pbsingledb"

	"google.golang.org/grpc"
)

// Server struct
type Server struct {
	pbsingledb.UnimplementedSingleDBServerServer
}

// NewServer new
func NewServer() *Server {
	s := &Server{}
	return s
}

// Run server
func (s *Server) Run(addr string) error {
	logger.Info("start listen... addr:", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("failed to listen, err:", err)
		return err
	}

	srv := grpc.NewServer()
	pbsingledb.RegisterSingleDBServerServer(srv, s)

	if err := srv.Serve(lis); err != nil {
		logger.Error("failed to serve, err:", err)
	}
	return err
}

// Save implements proto.
func (s *Server) Save(ctx context.Context, req *pbsingledb.SingleMsgChunk) (*pbsingledb.SingleMsgChunkReply, error) {
	logger.Debug("save singlemsg chunk key:", string(req.Key), " data.len:", len(req.Data))

	// 网络事件处理计数器，dump会通过配置将当前服务的网络事件吞吐量提交给监控服务
	dump.NetEventRecvIncr(0)
	defer dump.NetEventRecvDecr(0)

	// 构建回包 & 处理业务
	return &pbsingledb.SingleMsgChunkReply{}, nil
}
