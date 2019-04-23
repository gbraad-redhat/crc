package machine

import (
	"encoding/json"
	"io/ioutil"

	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/errors"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/output"

	"github.com/code-ready/crc/pkg/crc/machine/config"
	"github.com/code-ready/crc/pkg/crc/machine/libvirt"

	"github.com/code-ready/machine/libmachine"
	"github.com/code-ready/machine/libmachine/log"
	"github.com/code-ready/machine/libmachine/state"
	"github.com/code-ready/machine/libmachine/host"
)

func init() {
}

func Start(startConfig StartConfig) (StartResult, error) {
	result := &StartResult{Name: startConfig.Name}

	// Set libmachine logging
	setMachineLogging(startConfig.Debug)

	libMachineAPIClient := libmachine.NewClient(constants.MachineBaseDir, constants.MachineCertsDir)
	defer libMachineAPIClient.Close()

	machineConfig := config.MachineConfig{
		Name:       startConfig.Name,
		BundlePath: startConfig.BundlePath,
		VMDriver:   startConfig.VMDriver,
		CPUs:	    startConfig.CPUs,
		Memory:     startConfig.Memory,
	}

	exists, err := existVM(libMachineAPIClient, machineConfig);
	if !exists {
		output.Out("Creating VM")

		host, err := createHost(libMachineAPIClient, machineConfig)
		if err != nil {
			logging.ErrorF("Error creating host: %s", err)
			result.Error = err.Error()
		}

		vmState, err := host.Driver.GetState()
		if err != nil {
			logging.ErrorF("Error getting the state for host: %s", err)
			result.Error = err.Error()
		}
		result.Status = vmState.String()

		if vmState != state.Running {
			host.Driver.Start()
		}
	} else {
		output.Out("Starting stopped VM")
		host, err := libMachineAPIClient.Load(machineConfig.Name)
		s, err := host.Driver.GetState()
		if err != nil {
			logging.ErrorF("Error getting the state for host: %s", err)
			result.Error = err.Error()
		}

		if s != state.Running {
			if err := host.Driver.Start(); err != nil {
				logging.ErrorF("Error starting stopped VM: %s", err)
				result.Error = err.Error()
			}
			if err := libMachineAPIClient.Save(host); err != nil {
				logging.ErrorF("Error saving state for VM: %s", err)
				result.Error = err.Error()
			}
		}

		vmState, err := host.Driver.GetState()
		if err != nil {
			logging.ErrorF("Error getting the state for host: %s", err)
			result.Error = err.Error()
		}
		result.Status = vmState.String()
	}

	return *result, err
}

func Stop(stopConfig StopConfig) (StopResult, error) {
	result := &StopResult{Name: stopConfig.Name}

	libMachineAPIClient := libmachine.NewClient(constants.MachineBaseDir, constants.MachineCertsDir)
	host, err := libMachineAPIClient.Load(stopConfig.Name)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return *result, err
	}

	if err := host.Stop(); err != nil {
		result.Success = false
		result.Error = err.Error()
		return *result, err
	}

	result.Success = true
	return *result, nil
}

func Delete(deleteConfig DeleteConfig) (DeleteResult, error) {
	result := &DeleteResult{Name: deleteConfig.Name, Success: true}

	libMachineAPIClient := libmachine.NewClient(constants.MachineBaseDir, constants.MachineCertsDir)
	host, err := libMachineAPIClient.Load(deleteConfig.Name)
	
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return *result, err
	}

	m := errors.MultiError{}
	m.Collect(host.Driver.Remove())
	m.Collect(libMachineAPIClient.Remove(deleteConfig.Name))

	if m.ToError != nil {
		result.Success = false
		result.Error = m.ToError().Error()
		return *result, m.ToError()		
	}
	return *result, nil
}

func existVM(api libmachine.API, machineConfig config.MachineConfig) (bool, error) {
	exists, err := api.Exists(machineConfig.Name)
	if err != nil {
		return false, errors.NewF("Error checking if the host exists: %s", err)
	}
	return exists, nil
}

func createHost(api libmachine.API, machineConfig config.MachineConfig) (*host.Host, error) {
	driverOptions := getDriverOptions(machineConfig)
	jsonDriverConfig, err := json.Marshal(driverOptions)

	vm, err := api.NewHost(machineConfig.VMDriver, jsonDriverConfig)

	if err != nil {
		return nil, errors.NewF("Error creating new host: %s", err)
	}

	if err := api.Create(vm); err != nil {
		return nil, errors.NewF("Error creating the VM. %s", err)
	}

	return vm, nil
}

func getDriverOptions(machineConfig config.MachineConfig) interface{} {
	var driver interface{}

	// Supported drivers
	switch machineConfig.VMDriver {
	
	case "libvirt":
		driver = libvirt.CreateHost(machineConfig)
	
	default:
		errors.ExitWithMessage(1, "Unsupported driver: %s", machineConfig.VMDriver)
	}
	
	return driver
}

func setMachineLogging(logging bool) {
	if !logging {
		log.SetOutWriter(ioutil.Discard)
		log.SetErrWriter(ioutil.Discard)
	} else {
		log.SetDebug(true)
	}
}