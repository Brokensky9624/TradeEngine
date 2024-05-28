package member

import (
	"fmt"
	"net/http"
	"tradeengine/server/web/rest/param"
	"tradeengine/server/web/rest/resp"
	serviceInterfaces "tradeengine/service/interfaces"

	"github.com/gin-gonic/gin"
)

type MemberREST struct {
	mainGroup    *gin.RouterGroup
	currentGroup *gin.RouterGroup
	memberSrv    serviceInterfaces.IMemberSrv
}

func NewREST(mainGroup *gin.RouterGroup, srvMngr serviceInterfaces.IServiceManager) *MemberREST {
	rest := &MemberREST{
		mainGroup: mainGroup,
		memberSrv: srvMngr.MemberService(),
	}
	rest.currentGroup = rest.mainGroup.Group("/member")
	return rest
}

func (r *MemberREST) RegisterRoute() {
	r.currentGroup.PUT("/edit", r.Edit)
	r.currentGroup.DELETE("/delete/:account", r.Delete)
	r.currentGroup.GET("/list", r.Members)
	r.currentGroup.GET("/:account", r.Member)
}

func (r *MemberREST) Edit(c *gin.Context) {
	var user param.MemberEditParam
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// edit member
	if err := r.memberSrv.Edit(user); err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Edit member %s successful !", user.Account)
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, nil))
}

func (r *MemberREST) Delete(c *gin.Context) {
	account := c.Param("account")
	// delete member
	if err := r.memberSrv.Delete(param.MemberDeleteParam{
		Account: account,
	}); err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Delete member %s successful !", account)
	c.JSON(http.StatusOK, resp.SuccessRespObj(message, nil))
}

func (r *MemberREST) Members(c *gin.Context) {
	// get member list
	memberList, err := r.memberSrv.Members()
	if err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	memberListLen := len(memberList)
	dataList := make([]interface{}, memberListLen)
	for i, member := range memberList {
		dataList[i] = member
	}
	c.JSON(http.StatusOK, resp.SuccessRespObj("", dataList...))
}

func (r *MemberREST) Member(c *gin.Context) {
	account := c.Param("account")
	// query member
	member, err := r.memberSrv.Member(param.MemberInfoParam{
		Account: account,
	})
	if err != nil {
		c.JSON(http.StatusOK, resp.FailRespObj(err))
		return
	}
	c.JSON(http.StatusOK, resp.SuccessRespObj("", member))
}
