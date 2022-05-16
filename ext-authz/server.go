package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
)

var grpcport = flag.String("grpcport", ":50051", "grpcport")

var jwtKey = []byte("my_secret_key")

var expirationTime = time.Now().Add(10 * time.Minute)

type Claims struct {
	Deployment string `json:"deployment"`
	jwt.StandardClaims
}

type AuthorizationServer struct{}

func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Check() invoked from AuthorizationServer")

	_, err := json.Marshal(req.Attributes.Request.Http.Headers)
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
	authHeader, ok := req.Attributes.Request.Http.Headers["authorization"]
	if ok {
		var splitToken = strings.Split(authHeader, "Bearer ")
		fmt.Printf("Token : %v", splitToken[1])
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(splitToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil && !tkn.Valid {
			log.Printf("Error parsing token: %v", err)
			return &auth.CheckResponse{
				Status: &rpcstatus.Status{
					Code: int32(codes.Unauthenticated),
				},
				HttpResponse: &auth.CheckResponse_DeniedResponse{
					DeniedResponse: &auth.DeniedHttpResponse{
						Status: &envoy_type.HttpStatus{
							Code: envoy_type.StatusCode_Unauthorized,
						},
						Body: "Token expired or malformed",
					},
				},
			}, nil
		}
		log.Printf("Deployment: %v", claims.Deployment)
		return &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(codes.OK),
			},
			HttpResponse: &auth.CheckResponse_OkResponse{
				OkResponse: &auth.OkHttpResponse{
					Headers: []*envoy_config_core_v3.HeaderValueOption{
						{
							Header: &envoy_config_core_v3.HeaderValue{
								Key:   "x-wso2-cluster",
								Value: claims.Deployment,
							},
						},
					},
				},
			},
		}, nil
	}
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(codes.Unauthenticated),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: "No authorization header found",
			},
		},
	}, nil

}

func prodToken(w http.ResponseWriter, req *http.Request) {
	claims := &Claims{
		Deployment: "clusterProd",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	log.Printf("Production token: %v", tokenString)
	if err == nil {
		fmt.Fprint(w, tokenString+"\n")
	}
}

func sandToken(w http.ResponseWriter, req *http.Request) {
	claims := &Claims{
		Deployment: "clusterSand",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	log.Printf("Sandbox token: %v", tokenString)
	if err == nil {
		fmt.Fprint(w, tokenString+"\n")
	}
}

func main() {
	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		http.HandleFunc("/prod", prodToken)
		http.HandleFunc("/sand", sandToken)
		http.ListenAndServe(":8080", nil)
	}()

	s := grpc.NewServer()
	auth.RegisterAuthorizationServer(s, &AuthorizationServer{})

	log.Printf("Starting gRPC server at %s", *grpcport)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
