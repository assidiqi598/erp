package internal

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ForwardMetadataInterceptor intercepts and propagates tokens
func ForwardMetadataInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract incoming metadata
		incomingMD, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			incomingMD = metadata.New(nil) // Initialize empty metadata if none exists
		}

		// Create a new outgoing context with the incoming metadata
		outgoingCtx := metadata.NewOutgoingContext(ctx, incomingMD)

		// Call the handler with the modified context
		return handler(outgoingCtx, req)
	}
}
