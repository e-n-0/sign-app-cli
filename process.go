package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Function that executes a process with given arguments
func executeProcess(args []string) ([]byte, int, error) {
	if len(args) == 0 {
		return nil, -1, fmt.Errorf("executeProcess: no program given")
	}

	// Create a new process
	process := exec.Command(args[0], args[1:]...)

	// Set the process attributes
	//process.Stdout = stdout
	process.Stderr = os.Stderr

	bytes, err := process.Output()
	if err != nil {
		return nil, process.ProcessState.ExitCode(), err
	}

	// Start the process
	return bytes, process.ProcessState.ExitCode(), nil
}
