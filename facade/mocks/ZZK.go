package mocks

import "github.com/stretchr/testify/mock"

import "github.com/control-center/serviced/datastore"
import "github.com/control-center/serviced/domain/host"
import "github.com/control-center/serviced/domain/pool"
import "github.com/control-center/serviced/domain/registry"
import "github.com/control-center/serviced/domain/service"
import zkservice "github.com/control-center/serviced/zzk/service"

type ZZK struct {
	mock.Mock
}

func (_m *ZZK) UpdateService(ctx datastore.Context, tenantID string, svc *service.Service, setLockOnCreate bool, setLockOnUpdate bool) error {
	ret := _m.Called(ctx, tenantID, svc, setLockOnCreate, setLockOnUpdate)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, string, *service.Service, bool, bool) error); ok {
		r0 = rf(ctx, tenantID, svc, setLockOnCreate, setLockOnUpdate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) UpdateServices(ctx datastore.Context, tenantID string, svcs []*service.Service, setLockOnCreate bool, setLockOnUpdate bool) error {
	ret := _m.Called(ctx, tenantID, svcs, setLockOnCreate, setLockOnUpdate)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, string, []*service.Service, bool, bool) error); ok {
		r0 = rf(ctx, tenantID, svcs, setLockOnCreate, setLockOnUpdate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) SyncServiceRegistry(ctx datastore.Context, tenantID string, svc *service.Service) error {
	ret := _m.Called(ctx, tenantID, svc)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, string, *service.Service) error); ok {
		r0 = rf(ctx, tenantID, svc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) RemoveService(poolID string, serviceID string) error {
	ret := _m.Called(poolID, serviceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(poolID, serviceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) RemoveServiceEndpoints(serviceID string) error {
	ret := _m.Called(serviceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(serviceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) RemoveTenantExports(tenantID string) error {
	ret := _m.Called(tenantID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(tenantID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) WaitService(svc *service.Service, state service.DesiredState, cancel <-chan interface{}) error {
	ret := _m.Called(svc, state, cancel)

	var r0 error
	if rf, ok := ret.Get(0).(func(*service.Service, service.DesiredState, <-chan interface{}) error); ok {
		r0 = rf(svc, state, cancel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) WaitInstance(ctx datastore.Context, svc *service.Service, instanceID int, checkInstance func(*zkservice.State, bool) bool, cancel <-chan struct{}) error {
	ret := _m.Called(ctx, svc, instanceID, checkInstance, cancel)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, *service.Service, int, func(*zkservice.State, bool) bool, <-chan struct{}) error); ok {
		r0 = rf(ctx, svc, instanceID, checkInstance, cancel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) GetPublicPort(portAddress string) (string, string, error) {
	ret := _m.Called(portAddress)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(portAddress)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(portAddress)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(portAddress)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
func (_m *ZZK) GetVHost(subdomain string) (string, string, error) {
	ret := _m.Called(subdomain)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(subdomain)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(subdomain)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(subdomain)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
func (_m *ZZK) AddHost(_host *host.Host) error {
	ret := _m.Called(_host)

	var r0 error
	if rf, ok := ret.Get(0).(func(*host.Host) error); ok {
		r0 = rf(_host)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) UpdateHost(_host *host.Host) error {
	ret := _m.Called(_host)

	var r0 error
	if rf, ok := ret.Get(0).(func(*host.Host) error); ok {
		r0 = rf(_host)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) RemoveHost(_host *host.Host) error {
	ret := _m.Called(_host)

	var r0 error
	if rf, ok := ret.Get(0).(func(*host.Host) error); ok {
		r0 = rf(_host)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) GetActiveHosts(ctx datastore.Context, poolID string, hosts *[]string) error {
	ret := _m.Called(ctx, poolID, hosts)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, string, *[]string) error); ok {
		r0 = rf(ctx, poolID, hosts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) IsHostActive(poolID string, hostId string) (bool, error) {
	ret := _m.Called(poolID, hostId)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(poolID, hostId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(poolID, hostId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ZZK) UpdateResourcePool(_pool *pool.ResourcePool) error {
	ret := _m.Called(_pool)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pool.ResourcePool) error); ok {
		r0 = rf(_pool)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) RemoveResourcePool(poolID string) error {
	ret := _m.Called(poolID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(poolID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) AddVirtualIP(vip *pool.VirtualIP) error {
	ret := _m.Called(vip)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pool.VirtualIP) error); ok {
		r0 = rf(vip)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) RemoveVirtualIP(vip *pool.VirtualIP) error {
	ret := _m.Called(vip)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pool.VirtualIP) error); ok {
		r0 = rf(vip)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) GetVirtualIPHostID(poolID string, ip string) (string, error) {
	ret := _m.Called(poolID, ip)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(poolID, ip)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(poolID, ip)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ZZK) GetRegistryImage(id string) (*registry.Image, error) {
	ret := _m.Called(id)

	var r0 *registry.Image
	if rf, ok := ret.Get(0).(func(string) *registry.Image); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*registry.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ZZK) SetRegistryImage(rImage *registry.Image) error {
	ret := _m.Called(rImage)

	var r0 error
	if rf, ok := ret.Get(0).(func(*registry.Image) error); ok {
		r0 = rf(rImage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) DeleteRegistryImage(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) DeleteRegistryLibrary(tenantID string) error {
	ret := _m.Called(tenantID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(tenantID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) LockServices(ctx datastore.Context, svcs []service.ServiceDetails) error {
	ret := _m.Called(ctx, svcs)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, []service.ServiceDetails) error); ok {
		r0 = rf(ctx, svcs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) UnlockServices(ctx datastore.Context, svcs []service.ServiceDetails) error {
	ret := _m.Called(ctx, svcs)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, []service.ServiceDetails) error); ok {
		r0 = rf(ctx, svcs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) GetServiceStates(ctx datastore.Context, poolID string, serviceID string) ([]zkservice.State, error) {
	ret := _m.Called(ctx, poolID, serviceID)

	var r0 []zkservice.State
	if rf, ok := ret.Get(0).(func(datastore.Context, string, string) []zkservice.State); ok {
		r0 = rf(ctx, poolID, serviceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]zkservice.State)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(datastore.Context, string, string) error); ok {
		r1 = rf(ctx, poolID, serviceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ZZK) GetHostStates(ctx datastore.Context, poolID string, hostID string) ([]zkservice.State, error) {
	ret := _m.Called(ctx, poolID, hostID)

	var r0 []zkservice.State
	if rf, ok := ret.Get(0).(func(datastore.Context, string, string) []zkservice.State); ok {
		r0 = rf(ctx, poolID, hostID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]zkservice.State)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(datastore.Context, string, string) error); ok {
		r1 = rf(ctx, poolID, hostID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ZZK) GetServiceState(ctx datastore.Context, poolID string, serviceID string, instanceID int) (*zkservice.State, error) {
	ret := _m.Called(ctx, poolID, serviceID, instanceID)

	var r0 *zkservice.State
	if rf, ok := ret.Get(0).(func(datastore.Context, string, string, int) *zkservice.State); ok {
		r0 = rf(ctx, poolID, serviceID, instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*zkservice.State)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(datastore.Context, string, string, int) error); ok {
		r1 = rf(ctx, poolID, serviceID, instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ZZK) StopServiceInstance(poolID string, serviceID string, instanceID int) error {
	ret := _m.Called(poolID, serviceID, instanceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int) error); ok {
		r0 = rf(poolID, serviceID, instanceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) StopServiceInstances(ctx datastore.Context, poolID string, serviceID string) error {
	ret := _m.Called(ctx, poolID, serviceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, string, string) error); ok {
		r0 = rf(ctx, poolID, serviceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) RestartInstance(ctx datastore.Context, poolID string, serviceID string, instanceID int) error {
	ret := _m.Called(ctx, poolID, serviceID, instanceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, string, string, int) error); ok {
		r0 = rf(ctx, poolID, serviceID, instanceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) SendDockerAction(poolID string, serviceID string, instanceID int, command string, args []string) error {
	ret := _m.Called(poolID, serviceID, instanceID, command, args)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int, string, []string) error); ok {
		r0 = rf(poolID, serviceID, instanceID, command, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) GetServiceStateIDs(poolID string, serviceID string) ([]zkservice.StateRequest, error) {
	ret := _m.Called(poolID, serviceID)

	var r0 []zkservice.StateRequest
	if rf, ok := ret.Get(0).(func(string, string) []zkservice.StateRequest); ok {
		r0 = rf(poolID, serviceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]zkservice.StateRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(poolID, serviceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ZZK) GetServiceNodes() ([]zkservice.ServiceNode, error) {
	ret := _m.Called()

	var r0 []zkservice.ServiceNode
	if rf, ok := ret.Get(0).(func() []zkservice.ServiceNode); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]zkservice.ServiceNode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ZZK) RegisterDfsClients(clients ...host.Host) error {
	ret := _m.Called(clients)

	var r0 error
	if rf, ok := ret.Get(0).(func(...host.Host) error); ok {
		r0 = rf(clients...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) UnregisterDfsClients(clients ...host.Host) error {
	ret := _m.Called(clients)

	var r0 error
	if rf, ok := ret.Get(0).(func(...host.Host) error); ok {
		r0 = rf(clients...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ZZK) UpdateInstanceCurrentState(ctx datastore.Context, poolID, serviceID string, instanceID int, state service.InstanceCurrentState) error {
	ret := _m.Called(ctx, poolID, serviceID, instanceID, state)

	var r0 error
	if rf, ok := ret.Get(0).(func(datastore.Context, string, string, int, service.InstanceCurrentState) error); ok {
		r0 = rf(ctx, poolID, serviceID, instanceID, state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
