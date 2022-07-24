package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogRouter struct {
	Service BlogService
}

func (r BlogRouter) AddRoutesToGroup(group *gin.RouterGroup) {
	group.GET("/post", r.GetPostsPage)
	group.POST("/post", r.CreatePost)
}

func (r BlogRouter) CreatePost(ctx *gin.Context) {
	post := BlogPost{}
	err := ctx.BindJSON(&post)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Printf("failed to bind post: %s", err)
		return
	}

	err = post.Validate()
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Printf("invalid post: %s", err)
		return
	}

	err = r.Service.AddPost(post)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("error adding post: %s", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (r BlogRouter) GetPostsPage(ctx *gin.Context) {
	offsetText := ctx.Query("offset")
	limitText := ctx.Query("limit")
	if offsetText == "" {
		offsetText = "0"
	}
	if limitText == "" {
		limitText = "10"
	}
	offset, err := parseInt(offsetText)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Printf("failed to parse offset: %s", err)
		return
	}
	if offset < 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Printf("bad offset parameter: %v", offset)
		return
	}
	limit, err := parseInt(limitText)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Printf("failed to parse limit: %s", err)
		return
	}
	if limit < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.Printf("bad limit parameter: %v", limit)
		return
	}
	page, err := r.Service.GetPostsPage(offset, limit)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("error getting page: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, page)
}

func parseInt(text string) (int, error) {
	num, err := strconv.ParseInt(text, 10, 0)
	if err != nil {
		return 0, err
	}
	return int(num), nil
}
