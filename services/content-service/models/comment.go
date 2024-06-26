package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID        int        `gorm:"primary_key;column:id"`
	UserID    int        `gorm:"column:user_id"`
	PostID    int        `gorm:"column:post_id"`
	Content   string     `gorm:"column:content"`
	Votes     []UserVote `gorm:"foreignKey:UserVoteID"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Comment) TableName() string {
	return "comments"
}
