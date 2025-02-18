package webCoin_common

type BizCode int

const SuccessCode BizCode = 0

// 接口统一返回格式
type Result struct {
	Code    BizCode `json:"code"`
	Message string  `json:"message"`
	Data    any     `json:"data"`
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) Success(data any) {
	r.Code = SuccessCode
	r.Message = "success"
	r.Data = data
}

func (r *Result) Fail(code BizCode, msg string) {
	r.Code = code
	r.Message = msg
}

// 统一用该函数返回数据 这里的err可以再封装一层，自己定义一个error，包含code和msg两个属性
func (r *Result) Deal(data any, err error) *Result {
	if err != nil {
		r.Fail(-999, err.Error())
	} else {
		r.Success(data)
	}
	return r
}
