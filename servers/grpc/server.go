package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/repositories"
	"github.com/weni/whatsapp-router/servers/grpc/pb"
	"github.com/weni/whatsapp-router/services"
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
	chanelRepository := repositories.NewChannelRepositoryDb(s.Db)
	channelService := services.NewChannelService(chanelRepository)
	s.grpcServer = grpc.NewServer()
	pb.RegisterChannelServiceServer(s.grpcServer, channelService)
	reflection.Register(s.grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", s.config.Server.GRPCPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error(err.Error())
		log.Fatal()
		return err
	}

	logger.Info(fmt.Sprintf("Start grpc server :%v", s.config.Server.GRPCPort))

	go func() {
		err = s.grpcServer.Serve(listener)
		if err != nil {
			logger.Error(err.Error())
			log.Fatal()
		}
	}()

	return nil
}
