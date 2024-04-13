package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Reg struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	FirstName   string `json:"firstName"`
	MidddleName string `json:"middleName"`
	LastName    string `json:"lastName"`
	PAN         string `json:"pan"`
	DOB         string `json:dob`
	State       string `json:"state"`
	Gender      string `json:"gender"`
}

var db *gorm.DB

func main() {
	// Connect to MySQL database
	dsn := ("root:root@tcp(localhost:3306)/register?charset=utf8mb4&parseTime=True&loc=Local")
	//dsn := "your_mysql_user:your_mysql_password@tcp(your_mysql_host:3306)/your_mysql_database?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate Todo model
	db.AutoMigrate(&Reg{})

	// Set up Gin router
	r := gin.Default()

	// Middleware to enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// API routes
	r.GET("/reg", getRegs)
	r.POST("/reg", createReg)
	r.GET("/reg/:id", getReg)
	r.PUT("/reg/:id", updateReg)
	r.DELETE("/reg/:id", deleteReg)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Handlers
// func getRegs(c *gin.Context) {
// 	//var regis []Reg
// 	var (
// 		regis  []Reg
// 		limit  int
// 		offset int
// 		page   int
// 		total  int64
// 	)

// 	if pageStr := c.Query("page"); pageStr != "" {
// 		page, _ = strconv.Atoi(pageStr)
// 	} else {
// 		page = 1
// 	}
// 	if limitStr := c.Query("limit"); limitStr != "" {
// 		limit, _ = strconv.Atoi(limitStr)
// 	} else {
// 		limit = 10
// 	}

// 	//offset = (page - 1) * limit
// 	offset = (page - 1) * limit

// 	if err := db.Limit(limit).Offset(offset).Find(&regis).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println("regis : ", regis)
// 		return
// 	}
// 	 //count =db.Model(&Reg{}).Count(&total)
// 	totalPages := int(math.Ceil(float64(total) / float64(limit)))

// 	c.JSON(http.StatusOK, regis)
// 	c.JSON(http.StatusOK, gin.H{"data": regis, "totalPages": totalPages})
// }

// ===================================================

func getRegs(c *gin.Context) {
    var regis []Reg
    var total int64
    var totalPages int

    pageStr := c.DefaultQuery("page", "1")
    page, err := strconv.Atoi(pageStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }

    limitStr := c.DefaultQuery("limit", "10")
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
        return
    }

    offset := (page - 1) * limit

    if err := db.Model(&Reg{}).Count(&total).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count records"})
        return
    }

    if err := db.Limit(limit).Offset(offset).Find(&regis).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		// fmt.Println("regis : ", regis)
        return
    }

    totalPages = int(math.Ceil(float64(total) / float64(limit)))
	fmt.Println("totapages: ", totalPages)

    c.JSON(http.StatusOK, gin.H{"data": regis, "totalPages": totalPages})
}


func createReg(c *gin.Context) {
	var reg Reg
	if err := c.BindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&reg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, reg)
}

func getReg(c *gin.Context) {
	var reg Reg
	id := c.Param("id")
	if err := db.First(&reg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, reg)
}

func updateReg(c *gin.Context) {
	var reg Reg
	id := c.Param("id")
	if err := db.First(&reg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err := c.BindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&reg)
	c.JSON(http.StatusOK, reg)
}

func deleteReg(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Reg{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

