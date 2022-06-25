package rguser

import (
	"errors"
	"rgo/core/rgconfig"
	"rgo/core/rgjson"
	"rgo/core/rgrequest"
	"rgo/util/rghttp"
	"strconv"
)

const (
	configCCHost = "lib_cc_host"
)

/*
 * @Content : rguser
 * @Author  : LiJunDong
 * @Time    : 2022-06-10$
 */
type Result struct {
	Result bool `json:"result"`
	Code   int  `json:"code"`
	Data   User `json:"data"`
	Message string `json:"message"`
}

type ApiResult struct {
	Result bool `json:"result"`
	Code   int  `json:"code"`
	Data   interface{} `json:"data"`
	Message string `json:"message"`
}

type User struct {
	Id                       int      `json:"id"`
	TrueName                 string   `json:"true_name"`
	UserName                 string   `json:"user_name"`
	Gender                   int      `json:"gender"`
	Mobile                   string   `json:"mobile"`
	MaskMobile               string   `json:"mask_mobile"`
	MpMask                   string   `json:"mp_mask"`
	OrderType                int      `json:"order_type"`
	CurrentOrderType         int      `json:"current_order_type"`
	Token                    interface{}   `json:"token"`
	DeviceToken              string   `json:"device_token"`
	CustomerType             int      `json:"customer_type"`
	CurrentCustomerType      int      `json:"current_customer_type"`
	CurrentPriceLevel        int      `json:"current_price_level"`
	PriceLevelPart           int      `json:"price_level_part"`
	Location                 []string `json:"location"`
	AllowVirtualCustomerType []int    `json:"allow_virtual_customer_type"`
	Role                     []string `json:"role"`
	RoleName                 struct {
		MEMBER string `json:"MEMBER"`
	} `json:"-"`
	IsEmployee       int    `json:"is_employee"`
	CustomerTypeIcon string `json:"customer_type_icon"`
	Avatar           string `json:"avatar"`
	AssociatedId     interface{}    `json:"associated_id"`
	MemberLevel      int    `json:"member_level"`
	BposLevel        int    `json:"bpos_level"`
	IsLock           int    `json:"is_lock"`
	Store            []struct {
		IsLock      int      `json:"is_lock"`
		Name        string   `json:"name"`
		SaleUserId  int      `json:"sale_user_id"`
		SaleGroupId int      `json:"sale_group_id"`
		Location    []string `json:"location"`
	} `json:"store"`
	Tags           []string      `json:"tags"`
	CreditScore    int           `json:"credit_score"`
	Openid         interface{} `json:"openid"`
	AccountId      string        `json:"account_id"`
	MerchantNo     string        `json:"merchant_no"`
	MerchantId     string        `json:"merchant_id"`
	RegisterTime   string        `json:"register_time"`
	OriginType     interface{}        `json:"origin_type"`
	RsApprovedTime string        `json:"-"`
}

// GetUserInfoById 从CC接口获取用户详情
// @Param   : id int
// @Return  : data User
// @Author  : LiJunDong
// @Time    : 2022-06-10
func GetUserInfoById(this *rgrequest.Client, id int) (user User, err error) {
	host := rgconfig.GetStr(configCCHost)
	if host == "" {
		return user, errors.New("未配置" + configCCHost)
	}
	url := host + "/internal/member/baseById?id=" + strconv.Itoa(id)
	httpClient := rghttp.Client{
		Param:  nil,
		Method: "GET",
		Header: nil,
		Url:    url,
		This:   this,
	}
	data, err := httpClient.GetApi()
	if err != nil {
		return user, err
	}
	apiRes := new(ApiResult)
	err = rgjson.UnMarshel([]byte(data), apiRes)
	if err != nil {
		return user, err
	}
	if apiRes.Code != 200 {
		return user,errors.New("获取用户信息失败："+apiRes.Message)
	}

	result := new(Result)
	err = rgjson.UnMarshel([]byte(data), result)
	if err != nil {
		return user, err
	}
	if !result.Result || result.Code != 200 {
		return user, errors.New(result.Message)
	}
	return result.Data, err
}