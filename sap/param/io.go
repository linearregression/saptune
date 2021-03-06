package param

import (
	"github.com/HouzuoGuo/saptune/sap"
	"github.com/HouzuoGuo/saptune/system"
	"io/ioutil"
	"path"
)

// Change IO elevators on all IO devices
type BlockDeviceSchedulers struct {
	SchedulerChoice map[string]string
}

func (ioe BlockDeviceSchedulers) Inspect() (Parameter, error) {
	newIOE := BlockDeviceSchedulers{SchedulerChoice: make(map[string]string)}
	// List /sys/block and inspect the IO elevator of each one
	dirContent, err := ioutil.ReadDir("/sys/block")
	if err != nil {
		return nil, err
	}
	for _, entry := range dirContent {
		/*
			Remember: GetSysChoice does not accept the leading /sys/.
			The file "scheduler" may look like "[noop] deadline cfq", in which case the choice will be read successfully.
			If the file simply says "none", which means IO scheduling is not relevant to the block device, then
			the device name will not appear in return value, and there is no point in tuning it anyways.
		*/
		elev, _ := system.GetSysChoice(path.Join("block", entry.Name(), "queue", "scheduler"))
		if elev != "" {
			newIOE.SchedulerChoice[entry.Name()] = elev
		}
	}
	return newIOE, nil
}
func (ioe BlockDeviceSchedulers) Optimise(newElevatorName interface{}) (Parameter, error) {
	newIOE := BlockDeviceSchedulers{SchedulerChoice: make(map[string]string)}
	for k := range ioe.SchedulerChoice {
		newIOE.SchedulerChoice[k] = newElevatorName.(string)
	}
	return newIOE, nil
}
func (ioe BlockDeviceSchedulers) Apply() error {
	errs := make([]error, 0, 0)
	for name, elevator := range ioe.SchedulerChoice {
		errs = append(errs, system.SetSysString(path.Join("block", name, "queue", "scheduler"), elevator))
	}
	err := sap.PrintErrors(errs)
	return err
}

// Change IO nr_requests on all block devices
type BlockDeviceNrRequests struct {
	NrRequests map[string]int
}

func (ior BlockDeviceNrRequests) Inspect() (Parameter, error) {
	newIOR := BlockDeviceNrRequests{NrRequests: make(map[string]int)}
	// List /sys/block and inspect the number of requests of each one
	dirContent, err := ioutil.ReadDir("/sys/block")
	if err != nil {
		return nil, err
	}
	for _, entry := range dirContent {
		// Remember, GetSysString does not accept the leading /sys/
		nrreq, err := system.GetSysInt(path.Join("block", entry.Name(), "queue", "nr_requests"))
		if nrreq >= 0 && err == nil {
			newIOR.NrRequests[entry.Name()] = nrreq
		}
	}
	return newIOR, nil
}
func (ior BlockDeviceNrRequests) Optimise(newNrRequestValue interface{}) (Parameter, error) {
	newIOR := BlockDeviceNrRequests{NrRequests: make(map[string]int)}
	for k := range ior.NrRequests {
		newIOR.NrRequests[k] = newNrRequestValue.(int)
	}
	return newIOR, nil
}
func (ior BlockDeviceNrRequests) Apply() error {
	errs := make([]error, 0, 0)
	for name, nrreq := range ior.NrRequests {
		errs = append(errs, system.SetSysInt(path.Join("block", name, "queue", "nr_requests"), nrreq))
	}
	err := sap.PrintErrors(errs)
	return err
}
