package utils

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// func GetLoggerFromContext(c *gin.Context) *zap.Logger {
// 	logger, exists := c.Get("Logger")
// 	if !exists {
// 		return global.Logger
// 	}
// 	return logger.(*zap.Logger)
// }

func RemoveDuplicates(ids []int64) []int64 {
	unique := make(map[int64]struct{})
	var result []int64

	for _, id := range ids {
		if id != 0 {
			if _, exists := unique[id]; !exists {
				unique[id] = struct{}{}
				result = append(result, id)
			}
		}
	}
	return result
}

type WgGroup struct {
	wg     *sync.WaitGroup
	once   *sync.Once
	ctx    context.Context
	cancel context.CancelFunc
	err    error
}

func NewWgGroup() *WgGroup {
	g := WgGroup{
		wg: new(sync.WaitGroup),
	}
	ctx, cancel := context.WithCancel(context.Background())
	g.once = new(sync.Once)
	g.ctx = ctx
	g.cancel = cancel
	return &g
}
func (g *WgGroup) Wait() error {
	defer g.cancel()
	g.wg.Wait()
	return g.err
}

func (g *WgGroup) Go(f func() error) {
	g.wg.Add(1)
	go func(contextF context.Context) {
		defer g.wg.Done()
		err := f()
		if err != nil {
			g.once.Do(func() {
				g.err = err
				g.cancel()
			})
		}
		go func() {
			for {
				select {
				case <-contextF.Done():
					fmt.Println("Exiting go routine")
					return
				}
			}
		}()
	}(g.ctx)
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
