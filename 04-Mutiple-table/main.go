package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       int
	Name     string
	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt gorm.DeletedAt `gorm:"index"` // 软删除字段

	UserDetail *UserDetail `gorm:"foreignKey:UserID"` // 一对一关系
	VideoList  []*Video    `gorm:"foreignKey:UserID"` // 一对多关系
}

type UserDetail struct {
	ID      int
	Address string

	UserID int   `gorm:"not null;index"`    // 外键
	User   *User `gorm:"foreignKey:UserID"` // 一对一关系
}

type Video struct {
	ID    int
	Title string

	UserID int   `gorm:"Default:NULL;index"` // 外键
	User   *User `gorm:"foreignKey:UserID"`  // 一对多关系

	TagList []*Tag `gorm:"many2many:video_tags;joinForeignKey:VideoID;joinReferences:TagID"` // 多对多关系
}

type Tag struct {
	ID   int
	Name string

	VideoList []*Video `gorm:"many2many:video_tags;joinForeignKey:VideoID;joinReferences:TagID"` // 多对多关系
}

type VideoTag struct {
	// 联合主键
	VideoID int `gorm:"primaryKey"`
	TagID   int `gorm:"primaryKey"`

	// 其他字段
	Stutus int
	Sort   int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除字段
}

// 初始化
var DB *gorm.DB // 全局数据库连接对象

func init() {
	// 连接数据库
	db, err := gorm.Open(mysql.Open("root:123456@tcp(192.168.12.143:3306)/multi_tab"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}

func migrate() {
	// 自动迁移
	err := DB.AutoMigrate(
		&User{},
		&UserDetail{},
		&Video{},
		&Tag{},
		&VideoTag{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
}

func SelectAll() {
	var videos []Video
	DB.Preload("TagList").Find(&videos)
	for _, video := range videos {
		fmt.Println(video)
		for _, tag := range video.TagList {
			fmt.Print(*tag, " ")
		}
		fmt.Println()
	}
}

func main() {
	migrate()
	// 一对一关系
	// 插入数据
	// DB.Create(&User{
	// 	Name: "Alice",
	// 	UserDetail: &UserDetail{
	// 		Address: "123 Main St",
	// 	},
	// })

	// DB.Create(&UserDetail{
	// 	UserID:  1,
	// 	Address: "456 Elm St",
	// })

	// 查询数据
	// var user User
	// DB.Preload("UserDetail").First(&user, 1) // Preload 预加载关联数据
	// fmt.Println(user, user.UserDetail)

	// var userDetail UserDetail
	// DB.Preload("User").First(&userDetail, 1)
	// fmt.Println(userDetail, userDetail.User)

	// 删除数据
	// 级联删除：删除用户时，自动删除关联的用户详情
	// DB.Select("UserDetail").Delete(&User{}, 1)

	// 一对多关系
	// 插入数据
	// DB.Create(&User{
	// 	Name: "Brain",
	// 	VideoList: []*Video{
	// 		{Title: "Go Tutorial"},
	// 		{Title: "GORM Tutorial"},
	// 	},
	// })

	// DB.Create(&Video{
	// 	Title:  "Python Tutorial",
	// 	UserID: 2,
	// })

	// 查询数据
	// var user2 User
	// DB.Preload("VideoList").First(&user2, 2)

	// fmt.Println(user2)
	// for _, video := range user2.VideoList {
	// 	fmt.Println(*video)
	// }

	// var video Video
	// DB.Preload("User").First(&video, 1)
	// fmt.Println(video)

	// 更新数据
	// DB.Model(&Video{}).Where("id = ?", 2).Update("UserID", 1)

	// 删除数据
	// 删除用户时，设置关联视频的 UserID 为 NULL
	// DB.Select("VideoList").Delete(&User{}, 2)

	// 多对多关系
	// 插入数据
	// DB.Create(&Video{
	// 	Title: "Java Tutorial",
	// 	TagList: []*Tag{
	// 		{Name: "Programming"},
	// 		{Name: "Java"},
	// 		{Name: "Backend"},
	// 	},
	// })

	// DB.Create(&Tag{
	// 	Name: "GORM",
	// 	VideoList: []*Video{
	// 		{Title: "Go Tutorial"},
	// 		{Title: "GORM Tutorial"},
	// 	},
	// })

	// 查询数据
	// var video2 Video
	// DB.Preload("TagList").First(&video2, 1)

	// fmt.Println(video2)
	// for _, tag := range video2.TagList {
	// 	fmt.Print(*tag, " ")
	// }
	// fmt.Println()

	// var tag Tag
	// DB.Preload("VideoList").First(&tag, 2)

	// fmt.Println(tag)
	// for _, video := range tag.VideoList {
	// 	fmt.Println(*video, " ")
	// }
	// fmt.Println()

	// SelectAll()

	// 更新数据
	// 将已有的tag关联到指定的视频
	// var tagList []Tag
	// DB.Find(&tagList, "name IN ?", []string{"Programming", "Backend"})
	// tagArgs := make([]interface{}, 0, len(tagList))
	// for i := range tagList {
	// 	tagArgs = append(tagArgs, &tagList[i])
	// }
	// DB.Model(&Video{ID: 2}).Association("TagList").Append(tagArgs...)

	// SelectAll()

	// 删除数据
	// 删除视频时，自动删除关联的 video_tags 记录，但不删除标签
	// DB.Select("TagList").Delete(&Video{}, 4)

	// 自定义中间表
	// 插入数据
	// DB.Create(&VideoTag{
	// 	VideoID: 1,
	// 	TagID:   1,
	// 	Stutus:  1,
	// 	Sort:    1,
	// })

	// // 批量插入
	// DB.Create(&[]VideoTag{
	// 	{VideoID: 1, TagID: 2, Stutus: 1, Sort: 2},
	// 	{VideoID: 1, TagID: 3, Stutus: 1, Sort: 3},
	// 	{VideoID: 2, TagID: 1},
	// 	{VideoID: 2, TagID: 3},
	// 	{VideoID: 2, TagID: 4},
	// })

	// 操作主表或副表
	// DB.Model(&Video{ID: 1}).Association("TagList").Append(&Tag{ID: 11}, &Tag{ID: 12}) // 关联已有标签

	// 查询数据
	var video3 Video
	DB.Preload("TagList").First(&video3, 1)
	fmt.Println(video3.ID, video3.Title)
	for _, tag := range video3.TagList {
		fmt.Print(*tag, " ")
	}
	fmt.Println()

	var tag2 Tag
	DB.Preload("VideoList").First(&tag2, 1)
	fmt.Println(tag2.ID, tag2.Name)
	for _, video := range tag2.VideoList {
		fmt.Print(*video, " ")
	}
	fmt.Println()

	// 删除数据
	// DB.Model(&Video{ID: 3}).Association("TagList").Delete(&Tag{ID: 11}) // 删除关联，但不删除标签

	// 清空关联
	// DB.Model(&Video{ID: 3}).Association("TagList").Clear() // 清空关联，但不删除标签
}
