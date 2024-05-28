package order

import (
	"errors"
	"fmt"
	"net/http"
	"tradeengine/server/web/rest/param"
	"tradeengine/server/web/rest/resp"
	serviceInterfaces "tradeengine/service/interfaces"
	"tradeengine/service/order/types"
	"tradeengine/utils/tool"

	"github.com/gin-gonic/gin"
)

type OrderREST struct {
	mainGroup    *gin.RouterGroup
	currentGroup *gin.RouterGroup
	orderSrv     serviceInterfaces.IOrderSrv
}

func NewREST(mainGroup *gin.RouterGroup, srvMngr serviceInterfaces.IServiceManager) *OrderREST {
	rest := &OrderREST{
		mainGroup: mainGroup,
		orderSrv:  srvMngr.OrderService(),
	}
	rest.currentGroup = rest.mainGroup.Group("/order")
	return rest
}

func (r *OrderREST) RegisterRoute() {
	r.currentGroup.POST("/create", r.Create)
	r.currentGroup.PUT("/edit", r.Edit)
	r.currentGroup.DELETE("/delete/:id", r.Delete)
	r.currentGroup.GET("/:id", r.OrderInfo)
	r.currentGroup.GET("/list", r.OrderInfoList)
}

func (r *OrderREST) Create(c *gin.Context) {
	var param param.OrderCreateParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, resp.FailRespObj(err))
		return
	}
	// create order
	if err := r.orderSrv.Create(param); err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Create %s order successful !", types.GetOrderTypeStr(param.OrderType))
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, nil))
}

func (r *OrderREST) Edit(c *gin.Context) {
	user, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, resp.FailRespObj(errors.New("unauthorized")))
		return
	}
	order, ok := user.(*types.Order)
	if !ok {
		c.JSON(http.StatusInternalServerError, resp.FailRespObj(errors.New("internal Server Error")))
		return
	}
	var param param.OrderEditParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, resp.FailRespObj(err))
		return
	}
	param.OwnerID = order.OwnerID
	// edit member
	if err := r.orderSrv.Edit(param); err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Edit %s order %d successful !", types.GetOrderTypeStr(param.OrderType), param.ID)
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, nil))
}

func (r *OrderREST) Delete(c *gin.Context) {
	id, err := tool.ParseUintFromStr(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.FailRespObj(err))
		return
	}
	user, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, resp.FailRespObj(errors.New("unauthorized")))
		return
	}
	order, ok := user.(*types.Order)
	if !ok {
		c.JSON(http.StatusInternalServerError, resp.FailRespObj(errors.New("internal Server Error")))
		return
	}
	param := param.OrderDeleteParam{
		ID:      id,
		OwnerID: order.ID,
	}
	// delete member
	if err := r.orderSrv.Delete(param); err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Delete %s order %d successful !", types.GetOrderTypeStr(param.OrderType), param.ID)
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, nil))
}

func (r *OrderREST) OrderInfo(c *gin.Context) {
	id, err := tool.ParseUintFromStr(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.FailRespObj(err))
		return
	}
	// query member
	param := param.OrderInfoParam{
		ID: id,
	}
	member, err := r.orderSrv.OrderInfo(param)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Get %s order %d info successful !", types.GetOrderTypeStr(param.OrderType), param.ID)
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, member))
}

func (r *OrderREST) OrderInfoList(c *gin.Context) {
	var param param.OrderInfoListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, resp.FailRespObj(err))
		return
	}
	// get member list
	memberList, err := r.orderSrv.OrderInfoList(param)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	memberListLen := len(memberList)
	dataList := make([]interface{}, memberListLen)
	for i, member := range memberList {
		dataList[i] = member
	}
	message := fmt.Sprintf("Get %s order info list successful !", types.GetOrderTypeStr(param.OrderType))
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, dataList...))
}
