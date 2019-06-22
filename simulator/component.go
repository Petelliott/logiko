package simulator

type Component interface {
	// Read gets the value at a specific output port. The value returned will
	// be of the largest bus type with higher bits guarenteed to be zero
	Read(port int) PortType

	// Ports returns the number of ports
	Ports() int

	// Subscribe attaches a function to be called when the value of port
	// changes
	Subscribe(port int, fun func())

	// Update recalculates the outputs of the Component
	Update()
}

// AttachableComponents are components whose inputs can be attached to outputs
type AttachableComponent interface {
	Component

	// InPorts gets the number of input ports
	InPorts() int

	// Attach inport on this commponent to another component's Read
	Attach(inport int, fun func()PortType)
}




