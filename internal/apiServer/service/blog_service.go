// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	domain "github.com/lunarianss/Hurricane/internal/apiServer/domain/blog"
	blogDto "github.com/lunarianss/Hurricane/internal/apiServer/dto/blog"
	"github.com/lunarianss/Hurricane/internal/apiServer/model/v1"
	"github.com/lunarianss/Hurricane/internal/pkg/field"
)

type BlogService struct {
	BlogDomain *domain.BlogDomain
}

func NewBlogService(blogDomain *domain.BlogDomain) *BlogService {
	return &BlogService{BlogDomain: blogDomain}
}

func (s *BlogService) List(ctx context.Context, page, pageSize int) ([]*blogDto.Blog, int64, error) {
	blogs, count, err := s.BlogDomain.BlogRepo.List(ctx, page, pageSize)
	blogsResponse := make([]*blogDto.Blog, 0, 10)

	if err != nil {
		return nil, 0, err
	}

	for _, blog := range blogs {
		blogDto := blogDto.Convert(blog)
		blogsResponse = append(blogsResponse, blogDto)
	}

	return blogsResponse, count, nil
}

func (s *BlogService) Get(ctx context.Context, id int64) (*blogDto.Blog, error) {

	blog, err := s.BlogDomain.BlogRepo.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	blogDto := blogDto.Convert(blog)

	return blogDto, nil
}

func (s *BlogService) Delete(ctx context.Context, id int64) error {
	return s.BlogDomain.BlogRepo.Delete(ctx, id)
}

func (s *BlogService) Create(ctx context.Context, params *blogDto.CreateBlogRequest) (*model.Blog, error) {

	blog := &model.Blog{
		Title:       params.Title,
		BlogPic:     params.BlogPic,
		Content:     params.Content,
		Words:       params.Words,
		ReadTime:    params.ReadTime,
		CategoryId:  params.CategoryId,
		Description: params.Description,
		IsPublished: field.BitBool(params.IsPublished),
		IsTop:       field.BitBool(params.IsTop),
	}

	return s.BlogDomain.BlogRepo.Create(ctx, blog)
}

func (s *BlogService) Update(ctx context.Context, blogId int64, params *blogDto.UpdateBlogRequest) (*model.Blog, error) {

	blog := &model.Blog{
		Id:          blogId,
		Title:       params.Title,
		BlogPic:     params.BlogPic,
		Content:     params.Content,
		Words:       params.Words,
		ReadTime:    params.ReadTime,
		CategoryId:  params.CategoryId,
		Description: params.Description,
		IsPublished: field.BitBool(params.IsPublished),
		IsTop:       field.BitBool(params.IsTop),
	}

	return s.BlogDomain.BlogRepo.Update(ctx, blog)
}
