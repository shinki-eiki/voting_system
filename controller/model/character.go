package model

import (
	"time"

	"gorm.io/gorm"
)

type Music struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-" `
	UpdatedAt time.Time      `json:"-" `
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         // 音乐名称
	Poll      int            // 得票数
	Work_id   int            // 来源作品ID，可能有多个？
}

type Character struct {
	gorm.Model
	Name string // 角色名称
	Poll int    // 得票数
}

type Work struct {
	gorm.Model
	Name string // 作品名称
	Poll int    // 得票数
}

func (Music) TableName() string {
	return "music"
}
func (Work) TableName() string {
	return "work"
}

type Musics []Music

// 实现Music的排序
// 重写 Len() 方法
func (a Musics) Len() int {
	return len(a)
}

// 重写 Swap() 方法
func (a Musics) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// 重写 Less() 方法， 从大到小排序
func (a Musics) Less(i, j int) bool {
	return a[j].Poll < a[i].Poll
}
