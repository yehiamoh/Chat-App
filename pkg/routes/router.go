package routes

import "gorm.io/gorm"

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router{
	return &Router{
		db: db,
	}
}