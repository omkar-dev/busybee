package busybee

import (
	"errors"
	"fmt"
)

var statusServices = make(map[string]*StatusService)

//Status is an integer representing possable status.
type Status int

const (
	//StatusUnknown is used if a service could not figure out what the status of the user was.
	StatusUnknown Status = iota
	//StatusBusy is used to state that the user is busy.
	StatusBusy
	//StatusAvailable is used to state that the user is Available.
	StatusAvailable
)

//StatusService is the foundation for any service that can update the status of a user.
type StatusService struct {
	Service
	Create StatusServiceCreator
}

//UpdateStatuser is any service that can be used to update your status.
type UpdateStatuser interface {
	UpdateStatus(uid string, status Status) error
}

//StatusServiceCreator is a function that will create a new instance of a UpdateStatuser.
type StatusServiceCreator func(interface{}, *BusyBee) (UpdateStatuser, error)

//RegisterStatusService registers a StatusService instance with BusyBee
func RegisterStatusService(name string, statusService *StatusService) error {
	if name == "" {
		return errors.New("Must supply a name to register a Status Service")
	}
	if statusService == nil {
		return errors.New("Must supply a Status Service to register")
	}
	if _, exists := statusServices[name]; exists {
		return fmt.Errorf("%s is already registered as a Status Service", name)
	}
	statusServices[name] = statusService
	return nil
}

//CreateStatusService returns an instance of a registered StatusService by name.
func (bb *BusyBee) CreateStatusService(name string, conf interface{}) (UpdateStatuser, error) {
	statusService, err := GetStatusService(name)
	if err != nil {
		return nil, err
	}
	return statusService.Create(conf, bb)
}

//ListStatusService returns a slice of all registered status services.
func (bb *BusyBee) ListStatusService() []string {
	tmp := make([]string, 0, len(statusServices))
	for k := range statusServices {
		tmp = append(tmp, k)
	}
	return tmp
}

//GetStatusService will return the struct on how to create a new instance and a new config for that service.
func GetStatusService(name string) (*StatusService, error) {
	status, ok := statusServices[name]
	if !ok {
		return nil, fmt.Errorf("Could not find %s in the list of registered Event Services", name)
	}
	return status, nil
}
