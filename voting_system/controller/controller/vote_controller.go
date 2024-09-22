/* 用于与投票相关的方法 */
package controller

import (
	"fmt"
	"ginEssential/common"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

// 某用户为某音乐投票
func UserVoteMusic(c *gin.Context) {
	// 获取用户的主键以及音乐的主键
	uid, err := strconv.Atoi(c.PostForm("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "解析表单参数时出错"})
		return
	}
	mid, err := strconv.Atoi(c.PostForm("musicID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "解析表单参数时出错"})
		return
	}
	fmt.Println(uid, mid)

	err = common.UserVoteForMusic(uid, mid)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "error": err})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "投票成功"})
}

// 某用户取消为某音乐投票的记录
func UserCancelVoteMusic(c *gin.Context) {
	// 获取用户的主键
	uid, err := strconv.Atoi(c.PostForm("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "解析表单参数时出错"})
		return
	}

	// 获取音乐的主键
	mid, err := strconv.Atoi(c.PostForm("musicID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "解析表单参数时出错"})
		return
	}
	fmt.Println("Cancel:", uid, mid)

	err = common.CancelUserVoteForMusic(uid, mid)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "error": err})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "取消投票成功"})
}

// 获取音乐的排行榜
func MusicRank(c *gin.Context) {
	ms, err := common.GetMusciRank(20)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "error": err})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": ms,
		"msg":  "获取排行榜成功",
	})
}
