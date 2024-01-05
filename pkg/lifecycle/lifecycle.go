package lifecycle

import (
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/containrrr/watchtower/pkg/types"
	log "github.com/sirupsen/logrus"
)

type ExecCommandFunc func(client container.Client, container types.Container)

// ExecutePreCheckCommand tries to run the pre-check lifecycle hook for a single container.
func ExecutePreCheckCommand(client container.Client, container types.Container) {
	err := ExecuteLifeCyclePhaseCommand(types.PreCheck, client, container)
	if err != nil {
		log.WithField("container", container.Name()).Error(err)
	}
}

// ExecutePostCheckCommand tries to run the post-check lifecycle hook for a single container.
func ExecutePostCheckCommand(client container.Client, container types.Container) {
	err := ExecuteLifeCyclePhaseCommand(types.PostCheck, client, container)
	if err != nil {
		log.WithField("container", container.Name()).Error(err)
	}
}

// ExecutePreUpdateCommand tries to run the pre-update lifecycle hook for a single container.
func ExecutePreUpdateCommand(client container.Client, container types.Container) error {
	return ExecuteLifeCyclePhaseCommand(types.PreUpdate, client, container)
}

// ExecutePostUpdateCommand tries to run the post-update lifecycle hook for a single container.
func ExecutePostUpdateCommand(client container.Client, newContainerID types.ContainerID) {
	newContainer, err := client.GetContainer(newContainerID)
	if err != nil {
		log.WithField("containerID", newContainerID.ShortID()).Error(err)
		return
	}

	err = ExecuteLifeCyclePhaseCommand(types.PostUpdate, client, newContainer)
	if err != nil {
		log.WithField("container", newContainer.Name()).Error(err)
	}
}

// ExecuteLifeCyclePhaseCommand tries to run the corresponding lifecycle hook for a single container.
func ExecuteLifeCyclePhaseCommand(phase types.LifecyclePhase, client container.Client, container types.Container) error {

	timeout := container.GetLifecycleTimeout(phase)
	command := container.GetLifecycleCommand(phase)
	clog := log.WithField("container", container.Name())

	if len(command) == 0 {
		clog.Debugf("No %v command supplied. Skipping", phase)
		return nil
	}

	if !container.IsRunning() || container.IsRestarting() {
		clog.Debugf("Container is not running. Skipping %v command.", phase)
		return nil
	}

	clog.Debugf("Executing %v command.", phase)
	return client.ExecuteCommand(container.ID(), command, timeout)
}
