package hid

import (
	"github.com/gogf/gf/text/gstr"

	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	Epoch              = 1288834974657                                  // 毫秒，Thu, 04 Nov 2010 01:42:54 GMT
	DataCenterIdBits   = 5                                              // 数据标识id占位
	WorkerIdBits       = 5                                              // 机器id占位
	MaxDataCenterId    = -1 ^ (-1 << DataCenterIdBits)                  // 最大数据标识id
	MaxWorkerId        = -1 ^ (-1 << WorkerIdBits)                      // 最大机器id
	SequenceBits       = 12                                             // 序列占位
	TimestampLeftShift = SequenceBits + WorkerIdBits + DataCenterIdBits // 毫秒，左移22位
	DataCenterIdShift  = SequenceBits + WorkerIdBits                    // 数据标识id，左移17位
	WorkerIdShift      = SequenceBits                                   // 机器id，左移12位
	SequenceMask       = -1 ^ (-1 << SequenceBits)                      // 序列掩码
)

// SnowFlake 共64位
// 0       - 0000000000 0000000000 0000000000 0000000000 0 - 00000         - 00000     - 000000000000
// 正号|1位 - 毫秒|41位                                      - 数据标识id|5位 - 机器id|5位 - 序列|12位
type SnowFlake struct {
	mu            sync.Mutex // protects following fields
	lastTimestamp int64      // 上次毫秒
	dataCenterId  int64      // 数据标识id
	workerId      int64      // 机器id
	sequence      int64      // 序列
}

func NewSnowFlake(dataCenterId int64, workerId int64) (*SnowFlake, error) {
	if dataCenterId > MaxDataCenterId || dataCenterId < 0 {
		return nil, fmt.Errorf("data center id can't be greater than %d or less than 0", MaxDataCenterId)
	}

	if workerId > MaxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("worker id can't be greater than %d or less than 0", MaxWorkerId)
	}

	return &SnowFlake{
		lastTimestamp: -1,
		dataCenterId:  dataCenterId,
		workerId:      workerId,
		sequence:      0,
	}, nil
}

func (x *SnowFlake) NextId() (int64, error) {
	x.mu.Lock()
	defer x.mu.Unlock()

	timestamp := x.TimeGen()
	if timestamp < x.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards. refusing to generate id for %d milliseconds", x.lastTimestamp-timestamp)
	}

	if x.lastTimestamp == timestamp {
		x.sequence = (x.sequence + 1) & SequenceMask
		if x.sequence == 0 {
			timestamp = x.TilNextMillis(x.lastTimestamp)
		}
	} else {
		x.sequence = 0
	}

	x.lastTimestamp = timestamp

	return ((timestamp - Epoch) << TimestampLeftShift) |
		(x.dataCenterId << DataCenterIdShift) |
		(x.workerId << WorkerIdShift) |
		x.sequence, nil
}

func (x *SnowFlake) TilNextMillis(lastTimestamp int64) int64 {
	timestamp := x.TimeGen()
	for ; timestamp <= lastTimestamp; timestamp = x.TimeGen() {
		runtime.Gosched()
	}
	return timestamp
}

func (x *SnowFlake) TimeGen() int64 {
	return time.Now().UnixNano() / 1e6
}

// IdToTimestamp 雪花算法ID转时间戳，毫秒
func (x *SnowFlake) IdToTimestamp(snowflakeId int64) int64 {
	bid := fmt.Sprintf("%b", big.NewInt(snowflakeId))
	timestamp := len(bid) - (WorkerIdBits + DataCenterIdBits + SequenceBits)
	substring := gstr.SubStr(bid, 0, timestamp)
	parseInt, _ := strconv.ParseInt(substring, 2, 64)
	return parseInt + Epoch
}

// NextIdStr 获取雪花算法ID，字符串
func (x *SnowFlake) NextIdStr() string {
	nextId, _ := x.NextId()
	return strconv.FormatInt(nextId, 10)
}

// NextIdInt64 获取雪花算法ID，int64
func (x *SnowFlake) NextIdInt64() int64 {
	nextId, _ := x.NextId()
	return nextId
}
