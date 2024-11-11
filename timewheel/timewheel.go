package timewheel

// 时间轮
import (
	"container/list"
	"math"
	"sync"
	"time"
)

type Timewheel struct {
	unit         time.Duration
	numofslot    int
	currentIndex int //当前游标
	chain        []*list.List
	rw           *sync.Mutex
	ticker       *time.Ticker
	stop         chan bool
}
type OnTicker func(data any) bool
type Job struct {
	data any
	fun  OnTicker
}

// 检测周期
// 插槽数
func New(numofslot int, unit time.Duration) *Timewheel {
	chain := make([]*list.List, 0)
	for i := 0; i < numofslot; i++ {
		chain = append(chain, list.New())
	}
	return &Timewheel{
		unit:         unit,
		numofslot:    numofslot,
		currentIndex: 0,
		chain:        chain,
		rw:           &sync.Mutex{},
		ticker:       time.NewTicker(unit),
	}
}
func NewJob(data any, on OnTicker) Job {
	return Job{
		data: data,
		fun:  on,
	}
}

// tm 时间后再次检查,完了后执行fun函数
func (m *Timewheel) Add(job Job, tm time.Duration) {
	//
	m.rw.Lock()
	num := int(math.Ceil(float64(tm) / float64(m.unit)))
	index := (m.currentIndex + num) % m.numofslot
	m.chain[index].PushBack(job)
	m.rw.Unlock()
}

// tm 时间后再次检查,完了后执行fun函数
//
// 如果remove是true 则执行后删除该元素,否则不删除
//
//	func(ele any)remove bool{
//		if(ele.(int)%3==0){
//			remove = true
//		else{
//			remove = false
//		}
//	}
//
// }
func (m *Timewheel) Watch(ptr any, tm time.Duration, oncheckfunc func(ele any) (removeornot bool)) {
	task := Job{
		data: ptr,
		fun:  oncheckfunc,
	}
	//fmt.Println("add prt", ptr)
	m.Add(task, tm)
}
func (m *Timewheel) Cancel(ptr any) {
	m.rw.Lock()
	for i := 0; i < m.currentIndex; i++ {
		chain := m.chain[i]
		toRemove := []*list.Element{}
		for e := chain.Front(); e != nil; e = e.Next() {
			job := e.Value.(Job)
			if job.data == ptr {
				toRemove = append(toRemove, e)
			}
		}
		for _, e := range toRemove {
			chain.Remove(e)
		}
	}
	m.rw.Unlock()
}

// tm 时间后再次检查,完了后执行fun函数
func (m *Timewheel) run() {
	for {
		select {
		case <-m.ticker.C:
			m.rw.Lock()
			chain := m.chain[m.currentIndex]
			toRemove := []*list.Element{}
			for e := chain.Front(); e != nil; e = e.Next() {
				job := e.Value.(Job)
				if job.fun(job.data) {
					toRemove = append(toRemove, e)
				}
			}
			for _, e := range toRemove {
				chain.Remove(e)
			}
			//fmt.Printf("[%d] DELETE  [%d]\n", m.currentIndex, len(toRemove))
			m.currentIndex = (m.currentIndex + 1) % m.numofslot
			m.rw.Unlock()
		case <-m.stop:
			return
		}
	}
}

// tm 时间后再次检查,完了后执行fun函数
func (m *Timewheel) Start() {
	go m.run()
}

// tm 时间后再次检查,完了后执行fun函数
func (m *Timewheel) Shutdown() {
	m.stop <- true
}
