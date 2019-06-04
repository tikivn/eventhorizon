package oc

import (
	"context"
	"encoding/hex"

	eh "github.com/looplab/eventhorizon"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

type spanContextKeyType struct{}

var (
	spanContextKey = spanContextKeyType{}
)

const (
	spanContextKeyStr = "span_context"
)

func init() {
	// Register the SpanContext context.
	eh.RegisterContextMarshaler(MarshalSpanContext)
	eh.RegisterContextUnmarshaler(UnmarshalSpanContext)
}

func MarshalSpanContext(
	ctx context.Context, vals map[string]interface{},
) {
	if span := trace.FromContext(ctx); span != nil {
		spanContext := span.SpanContext()
		// TODO(giautm): Should use other format to easy reading?
		if bytes := propagation.Binary(spanContext); bytes != nil {
			vals["trace_id"] = hex.EncodeToString(spanContext.TraceID[:])
			vals[spanContextKeyStr] = hex.EncodeToString(bytes)
		}
	}
}

func UnmarshalSpanContext(
	ctx context.Context, vals map[string]interface{},
) context.Context {
	if val, exist := vals[spanContextKeyStr]; exist {
		if spanContextHex, ok := val.(string); ok {
			bytes, err := hex.DecodeString(spanContextHex)
			if err != nil {
				return ctx
			}

			if spanContext, ok := propagation.FromBinary(bytes); ok {
				return context.WithValue(ctx, spanContextKey, spanContext)
			}
		}
	}

	return ctx
}

func StartSpan(ctx context.Context, name string, opts ...trace.StartOption) (context.Context, *trace.Span) {
	val := ctx.Value(spanContextKey)
	if spanContext, ok := val.(trace.SpanContext); ok {
		return trace.StartSpanWithRemoteParent(ctx, name, spanContext, opts...)
	}
	return trace.StartSpan(ctx, name, opts...)
}
