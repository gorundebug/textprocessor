package functions

import (
    "example.com/textprocessor/services/wordsprocessor/generated/pb"
    "github.com/gorundebug/servicelib/runtime"
    "strings"
)

type SplitTextFunc struct{}

func (f *SplitTextFunc) FlatMap(_ runtime.Stream, value *pb.TextData, out runtime.Collect[string]) {
    for _, w := range strings.Fields(value.Text) {
        out.Out(w)
    }
}

func MakeSplitTextFunc(runtime.StreamExecutionRuntime) *SplitTextFunc {
    return &SplitTextFunc{}
}
