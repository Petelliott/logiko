package simulator

import (
	"testing"
)

func TestSim(t *testing.T) {
	fc := NewFuncComponent(add, 2, 1)
	sim := NewSim(fc)

	sim.Write(0, 8)
	sim.Write(1, 13)
	expect(t, PortType(21), sim.Read(0))

	expect(t, 1, sim.Ports())

	hit := false

	sim.Subscribe(0, func() {
		hit = true
	})

	sim.Update()
	expect(t, true, hit)
}
