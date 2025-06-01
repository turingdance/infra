package pkkit

import (
	"errors"
	"sync"
	"time"
)

const (
	workerIDBits     = 5
	datacenterIDBits = 5
	sequenceBits     = 12

	maxWorkerID     = -1 ^ (-1 << workerIDBits)
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits)
	maxSequence     = -1 ^ (-1 << sequenceBits)

	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
	sequenceMask       = maxSequence

	// 起始时间戳 (2020-01-01 00:00:00 UTC)
	epoch = int64(1577836800000)
)

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	sequence      int64
}

// NewSnowflake 创建一个雪花算法实例
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID must be between 0 and 31")
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("datacenter ID must be between 0 and 31")
	}

	return &Snowflake{
		lastTimestamp: -1,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
	}, nil
}

// NextID 生成下一个唯一 ID
func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano()/1e6 - epoch

	if timestamp < s.lastTimestamp {
		return 0, errors.New("clock moved backwards. Refusing to generate id for " +
			time.Duration(s.lastTimestamp-timestamp).String())
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 当前毫秒内序列号用完，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano()/1e6 - epoch
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return ((timestamp) << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence, nil
}
