package views

import (
	"sync"
)

func (v *Views) OpenSettings() error {
	v.WaitGroup.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
	}(v.WaitGroup)

	return nil
}
