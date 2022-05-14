package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
)

var grpcport = flag.String("grpcport", ":50051", "grpcport")

type AuthorizationServer struct{}

func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Check() invoked from AuthorizationServer")

	headers, err := json.Marshal(req.Attributes.Request.Http.Headers)
	log.Println(string(headers))
	if err != nil {
		log.Fatalf("Error marshalling request headers: %v", err)
		return &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(codes.Unauthenticated),
			},
			HttpResponse: &auth.CheckResponse_DeniedResponse{
				DeniedResponse: &auth.DeniedHttpResponse{
					Status: &envoy_type.HttpStatus{
						Code: envoy_type.StatusCode_Unauthorized,
					},
					Body: "Authorization Header malformed or not provided",
				},
			},
		}, nil
	}
	authHeader := req.Attributes.Request.Http.Headers["authorization"]
	var splitToken = strings.Split(authHeader, "Bearer ")
	fmt.Printf("Token is : %v", splitToken[1])
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(codes.OK),
		},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{},
		},
		DynamicMetadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"test": {
					Kind: &structpb.Value_StringValue{
						StringValue: "123456",
					},
				},
			},
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	auth.RegisterAuthorizationServer(s, &AuthorizationServer{})

	log.Printf("Starting gRPC server at %s", *grpcport)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
