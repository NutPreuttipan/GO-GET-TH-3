package main

import (
	"strings"
	"fmt"
	"github.com/labstack/echo/middleware"
	"net/http"
	
	"github.com/labstack/echo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

func main() {

	e := echo.New()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	mongoHost := viper.GetString("MONGO.HOST")
	mongoUser := viper.GetString("MONGO.USER")
	mongoPass := viper.GetString("MONGO.PASS")
	port := viper.GetString("port")

	conString := fmt.Sprintf("%s:%s@%s",mongoUser,mongoPass,mongoHost)

	session, err := mgo.Dial(conString)
	if err != nil {
		e.Logger.Fatal(err) 
	}

	handler := &handler { mongoDB:session }

	e.Use(middleware.Logger())
	e.GET("/todo/:id", handler.view)
	e.GET("/todos", handler.list)
	e.PUT("/todo/:id", handler.done)
	e.DELETE("/list:id", handler.delete)
	e.POST("/todo", handler.create)
	
	e.Logger.Fatal(e.Start(":"+port))
}

type todo struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"` 
	Topic string `json:"topic" bson:"topic"`
	Done bool `json:"done" bson:"done"`
}

type handler struct {
	mongoDB *mgo.Session
}

func (h *handler) create (c echo.Context) error {
	var newTodo todo
	if err := c.Bind(&newTodo); err != nil {
		return err
	}

	session := h.mongoDB.Copy()
	defer session.Close()

	collection := session.DB("workshop").C("PJ07")
	if err := collection.Insert(newTodo); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newTodo)
}

func (h *handler) list(c echo.Context) error {
	var listToto []todo

	session := h.mongoDB.Copy()
	defer session.Close()
	collection :=  session.DB("workshop").C("PJ07")
	if err := collection.Find(nil).All(&listToto); err != nil {
		return err
	}
	return c.JSON(http.StatusOK,listToto) 
}

func (h *handler) view(c echo.Context) error {
	var newToto todo
	id := bson.ObjectIdHex(c.Param("id"))

	session := h.mongoDB.Copy()
	defer session.Close()
	collection :=  session.DB("workshop").C("PJ07")
	if err := collection.FindId(id).One(&newToto); err != nil {
		return err
	}
	return c.JSON(http.StatusOK,newToto) 
}

func (h *handler) done(c echo.Context) error {
	var newTodo todo
	id := bson.ObjectIdHex(c.Param("id"))

	session := h.mongoDB.Copy()
	defer session.Close()

	collection :=  session.DB("workshop").C("PJ07")
	if err := collection.FindId(id).One(&newTodo); err != nil {
		return err
	}
	newTodo.Done = true
	if err := collection.UpdateId(id,newTodo); err != nil {
		return err
	}
	
	return c.JSON(http.StatusOK,newTodo) 
}

func (h *handler) delete(c echo.Context) error {
	id := bson.ObjectIdHex(c.Param("id"))

	session := h.mongoDB.Copy()
	defer session.Close()

	collection :=  session.DB("workshop").C("PJ07")
	if err := collection.RemoveId(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK,echo.Map {
		"reuslt":"Success",
	}) 
}