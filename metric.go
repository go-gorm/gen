package gen

//go:generate stringer -type=BasicMethod -trimprefix method

type BasicMethod uint

const (
	methodDiy BasicMethod = iota

	methodSelect
	methodWhere
	methodOrder
	methodDistinct
	methodOmit
	methodNot
	methodOr
	methodJoin
	methodGroup
	methodHaving
	methodLimit
	methodOffset
	methodScopes
	methodPreload
	methodAttrs
	methodAssign
	methodUnscoped

	methodCreate
	methodCreateInBatches
	methodSave
	methodFirst
	methodTake
	methodLast
	methodFind
	methodFindInBatches
	methodFirstOrInit
	methodFirstOrCreate
	methodUpdate
	methodUpdates
	methodUpdateColumn
	methodUpdateColumns
	methodDelete
	methodCount
	methodRow
	methodRows
	methodScan
	methodPluck
	methodScanRows
	methodTransaction
	methodBegin
	methodCommit
	methodRollback
	methodSavePoint
	methodRollbackTo
)

type metricsClient interface {
	Emit(BasicMethod)
}

var mClient metricsClient

// UseMetric specify a metric client
func UseMetric(client metricsClient) {
	mClient = client
}

// Emit metrics emit
func Emit(m BasicMethod) {
	if mClient != nil {
		mClient.Emit(m)
	}
}
