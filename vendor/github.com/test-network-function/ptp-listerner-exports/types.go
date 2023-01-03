package exports

import (
	"sync"
	"time"

	ptpEvent "github.com/redhat-cne/sdk-go/pkg/event/ptp"
)

type EventType int64

const (
	LockState EventType = iota
	Port9043            = 9043
)

type LockStateValue int64

const (
	AcquiringSync LockStateValue = iota
	AntennaDisconnected
	AntennaShortCircuit
	Booting
	Freerun
	Holdover
	Locked
	Synchronized
	Unlocked
)

type StoredEvent struct {
	TimeStamp time.Time
	Source    string
	Type      EventType
	Values     []int64
}

var (
	Mu               sync.Mutex
	AllEvents        []StoredEvent
	ToEventType      = map[string]EventType{string(ptpEvent.OsClockSyncStateChange): LockState}
	ToLockStateValue = map[string]LockStateValue{
		string(ptpEvent.ACQUIRING_SYNC):        AcquiringSync,
		string(ptpEvent.ANTENNA_DISCONNECTED):  AntennaDisconnected,
		string(ptpEvent.ANTENNA_SHORT_CIRCUIT): AntennaShortCircuit,
		string(ptpEvent.BOOTING):               Booting,
		string(ptpEvent.FREERUN):               Freerun,
		string(ptpEvent.HOLDOVER):              Holdover,
		string(ptpEvent.LOCKED):                Locked,
		string(ptpEvent.SYNCHRONIZED):          Synchronized,
		string(ptpEvent.UNLOCKED):              Unlocked,
	}
)

type EventReceivedCallback func(source string, eventType EventType, eventTime time.Time, data []byte)
