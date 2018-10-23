package actions

import (
	"log"

	"github.com/Filip/blog/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

// VisitorHandler is visitor page articles
func VisitorHandler(c buffalo.Context) error {
	articles := &models.Articles{}
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}
	q := tx.PaginateFromParams(c.Params())
	if err := q.All(articles); err != nil {
		return errors.WithStack(err)
	}
	c.Set("articles", articles)
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("visitor/index.html"))
}

// VisitorArticleShowHandler is visitor page for showing specific article
func VisitorArticleShowHandler(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Article
	article := &models.Article{}
	messages := &models.Messages{}

	// To find the Article the parameter article_id is used.
	if err := tx.Eager("User").Find(article, c.Param("article_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Eager("User").Where("article_id = ?", c.Param("article_id")).Order("created_at").All(messages); err != nil {
		return c.Error(404, err)
	}

	c.Set("article", article)
	c.Set("messages", messages)

	return c.Render(200, r.HTML("visitor/show.html"))
}

// VisitorSearchArticleHandler return list of search articles
func VisitorSearchArticleHandler(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	articles := &models.Articles{}

	q := tx.PaginateFromParams(c.Params())

	if err := q.Where(`LOWER(title) ~ LOWER(?)`, c.Request().Form.Get("title")).All(articles); err != nil {
		return errors.WithStack(err)
	}

	data := make(map[string]interface{})
	data["articles"] = articles
	data["pagination"] = q.Paginator

	return c.Render(200, r.JSON(data))
}

func CreateMessage(c buffalo.Context) error {
	message := &models.Message{}

	message.Msg = c.Request().Form.Get("msg")
	message.ArticleID = c.Request().Form.Get("article")
	message.UserID = c.Request().Form.Get("user")

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	verrs, err := tx.ValidateAndCreate(message)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(201, r.String("false"))
	}

	return c.Render(201, r.String("true"))
}

type Message struct {
	Msg       string `json:"msg"`
	Article   string `json:"article"`
	UserID    string `json:"user"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func MassageHandler(c buffalo.Context) error {
	var conn, err = upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	clients[conn] = true
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}

		user := &models.User{}
		tx, ok := c.Value("tx").(*pop.Connection)
		if !ok {
			return errors.WithStack(errors.New("no transaction found"))
		}

		if err := tx.Find(user, msg.UserID); err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}
		msg.FirstName = user.FirstName
		msg.LastName = user.LastName
		broadcast <- msg
	}

	return nil
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
