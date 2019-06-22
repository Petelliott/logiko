package simulator

import (
	"testing"
	"reflect"
)

func expect(t *testing.T, expected interface{}, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected '%v', got '%v'", expected, got)
	}
}

func add(in []PortType, out []PortType) {
	out[0] = in[0] + in[1]
}

func null(in []PortType, out []PortType) {}

func passthrough(in []PortType, out []PortType) {
	for i := 0; i < len(in) && i < len(out); i++ {
		out[i] = in[i]
	}
}

func TestRead(t *testing.T) {
	fc := NewFuncComponent(add, 2, 1)

	fc.Attach(0, func()PortType { return 5; })
	fc.Attach(1, func()PortType { return 7; })

	fc.Update()

	expect(t, PortType(12), fc.Read(0))
}

func TestPorts(t *testing.T) {
	fc := NewFuncComponent(null, 2, 1)

	expect(t, 1, fc.Ports())
	expect(t, 2, fc.InPorts())

	fc = NewFuncComponent(null, 0, 0)

	expect(t, 0, fc.Ports())
	expect(t, 0, fc.InPorts())
}

func TestSubscribe(t *testing.T) {
	fc := NewFuncComponent(passthrough, 1, 1)
	fc2 := NewFuncComponent(passthrough, 1, 1)
	fc3 := NewFuncComponent(passthrough, 1, 1)

	n := PortType(1)

	fc.Attach(0, func()PortType { return n; })
	fc2.Attach(0, func()PortType { return fc.Read(0); })
	fc.Subscribe(0, fc2.Update)
	fc3.Attach(0, func()PortType { return fc2.Read(0); })
	fc2.Subscribe(0, fc3.Update)

	expect(t, PortType(0), fc3.Read(0))

	fc.Update()

	expect(t, PortType(1), fc3.Read(0))

	n = PortType(2)

	fc.Update()

	expect(t, PortType(2), fc3.Read(0))
}

func TestUpdate(t *testing.T) {
	fc := NewFuncComponent(passthrough, 5, 1)
	// expect no null pointer expection
	fc.Update()
}
