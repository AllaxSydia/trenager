package executor

type Executor interface {
	Execute(code, language string) (map[string]interface{}, error)
}

// Контракт исполнителя, короче Абстракция
