// Copyright 2016 The Serviced Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/control-center/serviced/commons/docker"
	"github.com/control-center/serviced/dfs/registry"
	"github.com/control-center/serviced/domain/service"
	"github.com/control-center/serviced/zzk"
	zkservice "github.com/control-center/serviced/zzk/service2"
	dockerclient "github.com/fsouza/go-dockerclient"
)

// StopContainer stops running container or returns nil if the container does
// not exist or has already stopped.
func (a *HostAgent) StopContainer(serviceID string, instanceID int) error {
	logger := plog.WithFields(log.Fields{
		"serviceid":  serviceID,
		"instanceid": instanceID,
	})

	// find the container by name
	ctrName := fmt.Sprintf("%s-%d", serviceID, instanceID)
	ctr, err := docker.FindContainer(ctrName)
	if err == docker.ErrNoSuchContainer {
		logger.Debug("Could not stop, container not found")
		return nil
	} else if err != nil {
		logger.WithError(err).Debug("Could not look up container")
		return err
	}

	err = ctr.Stop(45 * time.Second)
	if _, ok := err.(*dockerclient.ContainerNotRunning); ok {
		logger.Debug("Container already stopped")
		return nil
	} else if err != nil {
		logger.WithError(err).Debug("Could not stop container")
		return err
	}

	return nil
}

// AttachContainer returns a channel that monitors the run state of a given
// container.
func (a *HostAgent) AttachContainer(containerID, serviceID string, instanceID int) (<-chan time.Time, error) {
	logger := plog.WithFields(log.Fields{
		"serviceid":   serviceID,
		"instanceid":  instanceID,
		"containerid": containerID,
	})

	// find the container by name
	ctrName := fmt.Sprintf("%s-%d", serviceID, instanceID)
	ctr, err := docker.FindContainer(ctrName)
	if err == docker.ErrNoSuchContainer {
		return nil, nil
	} else if err != nil {
		logger.WithError(err).Debug("Could not look up container")
		return nil, err
	}

	// verify that the container ids match, otherwise delete
	// the container.
	if ctr.ID != containerID {
		ctr.Kill()
		if err := ctr.Delete(true); err != nil {
			logger.WithError(err).Debug("Could not delete orphaned container")
			return nil, err
		}
		logger.WithField("currentcontainerid", ctr.ID).Warn("Removed orphaned container")
		return nil, nil
	}

	// monitor the container
	ev := a.monitorContainer(logger, ctr)

	// make sure the container is running at the time this event is set
	if !ctr.IsRunning() {
		logger.Debug("Could not capture event, container not running")
		ctr.CancelOnEvent(docker.Die)
		return nil, nil
	}
	return ev, nil
}

// StartContainer creates a new container and starts.  It returns info about
// the container, and an event monitor to track the running state of the
// service.
func (a *HostAgent) StartContainer(cancel <-chan interface{}, svc *service.Service, instanceID int) (*zkservice.ServiceState, <-chan time.Time, error) {
	logger := plog.WithFields(log.Fields{
		"serviceid":   svc.ID,
		"servicename": svc.Name,
		"imageid":     svc.ImageID,
		"instanceid":  instanceID,
	})

	// pull the service image
	imageUUID, imageName, err := a.pullImage(logger, cancel, svc.ImageID)
	if err != nil {
		logger.WithError(err).Debug("Could not pull the service image")
		return nil, nil, err
	}
	svc.ImageID = imageName

	// Establish a connection to the master
	// TODO: use the new rpc calls instead
	client, err := NewControlClient(a.master)
	if err != nil {
		logger.WithField("client", a.master).WithError(err).Debug("Could not connect to the master")
		return nil, nil, err
	}
	defer client.Close()

	// get the container configs
	conf, hostConf, err := a.setupContainer(client, svc, instanceID)
	if err != nil {
		logger.WithError(err).Debug("Could not setup container")
		return nil, nil, err
	}

	// create the container
	opts := dockerclient.CreateContainerOptions{
		Name:       fmt.Sprintf("%s-%d", svc.ID, instanceID),
		Config:     conf,
		HostConfig: hostConf,
	}

	ctr, err := docker.NewContainer(&opts, false, 10*time.Second, nil, nil)
	if err != nil {
		logger.WithError(err).Debug("Could not create container")
		return nil, nil, err
	}
	logger = logger.WithField("containerid", ctr.ID)
	logger.Debug("Created a new container")

	// start the container
	ev := a.monitorContainer(logger, ctr)

	if err := ctr.Start(); err != nil {
		logger.WithError(err).Debug("Could not start container")
		ctr.CancelOnEvent(docker.Die)
		return nil, nil, err
	}
	logger.Debug("Started container")

	dctr, err := ctr.Inspect()
	if err != nil {
		logger.WithError(err).Debug("Could not inspect container")
		ctr.CancelOnEvent(docker.Die)
		return nil, nil, err
	}

	state := &zkservice.ServiceState{
		ContainerID: ctr.ID,
		ImageID:     imageUUID,
		Paused:      false,
		PrivateIP:   ctr.NetworkSettings.IPAddress,
		HostIP:      a.ipaddress,
		Started:     dctr.State.StartedAt,
	}

	for _, ep := range svc.Endpoints {
		if ep.Purpose == "export" {
			state.Exports = append(state.Exports, zkservice.ExportBinding{
				Application: ep.Application,
				Protocol:    ep.Protocol,
				PortNumber:  ep.PortNumber,
			})
		} else {
			state.Imports = append(state.Imports, zkservice.ImportBinding{
				Application:    ep.Application,
				Purpose:        ep.Purpose,
				PortNumber:     ep.PortNumber,
				PortTemplate:   ep.PortTemplate,
				VirtualAddress: ep.VirtualAddress,
			})
		}
	}

	return state, ev, nil
}

