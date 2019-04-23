package machine

type StartConfig struct {
	Name     string

	// CRC system bundle
	BundlePath string

	// Hypervisor
	VMDriver string
	Memory   int
	CPUs	 int

	// Machine log output
	Debug bool
}

type StartResult struct {
	Name   string
	Status string
	Error  string
}

type StopConfig struct {
	Name   string
}

type StopResult struct {
	Name	string
	Success	bool
	Error   string
}

type DeleteConfig struct {
	Name	string
}

type DeleteResult struct {
	Name	string
	Success	bool
	Error   string
}