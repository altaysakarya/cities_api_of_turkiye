package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
)

var cities map[string]map[string]interface{}

func main() {
	data, err := os.ReadFile("cities.json")
	if err != nil {
		log.Fatalf("File cannot read %v", err)
	}

	err = json.Unmarshal(data, &cities)
	if err != nil {
		log.Fatalf("Json cannot parse %v", err)
	}

	router := gin.Default()

	router.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"detail": "Server is running"})
	})

	router.GET("cities", func(ctx *gin.Context) {
		var cityNames []string
		for key := range cities {
			cityNames = append(cityNames, key)
		}
		sort.Strings(sort.StringSlice(cityNames))
		ctx.JSON(http.StatusOK, cityNames)
	})

	router.GET("cities/:city", func(ctx *gin.Context) {
		cityName := ctx.Param("city")
		city, exists := cities[cityName]
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "City not found"})
			return
		}
		ctx.JSON(http.StatusOK, city["il"])
	})

	router.GET("cities/:city/districts", func(ctx *gin.Context) {
		cityName := ctx.Param("city")
		city, exists := cities[cityName]

		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "City not found"})
			return
		}
		var districts []interface{}
		for key := range city {
			if key != "il" {
				districts = append(districts, key)
			}
		}
		ctx.JSON(http.StatusOK, districts)
	})

	router.GET("cities/:city/districts/:district", func(ctx *gin.Context) {
		cityName := ctx.Param("city")
		districtName := ctx.Param("district")
		city, exists := cities[cityName]
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "City not found"})
			return
		}
		district, districtExists := city[districtName]

		if !districtExists {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "District not found"})
			return
		}
		ctx.JSON(http.StatusOK, district)
	})

	router.Run(":8080")
}
