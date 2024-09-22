package model

import (
	"gorm.io/gorm"
)

// 记录用户投票的结构体
type Vote struct {
	gorm.Model
	User_id int // 投票用户
	Item_id int // 目标对象
	Part    int // 投票部分
}

func (Vote) TableName() string {
	return "vote"
}
