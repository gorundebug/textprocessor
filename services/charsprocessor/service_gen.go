package charsprocessor

import (
    "github.com/gorundebug/servicelib/runtime"
    runtimeserde "github.com/gorundebug/servicelib/runtime/serde"
    runtimeconfig "github.com/gorundebug/servicelib/runtime/config"
    "github.com/gorundebug/servicelib/transformation"
    "reflect"
    "context"
    "time"
    "sync"
    "example.com/textprocessor/services/charsprocessor/config"
    "google.golang.org/grpc"
    "fmt"
    log "github.com/sirupsen/logrus"
    "net"
    "example.com/textprocessor/services/charsprocessor/generated/grpcsvc"
	"example.com/textprocessor/services/charsprocessor/functions"
	"github.com/gorundebug/servicelib/datasink"
)

type ServiceConfig struct {
    runtimeconfig.ServiceAppConfig `mapstructure:",squash"`
    config                   config.Config `yaml:"config"`
}

type Service struct {
    runtime.ServiceApp
    serviceConfig *ServiceConfig
    //stream functions
    //data source functions
    //data sink functions
    outputCharOutputCharEndpointFunc *functions.OutputCharOutputCharEndpointFunc
    //streams
    split runtime.TypedBinarySplitStream[string]
    splitWord runtime.TypedTransformConsumedStream[string, rune]
    outputChar runtime.TypedSinkStream[rune]
    //data sources
    //data sinks
    outputCharOutputCharEndpointSink runtime.Consumer[rune]
    grpcServer *grpc.Server
}

func (s *Service) GetSerde(valueType reflect.Type) (runtimeserde.Serializer, error) {
    if serde, err := s.getCustomSerde(valueType); err != nil {
        return nil, err
    } else if serde != nil {
        return serde, nil
    }
    switch valueType {
    }
    return nil, nil
}

func (s *Service) SetConfig(config runtimeconfig.Config) {
    s.serviceConfig = config.(*ServiceConfig)
}

func (s *Service) StreamsInit(ctx context.Context) {
    s.split = transformation.SplitInStub[string]("Split", s)
    s.splitWord = transformation.FlatMapIterable[string, rune]("SplitWord", s.split.AddStream())
    s.outputChar = transformation.Sink[rune]("OutputChar", s.splitWord)
    s.outputCharOutputCharEndpointFunc = functions.MakeOutputCharOutputCharEndpointFunc(s)
    s.outputCharOutputCharEndpointSink = datasink.CustomEndpointSink[rune](s.outputChar, s.outputCharOutputCharEndpointFunc)
    s.grpcServer = grpc.NewServer()
}

func (s *Service) StartService(ctx context.Context) error {
    var err error
    s.StreamsInit(ctx)
    if err = s.start(ctx); err != nil {
        return err
    }

    var charsProcessorGrpcListener net.Listener
    charsProcessorGrpcListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.GetServiceConfig().GrpcHost, s.GetServiceConfig().GrpcPort))
    if err != nil {
        return fmt.Errorf("failed to listen gRPC port: %v", err)
    }
    grpcsvc.RegisterCharsProcessorServer(s.grpcServer, &GrpcCharsProcessor{
        service: s,
    })
    go func() {
        log.Infof("gRPC server listening at %v", charsProcessorGrpcListener.Addr())
        if err := s.grpcServer.Serve(charsProcessorGrpcListener); err != nil {
            log.Fatalf("failed to serve gRPC: %v", err)
        }
    }()

    return s.ServiceApp.Start(ctx)
}

func (s *Service) StopService(ctx context.Context) {
    timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.GetServiceConfig().ShutdownTimeout) * time.Millisecond)
    defer cancel()
    wg := sync.WaitGroup{}
    done := make(chan struct{})
    wg.Add(1)
    go func() {
        defer wg.Done()
        s.ServiceApp.Stop(timeoutCtx)
    }()

    grpcStopped := false
    wg.Add(1)
    go func() {
        defer wg.Done()
        s.grpcServer.GracefulStop()
        grpcStopped = true
    }()
    wg.Add(1)
    go func() {
        defer wg.Done()
        s.stop(timeoutCtx)
    }()
    go func() {
        wg.Wait()
        close(done)
    }()
    select {
    case <-done:
    case <-timeoutCtx.Done():

        if !grpcStopped {
            log.Warnf("gRPC server stop timed out: %s", timeoutCtx.Err())
        }

    }
}

type GrpcCharsProcessor struct {
    grpcsvc.UnimplementedCharsProcessorServer
    service *Service
}

func (s *GrpcCharsProcessor) SplitToSplitword(ctx context.Context, value *grpcsvc.SplitToSplitwordRequest) (*grpcsvc.SplitToSplitwordResponse, error) {
    s.service.split.Consume(string(value.Val))
    return &grpcsvc.SplitToSplitwordResponse{}, nil
}