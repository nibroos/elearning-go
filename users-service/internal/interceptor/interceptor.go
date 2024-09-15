package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract metadata from the incoming context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Log the metadata for debugging purposes
		log.Printf("Received metadata: %v", md)

		// Add any custom logic here, for example, setting values in the context
		type contextKey string

		const filtersKey contextKey = "filters"

		newCtx := context.WithValue(ctx, filtersKey, md)

		// Call the handler with the new context
		return handler(newCtx, req)
	}
}
