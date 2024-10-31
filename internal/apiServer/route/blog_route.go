// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	controller "github.com/Ryan-eng-del/hurricane/internal/apiServer/controller/gin/v1/blog"
	"github.com/Ryan-eng-del/hurricane/internal/apiServer/dao"
	domain "github.com/Ryan-eng-del/hurricane/internal/apiServer/domain/blog"
	"github.com/Ryan-eng-del/hurricane/internal/apiServer/service"
	"github.com/Ryan-eng-del/hurricane/internal/pkg/mysql"
	"github.com/gin-gonic/gin"
)

type blogRoutes struct{}

func (r *blogRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// dao
	blogDao := dao.NewBlogDao(gormIns)

	// domain
	blogDomain := domain.NewBlogDomain(blogDao)

	// service
	blogService := service.NewBlogService(blogDomain)
	blogController := controller.NewBlogController(blogService)

	v1 := g.Group("/v1")
	blogV1 := v1.Group("/blog")

	blogV1.GET("", blogController.List)
	blogV1.GET("/:blogId", blogController.Get)
	blogV1.PUT("/:blogId", blogController.Update)
	blogV1.DELETE("/:blogId", blogController.Delete)
	blogV1.POST("", blogController.Create)
	return nil
}

func (r *blogRoutes) GetModule() string {
	return "blog"
}
