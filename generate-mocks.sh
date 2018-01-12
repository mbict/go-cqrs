#!/usr/bin/env bash
mockery -name AggregateRepository -testonly -inpkg -case=underscore
mockery -name Aggregate -testonly -inpkg -case=underscore
mockery -name AggregateHandlesCommands -testonly -inpkg -case=underscore
mockery -name AggregateContext -testonly -inpkg -case=underscore
mockery -name Validate -testonly -inpkg -case=underscore
mockery -name EventFactory -testonly -inpkg -case=underscore
mockery -name EventStore -testonly -inpkg -case=underscore
mockery -name EventStream -testonly -inpkg -case=underscore
mockery -name Event -testonly -inpkg -case=underscore
