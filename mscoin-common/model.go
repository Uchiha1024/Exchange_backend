package common


type BizCode int

const SuccessCode BizCode = 0


type Result struct {
	Code BizCode `json:"code"`
	Message string `json:"message"`
	Data any `json:"data"`
}


func NewResult() *Result{
	return &Result{}
}


func (r *Result) Success(data any){
	r.Code = SuccessCode
	r.Message = "success"
	r.Data = data
}

func (r *Result) Fail(code BizCode, message string){
	r.Code = code
	r.Message = message

}

func (r *Result) Deal(data any,err error) *Result {
	if err != nil {
		r.Fail(-999,err.Error())
	}else{
		r.Success(data)
	}
	return r
}
