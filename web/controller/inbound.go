package controller

import (
    "github.com/gin-gonic/gin"
    "x-ui/web/service"
    "x-ui/web/session"
)

// BatchAddInboundForm 批量添加入站的表单结构体
type BatchAddInboundForm struct {
    StartPort int    `json:"startPort" form:"startPort"`
    Count     int    `json:"count" form:"count"`
    Username  string `json:"username" form:"username"`
    Password  string `json:"password" form:"password"`
}

// InboundController 入站控制器结构体
type InboundController struct {
    inboundService *service.InboundService
}

// NewInboundController 创建一个新的入站控制器实例
func NewInboundController(g *gin.RouterGroup) *InboundController {
    a := &InboundController{
        inboundService: service.NewInboundService(),
    }
    a.initRouter(g)
    return a
}

// initRouter 初始化路由
func (a *InboundController) initRouter(g *gin.RouterGroup) {
    g = g.Group("/inbound")
    g.POST("/batchAdd", a.batchAddInbounds)
    g.GET("/list", a.getInbounds)
    // 其他原有路由...
}

// batchAddInbounds 处理批量添加入站的请求
func (a *InboundController) batchAddInbounds(c *gin.Context) {
    form := &BatchAddInboundForm{}
    // 绑定请求参数到表单结构体
    err := c.ShouldBind(form)
    if err != nil {
        session.JsonMsg(c, "批量添加入站", err)
        return
    }
    // 调用服务层的批量添加入站方法
    err = a.inboundService.BatchAddInbounds(form.StartPort, form.Count, form.Username, form.Password)
    session.JsonMsg(c, "批量添加入站", err)
}

// getInbounds 处理获取入站列表的请求
func (a *InboundController) getInbounds(c *gin.Context) {
    // 调用服务层的获取入站列表方法
    inbounds, err := a.inboundService.GetInbounds()
    session.JsonData(c, "获取入站列表", inbounds, err)
}
