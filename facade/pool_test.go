// Copyright 2014, The Serviced Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package facade

import (
	"github.com/zenoss/serviced/domain/host"
	"github.com/zenoss/serviced/domain/pool"
	. "gopkg.in/check.v1"
	"time"
)

func (ft *FacadeTest) Test_NewResourcePool(t *C) {
	poolID := "Test_NewResourcePool"
	defer ft.Facade.RemoveResourcePool(ft.CTX, poolID)

	rp := pool.ResourcePool{}
	err := ft.Facade.AddResourcePool(ft.CTX, &rp)
	if err == nil {
		t.Errorf("Expected failure to create resource pool %-v", rp)
	}

	rp.ID = poolID
	err = ft.Facade.AddResourcePool(ft.CTX, &rp)
	if err != nil {
		t.Errorf("Failure creating resource pool %-v with error: %s", rp, err)
		t.Fail()
	}

	err = ft.Facade.AddResourcePool(ft.CTX, &rp)
	if err == nil {
		t.Errorf("Expected error creating redundant resource pool %-v", rp)
		t.Fail()
	}
}

func (ft *FacadeTest) Test_UpdateResourcePool(t *C) {
	poolID := "Test_UpdateResourcePool"
	defer ft.Facade.RemoveResourcePool(ft.CTX, poolID)

	myPool := pool.New(poolID)
	ft.Facade.AddResourcePool(ft.CTX, myPool)

	myPool.Priority = 1
	myPool.CoreLimit = 1
	myPool.MemoryLimit = 1
	err := ft.Facade.UpdateResourcePool(ft.CTX, myPool)
	if err != nil {
		t.Errorf("Failure updating resource pool %-v with error: %s", myPool, err)
		t.Fail()
	}

	result, err := ft.Facade.GetResourcePool(ft.CTX, poolID)
	result.CreatedAt = myPool.CreatedAt
	result.UpdatedAt = myPool.UpdatedAt
	if myPool.Equals(result) {
		t.Errorf("%+v != %+v", myPool, result)
		t.Fail()
	}
}

func (ft *FacadeTest) Test_GetResourcePool(t *C) {
	poolID := "Test_UpdateResourcePool"
	defer ft.Facade.RemoveResourcePool(ft.CTX, poolID)

	ft.Facade.RemoveResourcePool(ft.CTX, poolID)
	rp := pool.New(poolID)
	rp.Priority = 1
	rp.CoreLimit = 1
	rp.MemoryLimit = 1
	if err := ft.Facade.AddResourcePool(ft.CTX, rp); err != nil {
		t.Fatalf("Failed to add resource pool: %v", err)
	}

	result, err := ft.Facade.GetResourcePool(ft.CTX, poolID)
	result.CreatedAt = rp.CreatedAt
	result.UpdatedAt = rp.UpdatedAt
	if err == nil {
		if rp.Equals(result) {
			t.Errorf("Unexpected ResourcePool: expected=%+v, actual=%+v", rp, result)
		}
	} else {
		t.Errorf("Unexpected Error Retrieving ResourcePool: %v", err)
	}
}

func (ft *FacadeTest) Test_RemoveResourcePool(t *C) {

	poolID := "Test_RemoveResourcePool"
	result, err := ft.Facade.GetResourcePool(ft.CTX, poolID)
	t.Assert(err, IsNil)
	t.Assert(result, IsNil)
	err = ft.Facade.RemoveResourcePool(ft.CTX, poolID)
	t.Assert(err, IsNil)

	rp := pool.New(poolID)
	err = ft.Facade.AddResourcePool(ft.CTX, rp)
	t.Assert(err, IsNil)

	rp, err = ft.Facade.GetResourcePool(ft.CTX, poolID)
	t.Assert(err, IsNil)
	t.Assert(rp.ID, Equals, poolID)

	err = ft.Facade.RemoveResourcePool(ft.CTX, poolID)
	t.Assert(err, IsNil)
	rp, err = ft.Facade.GetResourcePool(ft.CTX, poolID)
	t.Assert(err, IsNil)
	t.Assert(rp, IsNil)
}

