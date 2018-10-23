package cqrs

import (
	"testing"
)

func TestCallbackEventFactoryMake(t *testing.T) {
	f := NewCallbackEventFactory()

	if err := f.RegisterCallback(eventAFactory); err != nil {
		t.Errorf("expected a nil error but got error : %v", err)
	}

	event := f.MakeEvent("event:a")
	if event == nil {
		t.Fatal("expected the constructed event but got nil instead")
	}

	if string(event.EventType()) != "event:a" {
		t.Errorf("expected an event with name `%v` but got `%v`", "event.a", event.EventType())
	}

	if _, ok := event.(*eventA); !ok {
		t.Errorf("expected an event of type `%v` but got `%T`", "eventA", event)
	}
}

func TestCallbackEventFactoryMakeWithUnknownEvent(t *testing.T) {
	f := NewCallbackEventFactory()

	event := f.MakeEvent("this.event.is.not.registered")
	if event != nil {
		t.Fatalf("expected a nil response but got an event instead `%T`", event)
	}
}

func TestCallbackEventFactoryDuplicateRegister(t *testing.T) {
	f := NewCallbackEventFactory()

	if err := f.RegisterCallback(eventAFactory); err != nil {
		t.Errorf("expected a nil error but got error : %v", err)
	}

	err := f.RegisterCallback(eventAFactory)
	if err == nil {
		t.Error("expected a error but got none")
	}

	if _, ok := err.(ErrorEventFactoryAlreadyRegistered); !ok {
		t.Errorf("expected an error of type `%T` but got `%T`", ErrorEventFactoryAlreadyRegistered(""), err)
	}
}

func TestCallbackEventFactoryNonPointerEventRegister(t *testing.T) {
	f := NewCallbackEventFactory()

	badEventFactory := func() EventData {
		return eventB{}
	}

	err := f.RegisterCallback(badEventFactory)
	if err == nil {
		t.Error("expected an error but got none")
	}

	if _, ok := err.(ErrorEventFactoryNotReturningPointer); !ok {
		t.Errorf("expected an error of type `%T` but got `%T`", ErrorEventFactoryNotReturningPointer(""), err)
	}
}

func TestErrorEventFactoryAlreadyRegisteredToString(t *testing.T) {
	e := ErrorEventFactoryAlreadyRegistered("test")

	expected := "event factory callback/delegate already registered for type: `test`"
	if e.Error() != expected {
		t.Errorf("expected a error message `%s` but got message : `%s`", expected, e.Error())
	}
}

func TestErrorEventFactoryNotReturningPointerToString(t *testing.T) {
	e := ErrorEventFactoryNotReturningPointer("test")

	expected := "event factory callback/delegate does not return a pointer reference for type: `test`"
	if e.Error() != expected {
		t.Errorf("expected a error message `%s` but got message : `%s`", expected, e.Error())
	}
}
