package main

import (
    log "github.com/sirupsen/logrus"
    "github.com/gorundebug/servicelib/runtime"
    runtimeconfig "github.com/gorundebug/servicelib/runtime/config"
    "os"
    "os/signal"
    "syscall"
    "context"
    "example.com/textprocessor/services/charsprocessor"
    )

const (
    build_version = "v0.0.1"
    build_commit = ""
)

func main() {
    log.Printf("Starting service 'CharsProcessor' build version: %s, commit: %s", build_version, build_commit)
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    signal.Notify(stop, os.Kill)
    signal.Notify(stop, syscall.SIGTERM)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    configSettings := runtimeconfig.ConfigSettings{}
    service := runtime.MakeService[*charsprocessor.Service, *charsprocessor.ServiceConfig]("CharsProcessor", &configSettings)
    if err := service.StartService(ctx); err != nil {
        log.Fatalln(err)
    }
    log.Infof("Service '%s' started.", "CharsProcessor")
    <-stop
    log.Infof("Service '%s' stop signal received", "CharsProcessor")
    service.StopService(ctx)
    log.Infof("Service '%s' stopped.", "CharsProcessor")
}
