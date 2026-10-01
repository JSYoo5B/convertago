package conversion

// State tracks active pointer and slice traversals in either source reader.
// A fresh State belongs to one conversion; shared siblings are not cycles.
type State struct {
	active map[sourceVisit]bool
}

type sourceVisit struct {
	pointer any
	kind    string
}

func (state *State) Enter(pointer any, kind string) bool {
	if state.active == nil {
		state.active = make(map[sourceVisit]bool)
	}
	key := sourceVisit{pointer, kind}
	if state.active[key] {
		return false
	}
	state.active[key] = true
	return true
}

func (state *State) Leave(pointer any, kind string) {
	delete(state.active, sourceVisit{pointer, kind})
}
