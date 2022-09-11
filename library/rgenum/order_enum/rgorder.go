package order_enum

const (
	OrderTypeR1 = 1
	OrderTypeYS = 2
	OrderTypeCGMS = 3
	OrderTypeYSMS = 4
	OrderTypeDKH = 6
	OrderTypeCDYS = 7
	OrderTypeGCG = 8
	OrderTypeGCXG = 9
	OrderTypeTG = 11
)
var OrderTypeName = map[int]string{
	OrderTypeR1:"R1订单",
	OrderTypeYS:"预售订单",
	OrderTypeCGMS:"常规秒杀",
	OrderTypeYSMS:"预售秒杀",
	OrderTypeDKH:"大客户",
	OrderTypeCDYS:"拆单预售",
	OrderTypeGCG:"谷仓G单",
	OrderTypeGCXG:"谷仓GX单",
	OrderTypeTG :"团购",
}
/**
获取订单类型
 */
func GetOrderTypeName() map[int]string {
	return OrderTypeName
}

const (
	PayTypeAli = 1
	PayTypeWeChat = 2
	PayTypeUnionPay = 3
	PayTypeOffline  = 4
	PayTypeUnionApp = 5
	PayTypeCashBalance= 6
	PayTypeRgBt = 7
)

var PayTypeName = map[int]string{
	PayTypeAli: "支付宝",
	PayTypeWeChat: "微信",
	PayTypeUnionPay: "银联",
	PayTypeOffline: "线下",
	PayTypeUnionApp: "银联APP",
	PayTypeCashBalance: "现金余额",
	PayTypeRgBt: "锐锢白条",
}
/**
获取支付类型
 */
func GetPayTypeName() map[int]string {
	return PayTypeName
}

const (
	OsOrderDeliveryTypeR1 = 1
	OsOrderDeliveryTypeYS = 2//Deprecated
	OsOrderDeliveryTypeZF = 3
	OsOrderDeliveryTypeSD = 4
	OsOrderDeliveryTypeGC = 5
	OsOrderDeliveryTypeCDC = 7
)

var OsOrderDeliveryTypeName = map[int]string{
	OsOrderDeliveryTypeR1: "自营",
	OsOrderDeliveryTypeYS: "预售",
	OsOrderDeliveryTypeZF: "直发",
	OsOrderDeliveryTypeSD: "省代",
	OsOrderDeliveryTypeGC: "工厂通",
	OsOrderDeliveryTypeCDC: "产地仓",
}

func GetOrderDeliveryTypeName() map[int]string {
	return OsOrderDeliveryTypeName
}