package functions

import (
    "context"
    "fmt"
    "github.com/gorundebug/servicelib/runtime"
)

type OutputCharOutputCharEndpointFunc struct {
    values []rune
}

func (ep *OutputCharOutputCharEndpointFunc) Start(ctx context.Context) error {
    ep.values = make([]rune, 0)
    return nil
}

func (ep *OutputCharOutputCharEndpointFunc) Stop(ctx context.Context) {
}

func (ep *OutputCharOutputCharEndpointFunc) Consume(value rune) {
    ep.values = append(ep.values, value)
    fmt.Printf("cp: %c\n", value)
}

func MakeOutputCharOutputCharEndpointFunc(runtime.StreamExecutionRuntime) *OutputCharOutputCharEndpointFunc {
    return &OutputCharOutputCharEndpointFunc{}
}
