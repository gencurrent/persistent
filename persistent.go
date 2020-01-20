package persistent

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Copied from logrus logging library
const (
	PanicLevel int = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

const (
	EventOnToRestart string = "onToRestart"
	EventOnStart     string = "onStart"
	EventOnStop      string = "onStop"
	EventOnFail      string = "onFail"
)

const (
	StatusNew int = iota
	StatusRunning
	StatusDone
	StatusFail
)

type ServiceFunction func() (interface{}, error)

type LoggingFunction func(level int, eventType string, service *Service, message string)

func loggingDefault(level int, eventType string, service *Service, message string) {
	switch uint32(level) {
	case 0:
		log.WithFields(log.Fields{"service": service.Name}).Panic(message)
	case 1:
		log.WithFields(log.Fields{"service": service.Name}).Fatal(message)
	case 2:
		log.WithFields(log.Fields{"service": service.Name}).Error(message)
	case 3:
		log.WithFields(log.Fields{"service": service.Name}).Warn(message)
	case 4:
		log.WithFields(log.Fields{"service": service.Name}).Info(message)
	case 5:
		log.WithFields(log.Fields{"service": service.Name}).Debug(message)
	case 6:
		log.WithFields(log.Fields{"service": service.Name}).Trace(message)
	}
}

// Descriptor to use to describeevery service running behaviour
type ServiceDescriptor struct {
	TimesToRestart int             // Times to tolerate continious failures to try before giving up. -1: restart always
	OnToRestart       LoggingFunction // Event onToRestart
	OnStart           LoggingFunction // Event onStart
	OnStop            LoggingFunction // Event onStop
	OnFail            LoggingFunction // Event onToleration: the function failed to run until timeout
}

// Every service running descriptor
type Service struct {
	Name           string
	Function       ServiceFunction
	Descriptor     ServiceDescriptor
	Status         int // One of persistent.go:34
	TimesRestarted int // Times the service restarted at any reason
	TimesHealed    int // Times the service had been restarted due to fails

}

// NewService creates services just with 3 arguments: {name, function, descriptor} 
func NewService(name string, function ServiceFunction, descriptor ServiceDescriptor) Service {
	return Service{
		name,
		function,
		descriptor,
		StatusNew,
		0,
		0,
	}
}

// Just an array for services
type ServiceBundle struct {
	Services []Service
}

// Run the main loop to execute all the ServiceBundle Services
func (bundle ServiceBundle) Run() {
	var waitGroup sync.WaitGroup
	start := make(chan struct{})

	services := bundle.Services

	waitGroup.Add(len(services))
	for i := 0; i < len(services); i++ {
		fmt.Printf("Initialization loop #%v has started\n", i)

		go func(number int, wg *sync.WaitGroup) {
			defer wg.Done()
			service := services[number]
			service.Descriptor.OnStart(InfoLevel, EventOnStart, &service, "Success")
			startedTimes := 0
			for true {
				if startedTimes != 0 {
					service.Descriptor.OnToRestart(InfoLevel, EventOnToRestart, &service, fmt.Sprintf("About to restart for the #%v time", startedTimes-1))
				}
				_, err := service.Function()
				startedTimes++
				if err != nil {
					service.Descriptor.OnFail(WarnLevel, EventOnFail, &service, err.Error())
				}
				if service.Descriptor.TimesToRestart != -1 {
					if startedTimes > service.Descriptor.TimesToRestart - 1{
						break
					}
				}

			}
			service.Descriptor.OnStop(InfoLevel, EventOnStop, &service, "Stopped")
		}(i, &waitGroup)
	}
	waitGroup.Wait()

	close(start)

}
