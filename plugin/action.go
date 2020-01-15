package plugin

// MethodSignature defines what method signature is supported for a method.
type MethodSignature int

// Defines different method signatures.
const (
	UnsupportedSignature MethodSignature = iota
	DefaultSignature
	BlockReadableSignature
)

// Action represents a Wash action.
type Action struct {
	Name                         string `json:"name"`
	Protocol                     string `json:"protocol"`
	corePluginEntrySignatureFunc func(e Entry) MethodSignature
}

var actions = make(map[string]Action)

func newAction(name string, protocol string, corePluginEntrySignatureFunc func(e Entry) MethodSignature) Action {
	a := Action{
		Name:                         name,
		Protocol:                     protocol,
		corePluginEntrySignatureFunc: corePluginEntrySignatureFunc,
	}
	actions[a.Name] = a
	return a
}

// IsSupportedOn returns true if the action's supported
// on the specified entry, false otherwise.
func (a Action) IsSupportedOn(entry Entry) bool {
	return a.signature(entry) != UnsupportedSignature
}

func (a Action) signature(entry Entry) MethodSignature {
	switch t := entry.(type) {
	case externalPlugin:
		return t.MethodSignature(a.Name)
	default:
		return a.corePluginEntrySignatureFunc(entry)
	}
}

var listAction = newAction("list", "Parent", func(e Entry) MethodSignature {
	if _, ok := e.(Parent); ok {
		return DefaultSignature
	}
	return UnsupportedSignature
})

var readAction = newAction("read", "Readable", func(e Entry) MethodSignature {
	if _, ok := e.(Readable); ok {
		return DefaultSignature
	}
	if _, ok := e.(BlockReadable); ok {
		return BlockReadableSignature
	}
	return UnsupportedSignature
})

var streamAction = newAction("stream", "Streamable", func(e Entry) MethodSignature {
	if _, ok := e.(Streamable); ok {
		return DefaultSignature
	}
	return UnsupportedSignature
})

var writeAction = newAction("write", "Writable", func(e Entry) MethodSignature {
	if _, ok := e.(Writable); ok {
		return DefaultSignature
	}
	return UnsupportedSignature
})

var execAction = newAction("exec", "Execable", func(e Entry) MethodSignature {
	if _, ok := e.(Execable); ok {
		return DefaultSignature
	}
	return UnsupportedSignature
})

var deleteAction = newAction("delete", "Deletable", func(e Entry) MethodSignature {
	if _, ok := e.(Deletable); ok {
		return DefaultSignature
	}
	return UnsupportedSignature
})

var signalAction = newAction("signal", "Signalable", func(e Entry) MethodSignature {
	if _, ok := e.(Signalable); ok {
		return DefaultSignature
	}
	return UnsupportedSignature
})

// ListAction represents the list action
func ListAction() Action {
	return listAction
}

// ReadAction represents the read action
func ReadAction() Action {
	return readAction
}

// StreamAction represents the stream action
func StreamAction() Action {
	return streamAction
}

// WriteAction represents the append action
func WriteAction() Action {
	return writeAction
}

// ExecAction represents the exec action
func ExecAction() Action {
	return execAction
}

// DeleteAction represents the delete action
func DeleteAction() Action {
	return deleteAction
}

// SignalAction represents the signal action
func SignalAction() Action {
	return signalAction
}

// Actions returns all of the available Wash actions as a map
// of <action_name> => <action_object>.
func Actions() map[string]Action {
	// We create a clone of the actions map so that callers won't
	// be able to modify it.
	mp := make(map[string]Action)
	for k, v := range actions {
		mp[k] = v
	}
	return mp
}

// SupportedActionsOf returns all of the given
// entry's supported actions.
func SupportedActionsOf(entry Entry) []string {
	var supportedActions []string

	// Iterating over the actions and calling action.signature
	// results in multiple type switches, which can unnecessarily
	// slow down list. This implementation does the same thing
	// but in a single type-switch.
	switch t := entry.(type) {
	case externalPlugin:
		for _, action := range actions {
			signature := t.MethodSignature(action.Name)
			if signature != UnsupportedSignature {
				supportedActions = append(supportedActions, action.Name)
			}
		}
	default:
		for _, action := range actions {
			signature := action.corePluginEntrySignatureFunc(entry)
			if signature != UnsupportedSignature {
				supportedActions = append(supportedActions, action.Name)
			}
		}
	}

	return supportedActions
}
