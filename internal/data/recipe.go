package data

import "github.com/gin-gonic/gin"

var Recipes = map[int64]gin.H{
	1: {
		"id":   1,
		"name": "ขนมครก",
	},
	2: {
		"id":   2,
		"name": "ไข่เจียว",
	},
}
