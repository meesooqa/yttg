package job

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestJobQueue_AddJobAndGetJobsStatuses(t *testing.T) {
	queue := NewJobQueue()

	job1 := FakeJob{id: "job1"}
	job2 := FakeJob{id: "job2"}

	queue.AddJob(job1)
	queue.AddJob(job2)

	statuses := queue.GetJobsStatuses()
	if statuses["job1"] != StatusQueued {
		t.Errorf("Ожидался статус 'queued' для job1, получено %s", statuses["job1"])
	}
	if statuses["job2"] != StatusQueued {
		t.Errorf("Ожидался статус 'queued' для job2, получено %s", statuses["job2"])
	}
}

func TestWorker_ProcessJobs(t *testing.T) {
	queue := NewJobQueue()

	// Создаем две задачи: одну успешную, вторую с ошибкой
	jobSuccess := FakeJob{id: "success", execErr: nil}
	jobFail := FakeJob{id: "fail", execErr: errors.New("execution error")}

	// Добавляем задачи в очередь
	queue.AddJob(jobSuccess)
	queue.AddJob(jobFail)

	// Запускаем воркера в отдельной горутине
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// После обработки всех задач канал нужно закрыть, чтобы выйти из цикла
		Worker(1, queue)
		wg.Done()
	}()

	// Ждем, пока задачи попадут в очередь
	time.Sleep(100 * time.Millisecond)
	// Закрываем очередь, чтобы завершить воркера
	close(queue.queue)
	// Ожидаем завершения воркера
	wg.Wait()

	statuses := queue.GetJobsStatuses()
	if statuses["success"] != StatusDone {
		t.Errorf("Ожидался статус 'done' для задачи success, получено %s", statuses["success"])
	}
	if statuses["fail"] != StatusFailed {
		t.Errorf("Ожидался статус 'failed' для задачи fail, получено %s", statuses["fail"])
	}
}
