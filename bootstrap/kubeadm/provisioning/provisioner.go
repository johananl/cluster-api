package provisioning

// Provisioner represents a type which generates data necessary to provision a machine. The
// generated data, represented as a byte slice, is typically passed to the machine by an
// infrastructure provider using mechanisms such as an instance metadata service.
//
// The input argument in all methods is an interface{} representing an arbitrary value required for
// rendering the provisioning data to be returned. It's up to the implementation to define the
// structure of the data and how it's used. Typically, the input value is a data structure
// containing multiple values which get interpolated into a text template by the implementation.
//
// Implementations should assert the data to the expected concrete type and return an error if the
// assertion fails.
type Provisioner interface {
	// ControlPlaneInitData returns the data to be used for provisioning a control plane node
	// during the kubeadm init phase.
	ControlPlaneInitData(input interface{}) ([]byte, error)
	// ControlPlaneJoinData returns the data to be used for provisioning a control plane node
	// during the kubeadm join phase.
	ControlPlaneJoinData(input interface{}) ([]byte, error)
	// WorkerJoinData returns the data to be used for provisioning a worker node (worker nodes only
	// have a join phase).
	WorkerJoinData(input interface{}) ([]byte, error)
}
