package rgcategory

import (
	"errors"
	"rgo/core/rgconfig"
	"rgo/core/rgjson"
	"rgo/core/rgrequest"
	"rgo/util/rghttp"
)

const (
	configRBoxProductHost = "lib_rbox_product_host"
)

/*
 * @Content : rgcategory
 * @Author  : LiJunDong
 * @Time    : 2022-06-20$
 */

type apiResult struct {
	Code    int    `json:"code"`
	Result  bool   `json:"result"`
	Message string `json:"message"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		Content []CategoryInfo `json:"content"`
	} `json:"data"`
	TotalPages    int    `json:"totalPages"`
	TotalElements int    `json:"totalElements"`
	Size          int    `json:"size"`
	Page          int    `json:"page"`
	Msg           string `json:"msg"`
}

type CategoryInfo struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	ParentId        int    `json:"parentId"`
	ImageMobilePath string `json:"imageMobilePath"`
	ImageWebPath    string `json:"imageWebPath"`
	InImageList     int    `json:"inImageList"`
	NameInImageList string `json:"nameInImageList"`
	Sort            int    `json:"sort"`
	Alias           string `json:"alias"`
	TaxCode         string `json:"taxCode"`
	OnShelf         int    `json:"onShelf"`
	Delivery        int    `json:"delivery"`
	Status          int    `json:"status"`
	Leaf            int    `json:"leaf"`
	CatePath        string `json:"catePath"`

	//Id              int     `json:"id"`
	//Name            string  `json:"name"`
	//ParentId        int     `json:"parentId"`
	//ImageMobilePath *string `json:"imageMobilePath"`
	//ImageWebPath    *string `json:"imageWebPath"`
	//InImageList     *int    `json:"inImageList"`
	//NameInImageList *string `json:"nameInImageList"`
	//Sort            int     `json:"sort"`
	//Alias           *string `json:"alias"`
	//TaxCode         *string `json:"taxCode"`
	//OnShelf         int     `json:"onShelf"`
	//Delivery        *int    `json:"delivery"`
	//Status          int     `json:"status"`
	//Leaf            *int    `json:"leaf"`
	//CatePath        string  `json:"catePath"`
}

// GetAllCategory 获取所有的分类
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-20
func GetAllCategory(this *rgrequest.Client) (data []CategoryInfo, err error) {
	host := rgconfig.GetStr(configRBoxProductHost)
	if host == "" {
		return data, errors.New("未配置" + configRBoxProductHost)
	}
	url := host + "/categories?page=0&size=3000"
	httpClient := rghttp.Client{
		Param:  nil,
		Method: "GET",
		Header: nil,
		Url:    url,
		This:   this,
	}
	returnData, err := httpClient.GetApi()
	if err != nil {
		return data, err
	}
	result := new(apiResult)
	err = rgjson.UnMarshel([]byte(returnData), result)
	if err != nil {
		return data, err
	}
	if !result.Result || result.Code != 200 {
		return data, errors.New(result.Message)
	}
	return result.Data.Content, err
}
