package gen

const (
	diyM BasicMethod = iota // 用户自定义方法

	// modelM 不支持 Model 方法
	selectM
	whereM
	orderM
	distinctM
	omitM
	notM
	orM
	joinsM
	groupM
	havingM
	limitM
	offsetM
	scopesM
	preloadM
	attrsM
	assignM
	unscopedM

	createM
	createInBatchesM
	saveM
	firstM
	takeM
	lastM
	findM
	findInBatchesM
	firstOrInitM
	firstOrCreateM
	updateM
	updatesM
	updateColumnM
	updateColumnsM
	deleteM
	countM
	rowM
	rowsM
	scanM
	pluckM
	scanRowsM
	transactionM
	beginM
	commitM
	rollbackM
	savePointM
	rollbackToM
)

type BasicMethod uint

func (b BasicMethod) String() string {
	switch b {
	case diyM:
		return "diy"
	case selectM:
		return "select"
	case whereM:
		return "where"
	case orderM:
		return "order"
	case distinctM:
		return "distinct"
	case omitM:
		return "omit"
	case notM:
		return "not"
	case orM:
		return "or"
	case joinsM:
		return "joins"
	case groupM:
		return "group"
	case havingM:
		return "having"
	case limitM:
		return "limit"
	case offsetM:
		return "offset"
	case scopesM:
		return "scopes"
	case preloadM:
		return "preload"
	case attrsM:
		return "attrs"
	case assignM:
		return "assign"
	case unscopedM:
		return "unscoped"
	case createM:
		return "create"
	case createInBatchesM:
		return "create_in_batches"
	case saveM:
		return "save"
	case firstM:
		return "first"
	case takeM:
		return "take"
	case lastM:
		return "last"
	case findM:
		return "find"
	case findInBatchesM:
		return "find_in_batches"
	case firstOrInitM:
		return "first_or_init"
	case firstOrCreateM:
		return "first_or_create"
	case updateM:
		return "update"
	case updatesM:
		return "updates"
	case updateColumnM:
		return "update_column"
	case updateColumnsM:
		return "update_columns"
	case deleteM:
		return "delete"
	case countM:
		return "count"
	case rowM:
		return "row"
	case rowsM:
		return "rows"
	case scanM:
		return "scan"
	case pluckM:
		return "pluck"
	case scanRowsM:
		return "scan_rows"
	case transactionM:
		return "transaction"
	case beginM:
		return "begin"
	case commitM:
		return "commit"
	case rollbackM:
		return "rollback"
	case savePointM:
		return "savepoint"
	case rollbackToM:
		return "rollbackto"
	default:
		return "unknown"
	}
}

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
