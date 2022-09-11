package customer_enum

/**
	用户类型
 */
const (
	CustomerTypeDTB  = 7
	CustomerTypeQGB = 4
	CustomerTypeRDC = 1
	CustomerTypeC = 2
)

/**
	用户来源
 */
const (
	CustomerOriginTypeRS = 1
	CustomerOriginTypeQP = 2
)




var CustomerTypeName = map[int]string{
	CustomerTypeDTB : "地推B" ,
	CustomerTypeQGB: "全国B",
	CustomerTypeC: "C端",
	CustomerTypeRDC: "RDC用户",
}

var CustomerOriginTypeName = map[int]string{
	CustomerOriginTypeQP: "企拍",
	CustomerOriginTypeRS: "RS",

}

/**
	获取用户类型
 */
func GetCustomerTypeNames() map[int]string {
	return CustomerTypeName
}

/**
	获取用户来源
 */
func GetOriginTypeNames() map[int]string {
	return CustomerOriginTypeName
}

/**
	检查用户类型是否在枚举内
 */
func CheckCustomerType(customerType int) bool {
	_,ok := CustomerTypeName[customerType]
	return ok
}




