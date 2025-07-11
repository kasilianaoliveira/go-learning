package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"pizzaria/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var pizzas []models.Pizza 
const pizzaByIDRoute = "/pizzas/:id"

func main(){
	loadPizzas()
	router := gin.Default()
  router.GET("/pizzas", getPizzas)
	router.POST("/pizzas", createPizzas)
	router.GET(pizzaByIDRoute, getPizzaById)
	router.DELETE(pizzaByIDRoute, deletePizzaById)
	router.PUT(pizzaByIDRoute, updatePizzaById)
  router.Run()
}

func getPizzas(c *gin.Context) {

	if pizzas != nil {
		c.JSON(200, gin.H{
			"pizzas": pizzas,
		})
		return 
	}

	c.JSON(http.StatusNotFound, gin.H {
		"message": "No pizzas found",
		"pizzas": []models.Pizza{},
	})

}
func createPizzas(c *gin.Context){
	var newPizza models.Pizza
	// c.ShouldBindJSON(&newPizza) lê o corpo da requisição HTTP 
	// e tenta converter os dados JSON para preencher a struct newPizza.
	// Se houver erro (por exemplo, o JSON está malformado ou não bate com 
	// os tipos da struct), ele entra no if e retorna uma resposta de erro.
	if err := c.ShouldBindJSON(&newPizza); err != nil {
		c.JSON(400, gin.H{
			"ERRO":err.Error(),
		})
		return 
	}
	newPizza.ID = len(pizzas) + 1
	pizzas = append(pizzas, newPizza)
	savePizzas()

	c.JSON(201, newPizza) 
}
func getPizzaById(c *gin.Context){
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{
			"erro": err.Error(),
		})
		return
	}

	for _, p := range pizzas {
		if p.ID == id{
			c.JSON(http.StatusOK, gin.H{
				"pizza": p,
			})
			return
		} 
	}

	c.JSON(http.StatusNotFound, gin.H {
		"message": "Pizza not found",
	})
}

func deletePizzaById(c *gin.Context){
	idParam := c.Param("id")
	idConverted, err := strconv.Atoi(idParam)
	if(err != nil) {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": err.Error(),
		})
		return
	}
	for i, p := range pizzas {
		if p.ID == idConverted{
			pizzas = append(pizzas[:i],pizzas[i+1:]... )
			savePizzas()
			c.JSON(200, gin.H{
					"message": "Pizza deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Pizza not found",
	})
}
func updatePizzaById(c *gin.Context){
	
	idParam := c.Param("id")
	idConverted, err := strconv.Atoi(idParam)
	if(err != nil) {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": err.Error(),
		})
		return
	}

	var updatedPizza models.Pizza
	if err := c.ShouldBindJSON(&updatedPizza); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": err.Error(),
		})
		return
	}

	for i, p := range pizzas {
		if p.ID == idConverted{
			pizzas[i] = updatedPizza
			pizzas[i].ID = idConverted
			savePizzas()
			c.JSON(http.StatusOK, gin.H{
					"message": "Pizza updated successfully",
					"pizza": pizzas[i],
			})
			return
		}
	}
}


func loadPizzas(){
	file, err := os.Open("dados/pizzas.json")
	if err != nil {
		fmt.Println("Error file: ",err)
		return 
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pizzas); err != nil {
		fmt.Println("Error decoding json: ", err)
	}
}
func savePizzas() {
	file, err := os.Create("dados/pizzas.json")
	if err != nil {
		fmt.Println("Error file: ",err)
		return 
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&pizzas); err != nil {
		fmt.Println("Error encoding json: ", err)
	}
}