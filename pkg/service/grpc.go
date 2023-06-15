package service

import (
	"net"

	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
	"google.golang.org/grpc"
)

type GRPCServiceCollectionItem struct {
	ServiceDesc     *grpc.ServiceDesc
	ServerInterface interface{}
}

type GRPCServiceCollection []*GRPCServiceCollectionItem

func (c GRPCServiceCollection) Len() int {
	return len(c)
}

func (c *GRPCServiceCollection) Add(svc *GRPCServiceCollectionItem) {
	*c = append(*c, svc)
}

type GRPCServer struct {
	listenAddr string
	services   GRPCServiceCollection
}

func NewGRPCServer(listenAddr string) *GRPCServer {
	return &GRPCServer{
		listenAddr: listenAddr,
	}
}

func (s *GRPCServer) RegisterService(svc *grpc.ServiceDesc, srv interface{}) {
	s.services.Add(&GRPCServiceCollectionItem{
		ServiceDesc:     svc,
		ServerInterface: srv,
	})
}

func (s *GRPCServer) registerServices(srv *grpc.Server) {
	for _, item := range s.services {
		srv.RegisterService(item.ServiceDesc, item.ServerInterface)
	}
}

func (s *GRPCServer) Run() (err error) {
	lis, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return errors.Critical.Wrap(err, "failed on network listener on")
	}

	srv := grpc.NewServer()
	s.registerServices(srv)

	return srv.Serve(lis)
}
