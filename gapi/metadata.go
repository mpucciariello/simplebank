package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	md := &Metadata{}

	if m, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := m.Get(grpcGatewayUserAgentHeader)[0]; len(ua) > 0 {
			md.UserAgent = m.Get(grpcGatewayUserAgentHeader)[0]
		}
		if ua := m.Get(userAgentHeader)[0]; len(ua) > 0 {
			md.UserAgent = m.Get(userAgentHeader)[0]
		}
		if cip := m.Get(xForwardedForHeader)[0]; len(cip) > 0 {
			md.UserAgent = m.Get(xForwardedForHeader)[0]
		}
	}

	if ip, ok := peer.FromContext(ctx); ok {
		md.ClientIP = ip.Addr.String()
	}

	return md
}
