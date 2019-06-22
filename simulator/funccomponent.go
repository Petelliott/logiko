package simulator

// FuncComponent Wraps a function and is an AttachableComponent
type FuncComponent struct {
	fun func([]PortType, []PortType)
	in  []func()PortType
	out []PortType
	subs [][]func()
}

func NewFuncComponent(fun func([]PortType, []PortType), nin int, nout int) *FuncComponent {
	fc := &FuncComponent{
		fun:  fun,
		in:   make([]func()PortType, nin),
		out:  make([]PortType, nout),
		subs: make([][]func(), nout),
	}
	fc.Update()
	return fc
}

func (fc *FuncComponent) Read(port int) PortType {
	return fc.out[port]
}

func (fc *FuncComponent) Ports() int {
	return len(fc.out)
}

func (fc *FuncComponent) Subscribe(port int, fun func()) {
	fc.subs[port] = append(fc.subs[port], fun)
	fun()
}

func (fc *FuncComponent) Update() {
	oldout := make([]PortType, len(fc.out))
	copy(oldout, fc.out)

	newin := make([]PortType, len(fc.in))
	for port, fun := range fc.in {
		newin[port] = fun()
	}

	fc.fun(newin, fc.out)

	for port, subs := range fc.subs {
		if oldout[port] != fc.out[port] {
			for _, sub := range subs {
				sub()
			}
		}
	}
}

func (fc FuncComponent) InPorts() int {
	return len(fc.in)
}

func (fc *FuncComponent) Attach(inport int, fun func()PortType) {
	fc.in[inport] = fun
}
