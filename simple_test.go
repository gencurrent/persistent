package persistent

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func PrintStuff() (interface{}, error) {
	time.Sleep(time.Millisecond * 1500)
	fmt.Printf("The service is running \n")
	time.Sleep(time.Millisecond * 1500)
	return nil, nil
}

func PrintAndWait() (interface{}, error) {
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("The service is running \n")
	time.Sleep(time.Millisecond * 5000)
	return nil, nil
}

func PrintAndError() (interface{}, error) {
	time.Sleep(time.Millisecond * 1500)
	fmt.Printf("The service is running \n")
	a := 1
	b := 0
	c := a / b
	c = c / 2
	return nil, errors.New("A wild error appeared!")
	// return nil, nil
}
func TestAll(t *testing.T) {
	desc1 := ServiceDescriptor{2, loggingDefault, loggingDefault, loggingDefault, loggingDefault}
	desc2 := ServiceDescriptor{2, loggingDefault, loggingDefault, loggingDefault, loggingDefault}
	// desc3 := ServiceDescriptor{2, loggingDefault, loggingDefault, loggingDefault, loggingDefault}

	ser1 := NewService("ServiceOne", PrintStuff, desc1)
	ser2 := NewService("ServiceTwo", PrintAndWait, desc2)
	// ser3 := Service{"Service three", PrintAndError, desc3}

	bundle := ServiceBundle{[]Service{ser1, ser2}}
	bundle.Run()
}
