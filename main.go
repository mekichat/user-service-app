package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Role      string   `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}


var nextID = 1
var employees []Employee



func main() {

	r := gin.Default()

	r.GET("/employees", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, employees)
	})

	r.GET("/employees/:id", func(ctx *gin.Context) {

		idParam := ctx.Param("id")

		id, err := strconv.Atoi(idParam)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "id must be a number",
			})
			return
		}

		for _, e := range employees {
			if e.ID == id {
				ctx.JSON(http.StatusOK, e)
				return
			}
		}


		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "employee not found",
		})

	})

	r.POST("/employees", func(ctx *gin.Context) {

		var input struct {
			Name string `json:"name"`
			Role string `json:"role"`
		}

		if err := ctx.BindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON",
			})
			return
		}

		if input.Name == "" || input.Role == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "name and role are required",
			})
			return
		}

		newEmployee := Employee{
			ID: nextID,
			Name: input.Name,
			Role: input.Role,
			CreatedAt: time.Now(),
		}

		nextID++

		employees = append(employees, newEmployee)

		ctx.JSON(http.StatusCreated, newEmployee)

	})

	r.PUT("/employees/:id", func(ctx *gin.Context){

		var input struct {
			Name string `json:"name"`
			Role string `json:"role"`
		}
		
		idParam := ctx.Param("id")

		id, err := strconv.Atoi(idParam)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "id must be a number",
			})
			return
		}

		if err := ctx.BindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON",
			})
			return
		}

		if input.Name == "" || input.Role == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "name and role are required",
			})
			return
		}

		for i, e := range employees {
			if e.ID == id {
				e.Name = input.Name
				e.Role = input.Role
				employees[i] = e	
		        ctx.JSON(http.StatusCreated, e)			
				return
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "employee not found",
		})
	})

	r.DELETE("/employees/:id", func(ctx *gin.Context){

			idParam := ctx.Param("id")

			id, err := strconv.Atoi(idParam)

			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error" : "id must be a number",
				})
				return				
			}

			for i, e := range employees {
				if e.ID == id {
					//ctx.JSON(http.StatusOK, e)
					employees = append(employees[:i], employees[i+1:]...)
					ctx.JSON(http.StatusOK, gin.H{
						"message": "employee deleted successfully",
					})
					return
				}
		    }

			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "employee not found",
			})
			
		})

	r.Run(":8080")

}
