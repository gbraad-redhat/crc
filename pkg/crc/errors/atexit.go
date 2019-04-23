/*
Copyright (C) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errors

import (
	"os"
	"fmt"

	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/output"
)

const ExitHandlerPanicMessage = "At least one exit handler vetoed to exit program execution"

type exitHandlerFunc func(int) bool
// exitHandlers keeps track of the list of registered exit handlers. Handlers are applied in the order defined in this list.
var exitHandlers = []exitHandlerFunc{}

// Exit runs all registered exit handlers and then exits the program with the specified exit code using os.Exit.
func Exit(code int) {
	veto := runHandlers(code)
	if veto {
		panic(ExitHandlerPanicMessage)
	}
	os.Exit(code)
}

// ExitWithMessage runs all registered exit handlers, prints the specified message and then exits the program with the specified exit code.
// If the exit code is 0, the message is prints to stdout, otherwise to stderr.
func ExitWithMessage(code int, text string, args ...interface{}) {
	if code == 0 {
		output.OutW(os.Stdout, fmt.Sprintf(text, args))
	} else {
		output.OutW(os.Stderr, fmt.Sprintf(text, args))
	}
	Exit(code)
}

// Register registers an exit handler function which is run when Exit is called
func RegisterExitHandler(exitHandler exitHandlerFunc) {
	exitHandlers = append(exitHandlers, exitHandler)
}

// ClearExitHandler clears all registered exit handlers
func ClearExitHandler() {
	exitHandlers = []exitHandlerFunc{}
}

// runHandlers runs all registered exit handlers, passing on the intended exit code.
// Handlers can veto to exit the program. If at least one exit handlers casts a veto, the program does panic instead of exiting
// allowing for a potential recovering.
func runHandlers(code int) bool {
	var veto bool
	for _, handler := range exitHandlers {
		veto = veto || runHandler(handler, code)
	}
	return veto
}

// runHandler runs the single specified exit handler, returning whether this handler vetos the exit or not.
func runHandler(exitHandler exitHandlerFunc, code int) bool {
	defer func() {
		err := recover()
		if err != nil {
			logging.ErrorF("Error running exit handler: %v", err)
		}
	}()

	return exitHandler(code)
}