func (ft *FacadeTest) Test_GetResourcePools(t *C) {
	result, err := ft.Facade.GetResourcePools(ft.CTX)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		pools := make([]pool.ResourcePool, len(result))
		for i, pool := range result {
			pools[i] = *pool
		}

		t.Fatalf("unexpected pool found: %v", pools)
	}

	poolID := "Test_GetResourcePools"
	defer ft.Facade.RemoveResourcePool(ft.CTX, poolID)
	rp := pool.New(poolID)
	rp.Priority = 1
	rp.CoreLimit = 2
	rp.MemoryLimit = 3
	ft.Facade.AddResourcePool(ft.CTX, rp)
	time.Sleep(time.Second * 2)
	result, err = ft.Facade.GetResourcePools(ft.CTX)
	if err == nil && len(result) == 1 {
		result[0].CreatedAt = rp.CreatedAt
		result[0].UpdatedAt = rp.UpdatedAt
		if result[0].Equals(rp) {
			t.Fatalf("expected [%+v] actual=%s", rp, result)
		}
	} else {
		t.Fatalf("Unexpected Error Retrieving ResourcePools: %v", err)
	}
}

func (ft *FacadeTest) Test_GetPoolsIPs(t *C) {
	assignIPsPool := pool.New("assignIPsPoolID")
	err := ft.Facade.AddResourcePool(ft.CTX, assignIPsPool)
	defer func() {
		ft.Facade.RemoveResourcePool(ft.CTX, assignIPsPool.ID)
	}()

	if err != nil {
		t.Errorf("Failure creating resource pool %-v with error: %s", assignIPsPool, err)
		t.Fail()
	}

	hostID := "assignIPsHost"
	ipAddress1 := "192.168.100.10"
	ipAddress2 := "10.50.9.1"

	assignIPsHostIPResources := []host.HostIPResource{}
	oneHostIPResource := host.HostIPResource{}
	oneHostIPResource.HostID = hostID
	oneHostIPResource.IPAddress = ipAddress1
	oneHostIPResource.InterfaceName = "eth0"
	assignIPsHostIPResources = append(assignIPsHostIPResources, oneHostIPResource)
	oneHostIPResource.HostID = "A"
	oneHostIPResource.IPAddress = ipAddress2
	oneHostIPResource.InterfaceName = "eth1"
	assignIPsHostIPResources = append(assignIPsHostIPResources, oneHostIPResource)

	assignIPsHost, err := host.Build("", assignIPsPool.ID, []string{}...)
	if err != nil {
		t.Fatalf("could not build host for test: %v", err)
	}
	assignIPsHost.ID = hostID
	assignIPsHost.PoolID = assignIPsPool.ID
	assignIPsHost.IPs = assignIPsHostIPResources
	err = ft.Facade.AddHost(ft.CTX, assignIPsHost)
	if err != nil {
		t.Fatalf("failed to add host: %v", err)
	}
	defer func() {
		ft.Facade.RemoveHost(ft.CTX, assignIPsHost.ID)
	}()
	time.Sleep(2 * time.Second)
	IPs, err := ft.Facade.GetPoolIPs(ft.CTX, assignIPsPool.ID)
	if err != nil {
		t.Error("GetPoolIps failed")
	}
	if len(IPs.HostIPs) != 2 {
		t.Fatalf("Expected 2 addresses, found %v", len(IPs.HostIPs))
	}

	if IPs.HostIPs[0].IPAddress != ipAddress1 {
		t.Errorf("Unexpected IP address: %v", IPs.HostIPs[0].IPAddress)
	}
	if IPs.HostIPs[1].IPAddress != ipAddress2 {
		t.Errorf("Unexpected IP address: %v", IPs.HostIPs[1].IPAddress)
	}

}

