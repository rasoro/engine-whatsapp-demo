package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/servers/grpc/grpc_servers"
	"github.com/weni/whatsapp-router/servers/grpc/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config     config.Config
	Db         *mongo.Database
	grpcServer *grpc.Server
}

func NewServer(db *mongo.Database) *Server {
	conf := config.GetConfig()
	return &Server{
		Db:     db,
		config: *conf,
	}
}

func (s *Server) Start() error {
	channelServer := grpc_servers.NewChannelCServer(s.Db)
	s.grpcServer = grpc.NewServer()
	pb.RegisterChannelServiceServer(s.grpcServer, channelServer)
	reflection.Register(s.grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", s.config.Server.GRPCPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("Start grpc server :%v", s.config.Server.GRPCPort)

	// s.WaitGroup.Add(1)
	go func() {
		// defer s.WaitGroup.Done()
		err = s.grpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
