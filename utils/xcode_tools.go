package utils

import "fmt"

// Check if Xcode CLI tools are installed
func checkXcodeTools() bool {

	// Check xcode-select
	_, status, err := ExecuteProcess("xcode-select", "-p")
	if err != nil || status != 0 {
		return false
	}

	// Check pkgutil
	_, status, err = ExecuteProcess("pkgutil", "--pkg-info=com.apple.pkg.DeveloperToolsCLI")
	if err != nil || status != 0 {
		return false
	}

	return true
}

func installXcodeTools() error {
	_, status, err := ExecuteProcess("xcode-select", "--install")
	if err != nil || status != 0 {
		if err == nil {
			err = fmt.Errorf("failed to install Xcode CLI tools")
		}

		return err
	}
	return nil
}

func ManageXcodeTools() error {
	if !checkXcodeTools() {
		fmt.Println("Xcode CLI tools are not installed.")

		// Ask the user if he wants to install the tools
		if AskForConfirmation("Do you want to install them now?") {
			if err := installXcodeTools(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("install: Please install the Xcode command line tools and re-launch this application")
		}

		// Check if the tools are installed
		if !checkXcodeTools() {
			return fmt.Errorf("install: Xcode CLI tools failed to be installed")
		}
	}
	return nil
}
