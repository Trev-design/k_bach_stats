package credentialdistro

func (handler *grpcNewCredsStreamHandler) register(id string) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	handler.scheduled[id] = false
}

func (handler *grpcNewCredsStreamHandler) confirm(id string) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	handler.scheduled[id] = true
}

func (handler *grpcNewCredsStreamHandler) isScheduled() bool {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	for _, state := range handler.scheduled {
		if !state {
			return false
		}
	}

	return true
}

func (handler *grpcNewCredsStreamHandler) resetState() {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	for id := range handler.scheduled {
		handler.scheduled[id] = false
	}
}
