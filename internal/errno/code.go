package errno

var (
	ErrServer = New(10001, "服务异常，请联系管理员")
	ErrParam  = New(10002, "参数有误")
)