// ResumeContainer resumes a paused container
func (a *HostAgent) ResumeContainer(svc *service.Service, instanceID int) error {
	logger := plog.WithFields(log.Fields{
		"serviceid":     svc.ID,
		"servicename":   svc.Name,
		"resumecommand": svc.Snapshot.Resume,
		"instanceid":    instanceID,
	})
	ctrName := fmt.Sprintf("%s-%d", svc.ID, instanceID)

	// check to see if the container exists and is running
	ctr, err := docker.FindContainer(ctrName)
	if err == docker.ErrNoSuchContainer {
		// container has been deleted and the event monitor should catch this
		logger.Debug("Container not found")
		return nil
	}
	if !ctr.IsRunning() {
		// container has stopped and the event monitor should catch this
		logger.Debug("Container stopped")
		return nil
	}

	// resume the paused container
	if err := attachAndRun(ctrName, svc.Snapshot.Resume); err != nil {
		logger.WithError(err).Debug("Could not resume paused container")
		return err
	}
	logger.Debug("Resumed paused container")

	return nil
}

// PauseContainer pauses a running container
func (a *HostAgent) PauseContainer(svc *service.Service, instanceID int) error {
	logger := plog.WithFields(log.Fields{
		"serviceid":    svc.ID,
		"servicename":  svc.Name,
		"pausecommand": svc.Snapshot.Pause,
		"instanceid":   instanceID,
	})
	ctrName := fmt.Sprintf("%s-%d", svc.ID, instanceID)

	// check to see if the container exists and is running
	ctr, err := docker.FindContainer(ctrName)
	if err == docker.ErrNoSuchContainer {
		// container has been deleted and the event monitor should catch this
		logger.Debug("Container not found")
		return nil
	}
	if !ctr.IsRunning() {
		// container has stopped and the event monitor should catch this
		logger.Debug("Container stopped")
		return nil
	}

	// pause the running container
	if err := attachAndRun(ctrName, svc.Snapshot.Pause); err != nil {
		logger.WithError(err).Debug("Could not pause running container")
		return err
	}
	logger.Debug("Paused running container")
	return nil
}

// pullImage pulls the service image and returns the uuid string
// of the image and the fully qualified image name.
func (a *HostAgent) pullImage(logger *log.Entry, cancel <-chan interface{}, imageID string) (string, string, error) {
	conn, err := zzk.GetLocalConnection("/")
	if err != nil {
		logger.WithError(err).Debug("Could not connect to coordinator")

		// TODO: wrap error?
		return "", "", err
	}

	timeoutC := make(chan time.Time)
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-done:
		case <-cancel:
			select {
			case <-done:
			case timeoutC <- time.Now():
			}
		}
	}()

	a.pullreg.SetConnection(conn)
	if err := a.pullreg.PullImage(timeoutC, imageID); err != nil {
		logger.WithError(err).Debug("Could not pull image")

		// TODO: wrap error?
		return "", "", err
	}
	logger.Debug("Pulled image")

	uuid, err := registry.GetImageUUID(conn, imageID)
	if err != nil {
		logger.WithError(err).Debug("Could not load image id")

		// TODO: wrap error?
		return "", "", err
	}
	logger.Debug("Found image uuid")

	name, err := a.pullreg.ImagePath(imageID)
	if err != nil {
		logger.WithError(err).Debug("Could not get full image name")

		// TODO: wrap error?
		return "", "", err
	}

	return uuid, name, nil
}

// monitorContainer tracks the running state of the container.
func (a *HostAgent) monitorContainer(logger *log.Entry, ctr *docker.Container) <-chan time.Time {
	ev := make(chan time.Time, 1)
	ctr.OnEvent(docker.Die, func(_ string) {
		defer close(ev)
		dctr, err := ctr.Inspect()
		if err != nil {
			logger.WithError(err).Error("Could not look up container")
			ev <- time.Now()
			return
		}

		logger.WithFields(log.Fields{
			"terminated": dctr.State.FinishedAt,
			"exitcode":   dctr.State.ExitCode,
		}).Debug("Container exited")

		if dctr.State.ExitCode != 0 || log.GetLevel() == log.DebugLevel {
			// TODO: need to get logs from api
			output, err := exec.Command("docker", "logs", "--tail", "10000", ctr.ID).CombinedOutput()
			if err != nil {
				logger.WithField("output", string(output)).WithError(err).Warn("Could not get container logs")
			} else {
				prefix := fmt.Sprintf("ctr-%s: ", ctr.ID[0:5])
				split := strings.Split(string(output), "\n")
				for i, s := range split {
					split[i] = prefix + s
				}
				final := strings.Join(split, "\n")
				logger.WithField("output", string(final)).Info("Last 10000 lines of container")
			}
		}

		if err := ctr.Delete(true); err != nil {
			logger.WithError(err).Warn("Could not delete container")
		}

		// just in case something unusual happened
		if !dctr.State.FinishedAt.IsZero() {
			ev <- dctr.State.FinishedAt
		} else {
			ev <- time.Now()
		}
		return
	})
	return ev
}