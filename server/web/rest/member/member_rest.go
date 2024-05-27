package member

import (
	"fmt"
	"net/http"
	"tradeengine/server/web/rest/shared"

	"github.com/gin-gonic/gin"
)

type MemberREST struct {
	mainGroup    *gin.RouterGroup
	currentGroup *gin.RouterGroup
}

func NewREST(mainGroup *gin.RouterGroup) *MemberREST {
	rest := &MemberREST{
		mainGroup: mainGroup,
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
	memberSrv := r.SrvManager.MemberService()
	var user shared.MemberEditParam
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// edit member
	if err := memberSrv.Edit(user); err != nil {
		c.JSON(http.StatusOK, shared.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Edit member %s successful !", user.Account)
	c.JSON(http.StatusOK, shared.SuccessRespObj(message, nil))
}

func (r *MemberREST) Delete(c *gin.Context) {
	memberSrv := r.SrvManager.MemberService()
	account := c.Param("account")
	// delete member
	if err := memberSrv.Delete(shared.MemberDeleteParam{
		Account: account,
	}); err != nil {
		c.JSON(http.StatusOK, shared.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Delete member %s successful !", account)
	c.JSON(http.StatusOK, shared.SuccessRespObj(message, nil))
}

func (r *MemberREST) Members(c *gin.Context) {
	memberSrv := r.SrvManager.MemberService()
	// get member list
	memberList, err := memberSrv.Members()
	if err != nil {
		c.JSON(http.StatusOK, shared.FailRespObj(err))
		return
	}
	memberListLen := len(memberList)
	dataList := make([]interface{}, memberListLen)
	for i, member := range memberList {
		dataList[i] = member
	}
	c.JSON(http.StatusOK, shared.SuccessRespObj("", dataList...))
}

func (r *MemberREST) Member(c *gin.Context) {
	memberSrv := r.SrvManager.MemberService()
	account := c.Param("account")
	// query member
	member, err := memberSrv.Member(shared.MemberInfoParam{
		Account: account,
	})
	if err != nil {
		c.JSON(http.StatusOK, shared.FailRespObj(err))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessRespObj("", member))
}
