package charsprocessor

import (
	"github.com/gorundebug/servicelib/runtime"
    runtimeserde "github.com/gorundebug/servicelib/runtime/serde"
	"reflect"
    "context"
)

func (s *Service) getCustomSerde(valueType reflect.Type) (runtimeserde.Serializer, error) {
	return nil, nil
}

func (s *Service) start(ctx context.Context) error {
	return nil
}

func (s *Service) stop(ctx context.Context) {
}

func (s *Service) GetEndpointReader(endpoint runtime.Endpoint, stream runtime.Stream,
    valueType reflect.Type) runtime.EndpointReader {
    return nil
}

func (s *Service) GetEndpointWriter(endpoint runtime.Endpoint, stream runtime.Stream,
    valueType reflect.Type) runtime.EndpointWriter {
    return nil
}