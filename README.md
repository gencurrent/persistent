# Persistent #

Make your Service `Persistence` and `Runnable` at background. 
The Persistance lib alows you not to worry about times your service must be **re/started at parallel threads**.

## Usage ##
### 1. Install it ###

```sh
go get github.com/gencurrent/persistent
```

### 2. Import Persistent ###

```go
import (
    persistent "github.com/gencurrent/persistent"
)
```

### 3. Describe your services ### 

```go
// Service descriptors
desc1 := persistent.ServiceDescriptor{2, persistent.loggingDefault, persistent.loggingDefault, persistent.loggingDefault, persistent.loggingDefault} // Should remove loggingDefault repetitions ... [TODO: doit]
desc2 := persistent.ServiceDescriptor{2, persistent.loggingDefault, persistent.loggingDefault, persistent.loggingDefault, persistent.loggingDefault}

// Services themselves
// (YourServiceStarterFunc1 is function with return type itnerface{}, error)
ser1 := persistent.NewService("ServiceOne", YourServiceStarterFunc1, desc1)
ser2 := persistent.NewService("ServiceTwo", YourServiceStarterFunc2, desc2)

// Bind them all in a single bundle
bundle := persistent.ServiceBundle{[]persistent.Service{ser1, ser2}}

// Run it
bundle.Run()
```

### 4. Enjoy your services ###

The example of console output

```sh
Initialization loop #0 has started
Initialization loop #1 has started
INFO[0000] Success                                       service=ServiceTwo
INFO[0000] Success                                       service=ServiceOne
The service is running 
The service is running 
INFO[0003] About to restart ServiceOne for the #0 time   service=ServiceOne
The service is running 
INFO[0005] About to restart ServiceTwo for the #0 time   service=ServiceTwo
The service is running 
INFO[0006] Stopped                                       service=ServiceOne
INFO[0011] Stopped                                       service=ServiceTwo
PASS
```

## BTW ##
<details markdown="1">
Contact me
<summary>
My first OpenSource, so, please, do not kick me hard (or Persistently)
if have qustions, please send me an email: <gencurrent@gmail.com>
</summary>
