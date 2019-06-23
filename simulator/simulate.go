package simulator

// Sim is a Component that wraps an AttachableComponent for easy simulating
type Sim struct {
	comp   AttachableComponent
	inputs []PortType
}

// NewSim creates a simulator with all zero input, and updated
func NewSim(comp AttachableComponent) *Sim {
	sim := &Sim{
		comp: comp,
		inputs: make([]PortType, comp.InPorts()),
	}

	for idx, _ := range sim.inputs {
		// closure over value, not variable
		i := idx
		comp.Attach(idx, func() PortType {
			return sim.inputs[i]
		})
	}

	sim.Update()
	return sim
}

// Write sets the value of a specific input port
func (s *Sim) Write(port int, value PortType) {
	s.inputs[port] = value
	s.Update()
}

func (s *Sim) Read(port int) PortType {
	return s.comp.Read(port)
}

func (s *Sim) Ports() int {
	return s.comp.Ports()
}

func (s *Sim) Subscribe(port int, fun func()) {
	s.comp.Subscribe(port, fun)
}

func (s *Sim) Update() {
	s.comp.Update()
}
