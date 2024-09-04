package functions

import (
    "github.com/gorundebug/servicelib/runtime"
    "strings"
)

type FilterWordFunc struct{}

func (f *FilterWordFunc) Filter(_ runtime.Stream, value string) bool {
    return !strings.HasPrefix(value, "A")
}

func MakeFilterWordFunc(runtime.StreamExecutionRuntime) *FilterWordFunc {
    return &FilterWordFunc{}
}
