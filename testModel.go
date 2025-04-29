package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lee31802/comment_lib/ginerrors"
	"github.com/lee31802/comment_lib/gweb"
	"net/http"
)

type Module struct {
	client Client
}

type GetTagsResponse struct {
	Tags []string `json:"rating_star_to_tags"`
	Id   int      `json:"id"`
}

type Resp struct {
	Id   int             `json:"id"`
	Last *uint64         `json:"last"`
	User GetTagsResponse `json:"user"`
}

type Req struct {
	gweb.Request
	StoreID  int     `path:"store_id" desc:"stpreid"`
	LastID   *uint64 `json:"last_id"`
	PageSize *uint32 `json:"page_size"`
	Base     req     `json:"base"`
}

type req struct {
	Id int  `json:"id" desc:"id"`
	Ac Resp `json:"ac"`
}

type Client interface {
	GetTags(ctx context.Context, in *GetDriverTagsRequest) (*GetTagsResponse, error)
	Query(ctx context.Context, in *Req) (*Resp, error)
}
type GetDriverTagsRequest struct {
}

func NewRatingModule() *Module {
	return &Module{}
}
func (m *Module) Init(r gweb.Router) {
	group := r.Group("api/buyer/rating")
	{
		group.GET("/test", m.GetDriverTags)
		group.POST("/test/:store_id/listing", m.QueryStoreRatingV2)
	}
}

func (req *Req) Parse(c *gin.Context) ginerrors.Error {
	return nil
}

func (req *Req) Validate() ginerrors.Error {
	if req.StoreID == 0 || req.LastID == nil {
		return ginerrors.ErrorParamsInvalid
	}
	return ginerrors.Success
}

func (m *Module) QueryStoreRatingV2(ctx context.Context, req *Req) gweb.Response {
	getTagsResp := GetTagsResponse{
		Tags: []string{"llllxxx"},
	}

	return gweb.JSONResponse(http.StatusOK, ginerrors.Success, Resp{
		User: getTagsResp,
		Id:   req.StoreID,
		Last: req.LastID,
	})
}

func (m *Module) GetDriverTags(ctx context.Context) gweb.Response {
	return gweb.JSONResponse(http.StatusOK, ginerrors.Success, "sucess")
}
