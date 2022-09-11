package error_enum

const (
	/**
	token验证失败
	*/
	ErrorCodeTokenInvalid = 401
	/**

	 */
	ErrorCodeDuplicateEntry = -1062
	/**
	参数验证失败
	*/
	ErrorCodeParamValidate = -1001
	/**
	记录不存在
	*/
	ErrorCodeRecordNotExist = -1004
	/**
	记录重复
	*/
	ErrorCodeRecordRepeat = -1006
	/**
	创建出错
	*/
	ErrorCodeCreateError = -1007
	/**
	更新出错
	*/
	ErrorCodeUpdateError = -1008
	/**
	标记失效
	*/
	ErrorCodeSignInvalid = -102
	/**
	标记错误
	*/
	ErrorCodeSignError = -103
	/**
	默认错误码
	*/
	ErrorCodeDefault = -1
)
