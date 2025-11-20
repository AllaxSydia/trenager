package executor

// Executor интерфейс для выполнения кода
type Executor interface {
	Execute(code, language string, inputs []string) (map[string]interface{}, error)
}

// Cleaner интерфейс для очистки ресурсов
type Cleaner interface {
	Cleanup()
}
