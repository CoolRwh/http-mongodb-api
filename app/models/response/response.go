package response

type HttpCode int64

const (
	Success               HttpCode = 100200
	Failed                HttpCode = 100500
	DataNotExist          HttpCode = 100100
	UserPasswordError     HttpCode = 100202
	AuthNotExist          HttpCode = 100301
	AuthFail              HttpCode = 100302
	RequestParamError     HttpCode = 100501
	UserNameNotExist      HttpCode = 100502
	TokenBuildError       HttpCode = 100506
	TokenTimeOut          HttpCode = 100507
	AddDataError          HttpCode = 100508
	SqlExecuteError       HttpCode = 100510
	DeleteSuccess         HttpCode = 100511
	DataDeleteFail        HttpCode = 100515
	DataUpdateError       HttpCode = 100518
	RedisClientError      HttpCode = 100545
	RedisClientCloseError HttpCode = 100546
)

var Menus = map[HttpCode]string{
	Success:               "操作成功",
	Failed:                "操作失败",
	DataNotExist:          "数据不存在",
	UserPasswordError:     "账号密码不正确",
	AuthNotExist:          "认证信息不正确",
	AuthFail:              "校验认证信息失败",
	RequestParamError:     "请求参数错误",
	UserNameNotExist:      "用户名不存在",
	TokenBuildError:       "生成Token错误",
	TokenTimeOut:          "认证信息过期",
	AddDataError:          "增加数据失败",
	SqlExecuteError:       "SQL执行错误",
	DeleteSuccess:         "删除成功",
	DataDeleteFail:        "数据删除失败",
	DataUpdateError:       "数据更新失败",
	RedisClientError:      "Redis连接失败",
	RedisClientCloseError: "Redis连接关闭失败",
}

// Message 消息
type Message struct {
	Code HttpCode `json:"code"`
	Msg  string   `json:"message"`
	Data any      `json:"data"`
}

func Ok(data any) Message {
	return Message{
		Code: Success,
		Msg:  Menus[Success],
		Data: data,
	}
}

func Fail(data any) Message {
	return Message{
		Code: Failed,
		Msg:  Menus[Failed],
		Data: data,
	}
}

func Result(code HttpCode, data any) Message {
	return Message{
		Code: code,
		Msg:  Menus[code],
		Data: data,
	}
}
