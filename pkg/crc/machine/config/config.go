package config

type MachineConfig struct {
	// CRC system bundle
	BundlePath string

	// Hypervisor
	VMDriver string

	// Virtual machine configuration
	Name     string
	Memory   int
	CPUs     int
}
