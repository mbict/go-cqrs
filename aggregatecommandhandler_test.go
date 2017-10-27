package cqrs

import (
	"errors"
	"github.com/satori/go.uuid"
	. "github.com/stretchr/testify/mock"
	"testing"
)

func TestAggregateCommandHandler(t *testing.T) {
	command := &commandA{}

	aggregate := &MockAggregateComposition{}
	aggregate.On("HandleCommand", command).Return(nil)

	repo := &MockAggregateRepository{}
	repo.On("Load", Anything).Return(aggregate, nil)
	repo.On("Save", aggregate).Return(nil)

	err := AggregateCommandHandler(repo).Handle(nil, command)

	if err != nil {
		t.Fatalf("did not expecte a error but got `%v`", err)
	}

	aggregate.AssertNumberOfCalls(t, "HandleCommand", 1)
	aggregate.AssertCalled(t, "HandleCommand", command)

	repo.AssertNumberOfCalls(t, "Save", 1)
	repo.AssertCalled(t, "Save", aggregate)
}

func TestAggregateCommandHandler_NonAggregateCommand(t *testing.T) {
	repo := &MockAggregateRepository{}
	repo.On("Load", Anything).Return(nil, nil)
	command := &commandNonAggregate{}

	err := AggregateCommandHandler(repo).Handle(nil, command)

	if err == nil {
		t.Fatal("expected a error but got nil instead")
	}

	err, ok := err.(ErrorNotAnAggregateCommand)
	if !ok {
		t.Fatalf("expected error of type `%T` but got `%T`", "ErrorAggregateNotFound", err)
	}

	expectedErr := ErrorNotAnAggregateCommand(command.CommandName()).Error()
	if err.Error() != expectedErr {
		t.Errorf("expected error with message `%s` but got `%s`", expectedErr, err.Error())
	}
}

func TestAggregateCommandHandler_AggregateNotFound(t *testing.T) {
	repo := &MockAggregateRepository{}
	repo.On("Load", Anything).Return(nil, nil)
	command := &commandA{Id: uuid.NewV4()}

	err := AggregateCommandHandler(repo).Handle(nil, command)

	if err == nil {
		t.Fatal("expected a error but got nil instead")
	}

	err, ok := err.(ErrorAggregateNotFound)
	if !ok {
		t.Fatalf("expected error of type `%T` but got `%T`", "ErrorAggregateNotFound", err)
	}

	expectedErr := ErrorAggregateNotFound(command.Id.String()).Error()
	if err.Error() != expectedErr {
		t.Errorf("expected error with message `%s` but got `%s`", expectedErr, err.Error())
	}
}

func TestAggregateCommandHandler_RepositoryLoadError(t *testing.T) {
	repoError := errors.New("repo error")
	repo := &MockAggregateRepository{}
	repo.On("Load", Anything).Return(nil, repoError)
	command := &commandA{Id: uuid.NewV4()}

	err := AggregateCommandHandler(repo).Handle(nil, command)

	if err == nil {
		t.Fatal("expected a error but got nil instead")
	}

	if err != repoError {
		t.Fatalf("expected error of type `%v` but got `%v`", repoError, err)
	}

	repo.AssertCalled(t, "Load", command.Id)
	repo.AssertNotCalled(t, "Save")
}

func TestAggregateCommandHandler_RepositorySaveError(t *testing.T) {
	command := &commandA{}
	repoError := errors.New("repo error")

	aggregate := &MockAggregateComposition{}
	aggregate.On("HandleCommand", command).Return(nil)

	repo := &MockAggregateRepository{}
	repo.On("Load", Anything).Return(aggregate, nil)
	repo.On("Save", aggregate).Return(repoError)

	err := AggregateCommandHandler(repo).Handle(nil, command)

	if err == nil {
		t.Fatal("expected a error but got nil instead")
	}

	if err != repoError {
		t.Fatalf("expected error of type `%v` but got `%v`", repoError, err)
	}

	aggregate.AssertNumberOfCalls(t, "HandleCommand", 1)
	aggregate.AssertCalled(t, "HandleCommand", command)

	repo.AssertNumberOfCalls(t, "Save", 1)
	repo.AssertCalled(t, "Save", aggregate)
}

func TestAggregateCommandHandler_CommandLogicError(t *testing.T) {
	command := &commandA{}
	commandError := errors.New("command error")

	aggregate := &MockAggregateComposition{}
	aggregate.On("HandleCommand", command).Return(commandError)

	repo := &MockAggregateRepository{}
	repo.On("Load", Anything).Return(aggregate, nil)

	err := AggregateCommandHandler(repo).Handle(nil, command)

	if err == nil {
		t.Fatal("expected a error but got nil instead")
	}

	if err != commandError {
		t.Fatalf("expected error of type `%v` but got `%v`", commandError, err)
	}

	aggregate.AssertNumberOfCalls(t, "HandleCommand", 1)
	aggregate.AssertCalled(t, "HandleCommand", command)
	repo.AssertNotCalled(t, "Save", aggregate)
}

func TestAggregateCommandHandler_ValidateError(t *testing.T) {
	validateError := errors.New("validate error")
	validate := &MockValidate{}
	validate.On("Validate").Return(validateError)
	command := &commandAWithValidate{ValidateMock: validate}
	aggregate := &MockAggregateComposition{}
	aggregate.On("Validate").Return(validateError)
	repo := &MockAggregateRepository{}
	repo.On("Load", Anything).Return(aggregate, nil)

	err := AggregateCommandHandler(repo).Handle(nil, command)

	if err == nil {
		t.Fatal("expected a error but got nil instead")
	}

	if err != validateError {
		t.Fatalf("expected error of type `%v` but got `%v`", validateError, err)
	}

	validate.AssertNumberOfCalls(t, "Validate", 1)
	aggregate.AssertNotCalled(t, "HandleCommand", command)
	repo.AssertNotCalled(t, "Save", aggregate)
}
