package models

func AddToKeyMap(id, key string) {
	if Keys == nil {
		Keys = make(map[string]string)
	}

	_, exists := Keys[id]
	if !exists {
		Keys[id] = key
	}
}

func RemoveFromKeyMap(id string) {
	_, found := Keys[id]
	if found {
		delete(Keys, id)
	}
}

func AddToPortMap(port int, status bool) {
	if PortStatus == nil {
		PortStatus = make(map[int]bool)
	}

	PortStatus[port] = status
}

func RemoveFromPortMap(port int) {
	_, found := PortStatus[port]
	if found {
		delete(PortStatus, port)
	}
}

func AddToClientMap(id, key string) {
	if Keys == nil {
		Keys = make(map[string]string)
	}

	_, exists := Keys[id]
	if !exists {
		Keys[id] = key
	}
}

func RemoveFromClientMap(id string) {
	_, found := Keys[id]
	if found {
		delete(Keys, id)
	}
}
