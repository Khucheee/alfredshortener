package app

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

type worker struct {
	wg         *sync.WaitGroup
	cancelFunc context.CancelFunc
	storage    *Storage
	toProcess  chan ForDelete
}

type Worker interface {
	Start(pctx context.Context)
	Stop()
}

func NewWorker(storage *Storage) (Worker, chan ForDelete) {
	w := worker{
		wg:        new(sync.WaitGroup),
		storage:   storage,
		toProcess: make(chan ForDelete),
	}
	return &w, w.toProcess
}

func (w *worker) Start(pctx context.Context) {
	ctx, cancelFunc := context.WithCancel(pctx)
	w.cancelFunc = cancelFunc

	for i := 0; i <= runtime.NumCPU(); i++ {
		w.wg.Add(1)
		go w.spawnWorkers(ctx)
	}

}
func (w *worker) Stop() {
	w.cancelFunc()
	w.wg.Wait()
}

func (w *worker) spawnWorkers(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case value := <-w.toProcess:
			fmt.Println("Произошло чтение из канала, получены значения:", value)
			w.storage.DeleteUserLink(value.uid, value.hash)
		}
	}
}

/*func (w *worker) doWork(ctx context.Context) {
	fmt.Println("Сработал метод do work, это значит, что воркер отрабатывает каждые n секунд")
}
*/
