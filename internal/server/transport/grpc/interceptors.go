package grpc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func CheckSubnetInterceptor(trustedSubnet string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "failed to get peer information")
		}

		if trustedSubnet != "" {
			agentIP := net.ParseIP(strings.TrimSpace(p.Addr.String()))
			_, subnet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				return nil, status.Error(codes.Internal, "failed to parse trusted subnet")
			}
			if !subnet.Contains(agentIP) {
				return nil, status.Error(codes.PermissionDenied, "client IP is not in trusted subnet")
			}
		}

		return handler(ctx, req)
	}
}

func HMACInterceptor(secretKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "failed to get metadata from context")
		}

		h := md.Get("HashSHA256")
		if len(h) == 0 || h[0] == "none" {
			return handler(ctx, req)
		}

		body, err := proto.Marshal(req.(proto.Message))
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to marshal request")
		}

		h1 := hmac.New(sha256.New, []byte(secretKey))
		_, err = h1.Write(body)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to compute HMAC")
		}

		res := hex.EncodeToString(h1.Sum(nil))
		if !hmac.Equal([]byte(h[0]), []byte(res)) {
			return nil, status.Error(codes.PermissionDenied, "HMAC verification failed")
		}

		return handler(ctx, req)
	}
}
