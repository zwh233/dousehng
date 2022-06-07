package service

import (
	"demo1/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

//  ---favourite---

func FavoriteAction(c *gin.Context) {
	var req UserFavoriteRequest
	var user repository.User
	var video repository.Video

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "FavoriteAction should bind error",
			},
		})
		return
	}

	// 从token中提取相关信息
	user.Name = c.GetString("username")
	user.ID = c.GetUint("user_id")

	// 创建单例
	favoriteDAO := repository.NewFavoriteDAO()
	userDAO := repository.NewUserDAO()
	videoDAO := repository.NewVideoDAO()

	if err := userDAO.FindUserById(user.ID, (*repository.User)(&user)); err != nil {
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User can't find error",
			},
		})
		return
	}

	if err := videoDAO.FindVideoById(req.VideoId, &video); err != nil {
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Video can't find error",
			},
		})
	}

	if req.ActionType == 1 {
		//点赞操作
		//将该视频加入用户的点赞列表
		if err := favoriteDAO.Favorite(user.ID, req.VideoId); err != nil {
			c.JSON(http.StatusOK, UserFavoriteResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Add into favorite list error",
				},
			})
		}
		//video.FavoriteCount++
		if err := favoriteDAO.AddFavoriteCount(req.VideoId); err != nil {
			c.JSON(http.StatusOK, UserFavoriteResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Add favorite count error",
				},
			})
		}
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "Favorite successful",
			},
		})
	} else if req.ActionType == 2 {
		//取消点赞
		//将该视频从用户的点赞列表移除
		if err := favoriteDAO.UnFavorite(req.UserId, req.VideoId); err != nil {
			c.JSON(http.StatusOK, UserFavoriteResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Delete from favorite list error",
				},
			})
		}
		//video.FavoriteCount--
		if err := favoriteDAO.ReduceFavoriteCount(req.VideoId); err != nil {
			c.JSON(http.StatusOK, UserFavoriteResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Reduce favorite count error",
				},
			})
		}
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "UnFavorite successful",
			},
		})
	} else {
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Favorite action error",
			},
		})
	}
}

//  ---favourite list---

func FavoriteList(c *gin.Context) {
	var req UserFavoriteListRequest
	var videoList []repository.Video

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, UserFavoriteListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "FavoriteList Action should bind error",
			},
			VideoList: nil,
		})
	}

	// 创建单例
	favoriteDAO := repository.NewFavoriteDAO()

	if err := favoriteDAO.FindFavoriteVideoByUid(req.UserId, &videoList); err != nil {
		c.JSON(http.StatusOK, UserFavoriteListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Can't find favorite video list by uid",
			},
			VideoList: nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, UserFavoriteListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "Get favorite list success",
			},
			VideoList: videoList,
		})
	}
}
