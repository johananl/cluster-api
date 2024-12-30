package ignition

import (
	"fmt"
)

// Provisioner is a provisioner that generates Ignition data.
type Provisioner struct{}

// ControlPlaneInitData returns the Ignition data for initializing a control plane node using
// kubeadm.
func (c *Provisioner) ControlPlaneInitData(data interface{}) ([]byte, error) {
	input, ok := data.(*ControlPlaneInitInput)
	if !ok {
		return nil, fmt.Errorf("invalid input data type %T, expected %T", data, &ControlPlaneInitInput{})
	}

	b, _, err := controlPlaneInitData(input)

	return b, err
}

// ControlPlaneJoinData returns the Ignition data for joining a control plane node using kubeadm.
func (c *Provisioner) ControlPlaneJoinData(data interface{}) ([]byte, error) {
	input, ok := data.(*ControlPlaneJoinInput)
	if !ok {
		return nil, fmt.Errorf("invalid input data type %T, expected %T", data, &ControlPlaneJoinInput{})
	}

	b, _, err := controlPlaneJoinData(input)

	return b, err
}

// WorkerJoinData returns the Ignition data for joining a worker node using kubeadm.
func (c *Provisioner) WorkerJoinData(data interface{}) ([]byte, error) {
	input, ok := data.(*WorkerJoinInput)
	if !ok {
		return nil, fmt.Errorf("invalid input data type %T, expected %T", data, &WorkerJoinInput{})
	}

	b, _, err := workerJoinData(input)

	return b, err
}

// NewProvisioner returns a new Provisioner for Ignition data.
func NewProvisioner() *Provisioner {
	return &Provisioner{}
}
