package jaeger

import (
	"context"

	"google.golang.org/grpc"
)

type TracerGrpcInterceptor struct {
	tracer Tracer
}

func NewTracerGrpcInterceptor(tracer Tracer) *TracerGrpcInterceptor {
	return &TracerGrpcInterceptor{
		tracer: tracer,
	}
}

func (i *TracerGrpcInterceptor) UnaryClientInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, span := i.tracer.StartSpan(ctx, method)
		defer span.End()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
