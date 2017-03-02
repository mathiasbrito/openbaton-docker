package server

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
)

const (
	// monitoringDelay represents the default time between two monitoring checks.
	monitoringDelay = 30 * time.Second
)

// fetchDockerContainers fetches the ID and state of all the available Docker containers.
func (svc *service) fetchDockerContainers() (map[string]pop.Container_Status, error) {
	conts, err := svc.cln.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	m := make(map[string]pop.Container_Status, len(conts))

	for _, cont := range conts {
		m[cont.ID] = matchState(cont.State)
	}

	return m, nil
}

// refreshLoop implements a best-effort service that periodically monitors
// and refreshes the status of the containers based on their backing Docker containers.
func (svc *service) refreshLoop() {
	tag := util.FuncName()

	svc.WithFields(log.Fields{
		"tag": tag,
	}).Debug("refresh loop spawned")

	for {
		select {
		case <-svc.quitChan:
			svc.WithFields(log.Fields{
				"tag": tag,
			}).Debug("refresh loop stopping")

			close(svc.quitChan)
			return

		case <-time.After(monitoringDelay):
			svc.refreshStatuses()
		}
	}
}

func (svc *service) refreshStatuses() error {
	tag := util.FuncName()

	svc.WithFields(log.Fields{
		"tag": tag,
	}).Debug("refreshing containers")

	// get the lock
	svc.contsMux.Lock()
	defer svc.contsMux.Unlock()

	// fetch Docker containers and their states
	statuses, err := svc.fetchDockerContainers()
	if err != nil {
		svc.WithError(err).WithFields(log.Fields{
			"tag": tag,
		}).Error("error while updating containers")
		return err
	}

	// match Docker states with Pop container states
	svc.updateStatuses(statuses)

	return nil
}

// updateStatuses is executed under the container list lock, and updates the state
// of a container, matching with its Docker correspective container
func (svc *service) updateStatuses(states map[string]pop.Container_Status) {
	tag := util.FuncName()

	for _, cont := range svc.conts {
		// only running Pop containers have a Docker container
		if cont.DockerID != "" && cont.Status == pop.Container_RUNNING {
			svc.WithFields(log.Fields{
				"tag":            tag,
				"container-name": cont.Names[0],
			}).Debug("updating container")

			state, found := states[cont.DockerID]
			// The Docker container may have been shut down by any reason. In this case,
			// mark the Pop container as FAILED.
			if !found {
				cont.Status = pop.Container_FAILED
				cont.ExtendedStatus = "the Docker container terminated unexpectedly"

				continue
			}

			stopped := false

			switch state {
			case pop.Container_EXITED:
				cont.ExtendedStatus = "the container cleanly exited"
				cont.DockerID = ""
				stopped = true

			case pop.Container_FAILED:
				cont.ExtendedStatus = "the Docker container terminated unexpectedly"
				cont.DockerID = ""
				stopped = true
			}

			// if the container stopped, reclaim its IPs
			if stopped {
				if err := svc.releaseContIPs(cont); err != nil {
					svc.WithError(err).WithFields(log.Fields{
						"tag":             tag,
						"container-names": cont.Names,
					}).Error("couldn't reclaim container IPs")
				}
			}

			cont.Status = state
		}
	}
}

func matchState(dockerState string) pop.Container_Status {
	switch strings.ToLower(dockerState) {
	case "created":
		return pop.Container_CREATED

	case "running":
		return pop.Container_RUNNING

	case "exited":
		return pop.Container_EXITED

	case "dead":
		return pop.Container_FAILED

	case "stopping":
		return pop.Container_STOPPING

	default:
		return pop.Container_UNAVAILABLE
	}
}
