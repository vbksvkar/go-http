package propagation

import (
	"context"
	"net/http"

	gologger "github.com/vbksvkar/go-logger"
)

type propagationCtxKey struct{}

type PropagationValues struct {
	RequestId string
}

func WithPropagationValues(ctx context.Context, pv PropagationValues) context.Context {
	return context.WithValue(ctx, propagationCtxKey{}, pv)
}

func FromContext(ctx context.Context) PropagationValues {
	pv, ok := ctx.Value(propagationCtxKey{}).(PropagationValues)
	if !ok {
		logger := gologger.FromContext(ctx)
		logger.Warn("propagation values not found in context")
		return PropagationValues{}
	}
	return pv
}

func ExtractFromHeaders(h http.Header) PropagationValues {

	smallRequestId := h.Get("x-request-id")
	if smallRequestId != "" {
		return PropagationValues{
			RequestId: smallRequestId,
		}
	}

	caseRequestId := h.Get("X-Request-Id")
	if caseRequestId != "" {
		return PropagationValues{
			RequestId: caseRequestId,
		}
	}

	return PropagationValues{
		RequestId: "",
	}
}

func (pv PropagationValues) AddHeaders(h http.Header) {
	h.Set("x-request-id", pv.RequestId)
	h.Set("X-Request-Id", pv.RequestId)
}
