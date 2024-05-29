package stockinfo

import (
	"fmt"
	"net/http"
	"tradeengine/server/web/rest/param"
	"tradeengine/server/web/rest/resp"
	serviceInterfaces "tradeengine/service/interfaces"
	"tradeengine/utils/tool"

	"github.com/gin-gonic/gin"
)

type StockInfoREST struct {
	mainGroup    *gin.RouterGroup
	currentGroup *gin.RouterGroup
	stockInfoSrv serviceInterfaces.IStockInfoSrv
}

func NewREST(mainGroup *gin.RouterGroup, srvMngr serviceInterfaces.IServiceManager) *StockInfoREST {
	rest := &StockInfoREST{
		mainGroup:    mainGroup,
		stockInfoSrv: srvMngr.StockInfoService(),
	}
	rest.currentGroup = rest.mainGroup.Group("/stockinfo")
	return rest
}

func (r *StockInfoREST) RegisterRoute() {
	r.currentGroup.POST("/create", r.Create)
	r.currentGroup.GET("/:id", r.StockInfoInfo)
	r.currentGroup.GET("/list", r.StockInfoList)
}

func (r *StockInfoREST) Create(c *gin.Context) {
	var param param.StockInfoCreateParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, resp.FailRespObj(err))
		return
	}
	// create stock
	if err := r.stockInfoSrv.Create(param); err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	message := "Create stock info successful !"
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, nil))
}

func (r *StockInfoREST) StockInfoInfo(c *gin.Context) {
	id, err := tool.ParseUintFromStr(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.FailRespObj(err))
		return
	}
	// query stock info
	param := param.StockInfoParam{
		ID: id,
	}
	stockInfo, err := r.stockInfoSrv.StockInfo(param)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Get stock info %d successful !", param.ID)
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, stockInfo))
}

func (r *StockInfoREST) StockInfoList(c *gin.Context) {
	var param param.StockInfoListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, resp.FailRespObj(err))
		return
	}
	// get stock info list
	stockInfoList, err := r.stockInfoSrv.StockInfoList(param)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	stockInfoListLen := len(stockInfoList)
	dataList := make([]interface{}, stockInfoListLen)
	for i, stockInfo := range stockInfoList {
		dataList[i] = stockInfo
	}
	message := "Get stock info list successful !"
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, dataList...))
}