func (ft *FacadeTest) Test_AddVirtualIP(t *C) {
	myPoolID := "aPoolID"
	assignIPsPool := pool.New(myPoolID)
	err := ft.Facade.AddResourcePool(ft.CTX, assignIPsPool)
	defer func() {
		ft.Facade.RemoveResourcePool(ft.CTX, assignIPsPool.ID)
	}()

	if err != nil {
		t.Errorf("Failure creating resource pool %-v with error: %s", assignIPsPool, err)
		t.Fail()
	}

	hostID := "aHost"
	ipAddress1 := "192.168.100.10"

	assignIPsHostIPResources := []host.HostIPResource{}
	oneHostIPResource := host.HostIPResource{}
	oneHostIPResource.HostID = hostID
	oneHostIPResource.IPAddress = ipAddress1
	myInterfaceName := "eth0"
	oneHostIPResource.InterfaceName = myInterfaceName
	assignIPsHostIPResources = append(assignIPsHostIPResources, oneHostIPResource)

	assignIPsHost, err := host.Build("", assignIPsPool.ID, []string{}...)
	if err != nil {
		t.Fatalf("could not build host for test: %v", err)
	}
	assignIPsHost.ID = hostID
	assignIPsHost.PoolID = assignIPsPool.ID
	assignIPsHost.IPs = assignIPsHostIPResources
	err = ft.Facade.AddHost(ft.CTX, assignIPsHost)
	if err != nil {
		t.Fatalf("failed to add host: %v", err)
	}
	defer func() {
		ft.Facade.RemoveHost(ft.CTX, assignIPsHost.ID)
	}()
	time.Sleep(2 * time.Second)
	if err := ft.Facade.AddVirtualIP(ft.CTX, pool.VirtualIP{PoolID: myPoolID, IP: "192.168.100.20", Netmask: "255.255.255.0", BindInterface: myInterfaceName, InterfaceIndex: ""}); err != nil {
		t.Error("AddVirtualIP failed: %v", err)
	}
	IPs, err := ft.Facade.GetPoolIPs(ft.CTX, assignIPsPool.ID)
	if err != nil {
		t.Error("GetPoolIps failed")
	}
	if len(IPs.VirtualIPs) != 1 {
		t.Fatalf("Expected 1 address, found %v", len(IPs.VirtualIPs))
	}
	if err := ft.Facade.AddVirtualIP(ft.CTX, pool.VirtualIP{PoolID: myPoolID, IP: "192.168.100.30", Netmask: "255.255.255.0", BindInterface: myInterfaceName, InterfaceIndex: ""}); err != nil {
		t.Error("AddVirtualIP failed: %v", err)
	}
	IPs, err := ft.Facade.GetPoolIPs(ft.CTX, assignIPsPool.ID)
	if err != nil {
		t.Error("GetPoolIps failed")
	}
	if len(IPs.VirtualIPs) != 2 {
		t.Fatalf("Expected 2 address, found %v", len(IPs.VirtualIPs))
	}
}

func (ft *FacadeTest) Test_PoolCapacity(t *C) {
	hostid := "host-id"
	poolid := "pool-id"

	//create pool for test
	rp := pool.New(poolid)
	if err := ft.Facade.AddResourcePool(ft.CTX, rp); err != nil {
		t.Fatalf("Could not add pool for test: %v", err)
	}

	//fill host with required values
	h, err := host.Build("", poolid, []string{}...)
	h.ID = hostid
	if err != nil {
		t.Fatalf("Unexpected error building host: %v", err)
	}

	err = ft.Facade.AddHost(ft.CTX, h)
	if err != nil {
		t.Errorf("Unexpected error adding host: %v", err)
	}

	// load pool with calculated capacity
	loadedPool, err := ft.Facade.GetResourcePool(ft.CTX, poolid)

	if err != nil {
		t.Fatalf("Unexpected error calculating pool capacity: %v", err)
	}

	if loadedPool.CoreCapacity <= 0 || loadedPool.MemoryCapacity <= 0 {
		t.Fatalf("Unexpected values calculated for %s capacity: CPU - %v : Memory - %v", loadedPool.ID, loadedPool.CoreCapacity, loadedPool.MemoryCapacity)
	}
}

func (ft *FacadeTest) Test_PoolCommitment(t *C) {
	poolid := "pool-id"

	//create pool for test
	rp := pool.New(poolid)
	if err := ft.Facade.AddResourcePool(ft.CTX, rp); err != nil {
		t.Fatalf("Could not add pool for test: %v", err)
	}

	// load pool with calculated capacity
	loadedPool, err := ft.Facade.GetResourcePool(ft.CTX, poolid)

	commitmentErr := ft.Facade.calcPoolCommitment(ft.CTX, loadedPool)

	if commitmentErr != nil {
		t.Fatalf("Unexpected error calculating pool commitment: %v", err)
	}
}
