package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwadedForHeader         = "x-forwarded-for"
)

// Metadata holds the metadata we need
type Metadata struct {
	UserAgent string
	ClientIP  string
}

// extractMetadata parses the metadata we want from the request
func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// to parse user agent from gateway requests
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		// to parse user agent from gRPC requests
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		// to parse client ip from gateway requests
		if clientIPs := md.Get(xForwadedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	// here we added peer for directli gRPC requests to parse client ip
	if pr, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = pr.Addr.String()
	}

	return mtdt
}
