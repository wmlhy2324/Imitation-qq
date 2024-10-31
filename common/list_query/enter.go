package list_query

import (
	"fmt"
	"gorm.io/gorm"
	"lhyim_server/common/models"
)

type Option struct {
	PageInfo models.PageInfo
	Where    *gorm.DB
	Join     string
	Likes    []string
	Preload  []string
}

func ListQuery[T any](db *gorm.DB, model T, option Option) (list []T, count int64, err error) {
	query := db.Where(model) //把结构体自己的查询条件查了
	if option.PageInfo.Key != "" && len(option.Likes) > 0 {
		likeQuery := db.Where("")
		for index, column := range option.Likes {
			if index == 0 {
				likeQuery.Where(fmt.Sprintf("%s like '%%?%%'", column), option.PageInfo.Key)
			} else {
				likeQuery.Or(fmt.Sprintf("%s like '%%?%%'", column), option.PageInfo.Key)
			}

		}
		query.Where(likeQuery)
	}
	if option.Join != "" {
		query = query.Joins(option.Join)
	}
	if option.Where != nil {
		query = query.Where(option.Where)
	}
	//求总数
	query.Model(model).Count(&count)
	//预加载
	for _, s := range option.Preload {
		query = query.Preload(s)
	}
	//分页
	if option.PageInfo.Page <= 0 {
		option.PageInfo.Page = 1
	}
	if option.PageInfo.Limit <= 0 {
		option.PageInfo.Limit = 10
	}
	if option.PageInfo.Sort != "" {
		query = query.Order(option.PageInfo.Sort)
	}
	offset := (option.PageInfo.Page - 1) * option.PageInfo.Limit

	err = query.Limit(option.PageInfo.Limit).Offset(offset).Find(&list).Error
	return
}
