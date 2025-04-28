package ginerrors

// Error definitions
var (
	//common error : [0,999]
	Success           = New(0, "success")
	ErrorUnKnown      = New(1, "unknown")
	ErrorParseRequest = New(2, "parse request error")

	ErrorParamsInvalid    = New(100, "params invalid")
	ErrorDuplicateRequest = New(101, "duplicate request")
	ErrorUserSigCal       = New(102, "usersig cal error")
	ErrorJSONMarshal      = New(103, "json marshal error")
	ErrorJSONUnMarshal    = New(104, "json unmarshal error")

	ErrorDBInit       = New(200, "database init error")
	ErrorRedisInit    = New(201, "redis init error")
	ErrorCacheInit    = New(202, "cache init error")
	ErrorDBOperate    = New(203, "database operate error")
	ErrorRedisOperate = New(204, "redis operate error")
	ErrorCacheOperate = New(205, "cache operate error")
	ErrorDBTxBegin    = New(206, "database tx begin error")
	ErrorDBTxCommit   = New(207, "database tx commit error")
	ErrorDataNotFound = New(208, "data not found")

	ErrorRPCCall     = New(300, "rpc call error")
	ErrorCallTimeout = New(301, "call timeout error")
	ErrorBrokenPipe  = New(302, "broken pipe")

	// others
	ErrorUnknown           = New(500, "server internal error")
	ErrorInvalidValidation = New(501, "Invalid Validation Error")
	ErrorValidation        = New(502, "Validation Error")
)
