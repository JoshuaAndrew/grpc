package service

import (
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"github.com/sirupsen/logrus"
	"github.com/JoshuaAndrew/grpc/uuid"
	"github.com/JoshuaAndrew/grpc/api"
)

var (
	emptyResponse = &google_protobuf.Empty{}
)

type GreetingServiceImpl struct {
	id   string;
	name string;
	age  int64;
}

func NewGreetingService() (*GreetingServiceImpl, error) {
	s := &GreetingServiceImpl{}
	return s, nil
}

func (g *GreetingServiceImpl) Say(context context.Context, request *api.Request) (*google_protobuf.Empty, error) {
	return emptyResponse, nil
}

func (g *GreetingServiceImpl) SayHello(context context.Context, request *api.Request) (*api.Response, error) {
	r := &api.Response{
		Message:uuid.Rand().Hex(),
	}
	logrus.Info("Request Message:", request.GetId(), request.GetName(), request.GetAge())
	logrus.Info("Response Message:", r.GetMessage())
	return r, nil
}