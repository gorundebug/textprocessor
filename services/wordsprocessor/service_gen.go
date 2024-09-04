package wordsprocessor

import (
    "github.com/gorundebug/servicelib/runtime"
    runtimeserde "github.com/gorundebug/servicelib/runtime/serde"
    runtimeconfig "github.com/gorundebug/servicelib/runtime/config"
    "github.com/gorundebug/servicelib/transformation"
    "reflect"
    "context"
    "time"
    "sync"
    "example.com/textprocessor/services/wordsprocessor/config"
    "google.golang.org/grpc/credentials/insecure"
	charsprocessor_grpcsvc "example.com/textprocessor/services/charsprocessor/generated/grpcsvc"
    "google.golang.org/grpc"
    "fmt"
    log "github.com/sirupsen/logrus"
	"example.com/textprocessor/services/wordsprocessor/serdes"
	"example.com/textprocessor/services/wordsprocessor/generated/pb"
	"example.com/textprocessor/services/wordsprocessor/functions"
	"github.com/gorundebug/servicelib/datasource"
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
    splitTextFuncFunc *functions.SplitTextFunc
    filterWordFuncFunc *functions.FilterWordFunc
    //data source functions
    //data sink functions
    outputWordOutputWordEndpointFunc *functions.OutputWordOutputWordEndpointFunc
    //streams
    inputText runtime.TypedInputStream[*pb.TextData]
    splitText runtime.TypedTransformConsumedStream[*pb.TextData, string]
    split runtime.TypedSplitStream[string]
    filterWord runtime.TypedConsumedStream[string]
    outputWord runtime.TypedSinkStream[string]
    fromSplitToSplitWord runtime.TypedStreamConsumer[string]
    fromSplitToSplitWordTimeout time.Duration
    //data sources
    inputTextInputTextEndpointSource runtime.Consumer[*pb.TextData]
    //data sinks
    outputWordOutputWordEndpointSink runtime.Consumer[string]
    charsProcessorGrpcClientConn *grpc.ClientConn
    charsProcessorGrpcClient charsprocessor_grpcsvc.CharsProcessorClient

}

func (s *Service) GetSerde(valueType reflect.Type) (runtimeserde.Serializer, error) {
    if serde, err := s.getCustomSerde(valueType); err != nil {
        return nil, err
    } else if serde != nil {
        return serde, nil
    }
    switch valueType {
    case runtimeserde.GetSerdeType[pb.TextData]():
    {
        var serde runtimeserde.Serde[*pb.TextData] = &serdes.TextDataSerde{}
        return serde, nil
    }

    }
    return nil, nil
}

func (s *Service) SetConfig(config runtimeconfig.Config) {
    s.serviceConfig = config.(*ServiceConfig)
}

func (s *Service) StreamsInit(ctx context.Context) {
    s.inputText = transformation.Input[*pb.TextData]("InputText", s)
    s.splitTextFuncFunc = functions.MakeSplitTextFunc(s)
    s.splitText = transformation.FlatMap[*pb.TextData, string]("SplitText", s.inputText, s.splitTextFuncFunc)
    s.split = transformation.Split[string]("Split", s.splitText)
    s.filterWordFuncFunc = functions.MakeFilterWordFunc(s)
    s.filterWord = transformation.Filter[string]("FilterWord", s.split.AddStream(), s.filterWordFuncFunc)
    s.outputWord = transformation.Sink[string]("OutputWord", s.filterWord)
    s.fromSplitToSplitWordTimeout = s.GetConsumeTimeout(3, 6)
    s.fromSplitToSplitWord = transformation.OutStub[string]("SplitWord", s.split.AddStream(),
        func(value string) error {
            var callCtx context.Context
            if s.fromSplitToSplitWordTimeout > 0 {
                ctxTimeout, cancel := context.WithTimeout(ctx, s.fromSplitToSplitWordTimeout)
                defer cancel()
                callCtx = ctxTimeout
            } else {
                callCtx = ctx
            }
            message := &charsprocessor_grpcsvc.SplitToSplitwordRequest{ Val:value}
            _, err := s.charsProcessorGrpcClient.SplitToSplitword(callCtx, message)
            return err
        })
    s.inputTextInputTextEndpointSource = datasource.NetHTTPEndpointConsumer[*pb.TextData](s.inputText)
    s.outputWordOutputWordEndpointFunc = functions.MakeOutputWordOutputWordEndpointFunc(s)
    s.outputWordOutputWordEndpointSink = datasink.CustomEndpointSink[string](s.outputWord, s.outputWordOutputWordEndpointFunc)
}

func (s *Service) StartService(ctx context.Context) error {
    var err error
    s.StreamsInit(ctx)
    if err = s.start(ctx); err != nil {
        return err
    }

    charsProcessorCfg := s.serviceConfig.GetServiceConfigByName("CharsProcessor")
    if charsProcessorCfg == nil {
        return fmt.Errorf("config for service 'CharsProcessor' was not found")
    }
    s.charsProcessorGrpcClientConn, err = grpc.NewClient(fmt.Sprintf("%s:%d", charsProcessorCfg.GrpcHost, charsProcessorCfg.GrpcPort),
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return fmt.Errorf("gRPC client did not connect to server 'CharsProcessor': %v", err)
    }
    s.charsProcessorGrpcClient = charsprocessor_grpcsvc.NewCharsProcessorClient(s.charsProcessorGrpcClientConn)

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

    charsProcessorGrpcClientStopped := false
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := s.charsProcessorGrpcClientConn.Close(); err != nil {
            log.Warnln(err)
        }
        charsProcessorGrpcClientStopped = true
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

        if !charsProcessorGrpcClientStopped {
            log.Warnf("gRPC client to server 'CharsProcessor' stop timed out: %s", timeoutCtx.Err())
        }

    }
}