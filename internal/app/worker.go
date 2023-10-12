package app

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
Создаю пул воркеров (получается 1 соединение к базе)
В воркер нужно передать storage
В методе обработчике у стораджа вызываю удаление урла, в параметре передаю идентификатор ссылки и удаляю её в go рутине
*/
type worker struct {
	wg         *sync.WaitGroup
	cancelFunc context.CancelFunc
	storage    *Storage
}

type Worker interface {
	Start(pctx context.Context)
	Stop()
}

func NewWorker(storage *Storage) Worker {
	w := worker{
		wg:      new(sync.WaitGroup),
		storage: storage,
	}
	return &w
}
func (w *worker) Start(pctx context.Context) {
	ctx, cancelFunc := context.WithCancel(pctx)
	w.cancelFunc = cancelFunc
	w.wg.Add(1)
	go w.spawnWorkers(ctx)
}
func (w *worker) Stop() {
	w.cancelFunc()
	w.wg.Wait()
}

func (w *worker) spawnWorkers(ctx context.Context) {
	defer w.wg.Done()
	t := time.NewTicker(10000 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			w.doWork(ctx)
		}
	}
}

func (w *worker) doWork(ctx context.Context) {
	//	rnd:=rand.int63()
	//w.storage.
	//здесь должна быть реализация удаления урлов в кипере
	fmt.Println("Сработал метод do work")
}
