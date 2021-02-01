package hid

import (
	"errors"
	"net"
	"sync"
	"time"
)

// 包hid实现了SonyFlake，这是一个受Twitter的Snowflake启发的分布式唯一ID生成器。
//
// SonFlake ID由以下组成
//     39位，以10毫秒为单位的时间
//      8位为序列号
//     16位机器ID
//
// +-----------------------------------------------------------------------------+
// | 1 Bit Unused | 39 Bit Timestamp |  8 Bit Sequence ID  |   16 Bit Machine ID |
// +-----------------------------------------------------------------------------+
// 39bit 时间戳，SonyFlake是以10毫秒为单位保存时间，使用年限为174年。
//  8bit 序列号，每10毫秒最大生成256个，1秒最多生成25600个。
// 16bit 机器号，默认的是当前机器的私有IP的最后两位。
//

// 这些常数是SonyFlake ID部分的位长度。
const (
	BitLenTime      = 39                               // 时间位长
	BitLenSequence  = 8                                // 序列号的位长
	BitLenMachineID = 63 - BitLenTime - BitLenSequence // 机器ID的位长
)

// Settings 配置SonyFlake:
//
//  StartTime 是将SonyFlake时间定义为经过时间的时间。
//  如果StartTime为0，则将SonyFlake的开始时间设置为"2014-09-01 00:00:00 +0000 UTC"。
//  如果StartTime早于当前时间，则不会创建SonyFlake。
//
//  MachineID 返回SonyFlake实例的唯一ID。
//  如果MachineID返回错误，则不会创建SonyFlake。
//  如果MachineID为nil，则使用默认的MachineID。
//  默认的MachineID返回私有IP地址的低16位。
//
//  CheckMachineID 验证机器ID的唯一性。
//  如果CheckMachineID返回false，说明SonyFlake未创建。
//  如果CheckMachineID为nil，则不进行验证。
type Settings struct {
	StartTime      time.Time
	MachineID      func() (uint16, error)
	CheckMachineID func(uint16) bool
}

// SonyFlake 是一个分布式唯一ID生成器。
type SonyFlake struct {
	mutex       *sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
	machineID   uint16
}

// NewSonyFlake 返回用给定设置配置的新SonyFlake。
//  NewSonyFlake 在以下情况下返回nil：
//  - Settings.StartTime 在当前时间之前。
//  - Settings.MachineID 返回一个错误。
//  - Settings.CheckMachineID 返回false。
func NewSonyFlake(st Settings) *SonyFlake {
	sf := new(SonyFlake)
	sf.mutex = new(sync.Mutex)
	sf.sequence = uint16(1<<BitLenSequence - 1)

	if st.StartTime.After(time.Now()) {
		return nil
	}
	if st.StartTime.IsZero() {
		sf.startTime = toSonyFlakeTime(time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC))
	} else {
		sf.startTime = toSonyFlakeTime(st.StartTime)
	}

	var err error
	if st.MachineID == nil {
		sf.machineID, err = lower16BitPrivateIP()
	} else {
		sf.machineID, err = st.MachineID()
	}
	if err != nil || (st.CheckMachineID != nil && !st.CheckMachineID(sf.machineID)) {
		return nil
	}

	return sf
}

// NextID 生成下一个唯一ID。
//  SonyFlake时间溢出后，NextID返回一个错误。
func (sf *SonyFlake) NextID() (uint64, error) {
	const maskSequence = uint16(1<<BitLenSequence - 1)

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	current := currentElapsedTime(sf.startTime)
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else { // sf.elapsedTime >= current
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			time.Sleep(sleepTime(overtime))
		}
	}

	return sf.toID()
}

const sonyFlakeTimeUnit = 1e7 // 纳秒，换算10毫秒

func toSonyFlakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / sonyFlakeTimeUnit
}

func currentElapsedTime(startTime int64) int64 {
	return toSonyFlakeTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%sonyFlakeTimeUnit)*time.Nanosecond
}

func (sf *SonyFlake) toID() (uint64, error) {
	if sf.elapsedTime >= 1<<BitLenTime {
		return 0, errors.New("超过时间限制")
	}

	return uint64(sf.elapsedTime)<<(BitLenSequence+BitLenMachineID) |
		uint64(sf.sequence)<<BitLenMachineID |
		uint64(sf.machineID), nil
}

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipNet, ok := a.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}

		ip := ipNet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("没有私有ip地址")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

// Decompose 返回一组SonyFlake ID部分。
func Decompose(id uint64) map[string]uint64 {
	const maskSequence = uint64((1<<BitLenSequence - 1) << BitLenMachineID)
	const maskMachineID = uint64(1<<BitLenMachineID - 1)

	msb := id >> 63
	sTime := id >> (BitLenSequence + BitLenMachineID)
	sequence := id & maskSequence >> BitLenMachineID
	machineID := id & maskMachineID
	return map[string]uint64{
		"id":         id,
		"msb":        msb,
		"time":       sTime,
		"sequence":   sequence,
		"machine-id": machineID,
	}
}
