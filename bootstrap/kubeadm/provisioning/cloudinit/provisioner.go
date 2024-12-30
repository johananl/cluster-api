package cloudinit

import "fmt"

// Provisioner is a provisioner that generates cloud-init data.
type Provisioner struct{}

// ControlPlaneInitData returns the cloud-init data for initializing a control plane node using
// kubeadm.
func (c *Provisioner) ControlPlaneInitData(data interface{}) ([]byte, error) {
	input, ok := data.(*ControlPlaneInitInput)
	if !ok {
		return nil, fmt.Errorf("invalid input data type %T, expected %T", data, &ControlPlaneInitInput{})
	}

	return controlPlaneInitData(input)
}

// ControlPlaneJoinData returns the cloud-init data for joining a control plane node using kubeadm.
func (c *Provisioner) ControlPlaneJoinData(data interface{}) ([]byte, error) {
	input, ok := data.(*ControlPlaneJoinInput)
	if !ok {
		return nil, fmt.Errorf("invalid input data type %T, expected %T", data, &ControlPlaneJoinInput{})
	}

	return controlPlaneJoinData(input)
}

// WorkerJoinData returns the cloud-init data for joining a worker node using kubeadm.
func (c *Provisioner) WorkerJoinData(data interface{}) ([]byte, error) {
	input, ok := data.(*WorkerJoinInput)
	if !ok {
		return nil, fmt.Errorf("invalid input data type %T, expected %T", data, &WorkerJoinInput{})
	}

	return workerJoinData(input)
}

// NewProvisioner returns a new Provisioner for cloud-init data.
func NewProvisioner() *Provisioner {
	return &Provisioner{}
}
