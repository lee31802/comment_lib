package main

import (
	"context"
	"github.com/lee31802/comment_lib/errors"
	"github.com/lee31802/comment_lib/ginserver"
	"net/http"
)

type RatingModule struct {
	client Client
}

type GetTagsResponse struct {
	RatingStarToTags []string `json:"rating_star_to_tags"`
	Id               int      `json:"id"`
}

type Resp struct {
	Id   int             `json:"id"`
	Last *uint64         `json:"last"`
	User GetTagsResponse `json:"user"`
}

type QueryStoreRatingReq struct {
	ginserver.Request
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
	GetDriverTags(ctx context.Context, in *GetDriverTagsRequest) (*GetTagsResponse, error)
	QueryStoreRatingV2(ctx context.Context, in *QueryStoreRatingReq) (*Resp, error)
}
type GetDriverTagsRequest struct {
}

func NewRatingModule() *RatingModule {
	return &RatingModule{}
}
func (m *RatingModule) Init(r ginserver.Router) {
	group := r.Group("api/buyer/rating")
	{
		group.GET("/tags/driver", m.GetDriverTags)
		group.POST("/store/:store_id/store-rating/action/listing", m.QueryStoreRatingV2)
	}
}

func (req *QueryStoreRatingReq) Validate() errors.Error {
	if req.StoreID == 0 || req.LastID == nil {
		return errors.ErrorParamsInvalid
	}
	return nil
}

func (m *RatingModule) QueryStoreRatingV2(ctx context.Context, req *QueryStoreRatingReq) ginserver.Response {
	getTagsResp := GetTagsResponse{
		RatingStarToTags: []string{"llllxxx"},
	}

	return ginserver.JSONResponse(http.StatusOK, errors.Success, Resp{
		User: getTagsResp,
		Id:   req.StoreID,
		Last: req.LastID,
	})
}

func (m *RatingModule) GetDriverTags(ctx context.Context) string {

	//getTagsResp := &GetTagsResponse{
	//	Id:               11,
	//	RatingStarToTags: []string{"llll"},
	//}

	return "success"
}
