package timewheel_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/turingdance/infra/timewheel"
)

type teststruct struct {
	index int
}

func TestMain(t *testing.T) {
	monitor := timewheel.New(20, time.Millisecond*500)
	for i := 0; i < 140; i++ {
		monitor.Watch(teststruct{
			index: i,
		}, time.Second*time.Duration(i), func(ele any) (remove bool) {
			stru := ele.(teststruct)
			return stru.index%3 == 0
		})
	}
	monitor.Start()
	fmt.Scanln()
}
