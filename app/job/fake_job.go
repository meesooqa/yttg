package job

// FakeJob — тестовая реализация интерфейса Job.
type FakeJob struct {
	id      string
	execErr error // если не nil, Execute возвращает ошибку
}

func (fj FakeJob) Execute() error {
	return fj.execErr
}

func (fj FakeJob) GetID() string {
	return fj.id
}

// GetStatus не используется в воркере, так как статус записывается в очередь
func (fj FakeJob) GetStatus() JobStatus {
	return ""
}
