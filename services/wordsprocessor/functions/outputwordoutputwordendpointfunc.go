package functions

import (
    "context"
    "fmt"
    "github.com/gorundebug/servicelib/runtime"
)

type OutputWordOutputWordEndpointFunc struct {
    values []string
}

func (ep *OutputWordOutputWordEndpointFunc) Start(ctx context.Context) error {
    ep.values = make([]string, 0)
    return nil
}

func (ep *OutputWordOutputWordEndpointFunc) Stop(ctx context.Context) {
}

func (ep *OutputWordOutputWordEndpointFunc) Consume(value string) {
    ep.values = append(ep.values, value)
    fmt.Printf("wp: %s\n", value)
}

func MakeOutputWordOutputWordEndpointFunc(runtime.StreamExecutionRuntime) *OutputWordOutputWordEndpointFunc {
    return &OutputWordOutputWordEndpointFunc{}
}
