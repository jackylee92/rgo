package product_enum

const (
	EntityTypeGoods     = 1
	EntityTypeAccessory = 2
	EntityTypeComponent = 256
)

var EntityTypeName = map[int]string{
	EntityTypeGoods:     "商品",
	EntityTypeAccessory: "配件",
	EntityTypeComponent: "组件",
}

func GetEntityTypeNames() map[int]string {
	return EntityTypeName
}

const (
	SalesModeR1 = 1
	SalesModeR4 = 3
	SalesModeR5 = 5
)

var EntitySalesMode = map[int]string{
	SalesModeR1: "自营",
	SalesModeR4: "大菠萝",
	SalesModeR5: "工厂通",
}

/**
	获取商品的销售模式
 */
func GetEntitySalesMode() map[int]string {
	return EntitySalesMode
}

const (
	TradeTypeM = 1
	TradeTypeP = 2
)

const (
	CommodityGoods = 1
	NotCommodityGoods = 0
)

var EntityCommodityType = map[int]string{
	CommodityGoods: "大通",
	NotCommodityGoods: "非大通",
}


var EntityTradeType = map[int]string{
	TradeTypeM: "贸易",
	TradeTypeP: "平台",
}

/**
获取商品的tradeType
 */
func GetEntityTradeType() map[int]string {
	return EntityTradeType
}

/**
获取商品大通非大通属性
 */
func GetEntityCommodityName() map[int]string {
	return EntityCommodityType
}
