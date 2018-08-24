package actions

import (
	"github.com/Filip/blog/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Article)
// DB Table: Plural (articles)
// Resource: Plural (Articles)
// Path: Plural (/articles)
// View Template Folder: Plural (/templates/articles/)

// ArticlesResource is the resource for the Article model
type ArticlesResource struct {
	buffalo.Resource
}

// List gets all Articles. This function is mapped to the path
// GET /articles
func (v ArticlesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	articles := &models.Articles{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Articles from the DB
	if err := q.All(articles); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, articles))
}

// Show gets the data for one Article. This function is mapped to
// the path GET /articles/{article_id}
func (v ArticlesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Article
	article := &models.Article{}

	// To find the Article the parameter article_id is used.
	if err := tx.Find(article, c.Param("article_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, article))
}

// New renders the form for creating a new Article.
// This function is mapped to the path GET /articles/new
func (v ArticlesResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Article{}))
}

// Create adds a Article to the DB. This function is mapped to the
// path POST /articles
func (v ArticlesResource) Create(c buffalo.Context) error {
	// Allocate an empty Article
	article := &models.Article{}

	// Bind article to the html form elements
	if err := c.Bind(article); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	users := []models.User{}
	err1 := tx.Where("email = ?", "splinter1231@gmail.com").All(&users)
	if err1 == nil {
		article.UserID = users[0].ID
		article.User = users[0]
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(article)
	if err != nil {
		return errors.WithStack(err)
	}

	// article2 := &models.Article{}
	// article2.Title = c.Param("content")
	// article2.Content = c.Param("content")
	// tx.Create(article2)

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, article))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Article was created successfully")

	// and redirect to the articles index page
	return c.Render(201, r.Auto(c, article))
}

// Edit renders a edit form for a Article. This function is
// mapped to the path GET /articles/{article_id}/edit
func (v ArticlesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Article
	article := &models.Article{}

	if err := tx.Find(article, c.Param("article_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, article))
}

// Update changes a Article in the DB. This function is mapped to
// the path PUT /articles/{article_id}
func (v ArticlesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Article
	article := &models.Article{}

	if err := tx.Find(article, c.Param("article_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Article to the html form elements
	if err := c.Bind(article); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(article)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, article))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Article was updated successfully")

	// and redirect to the articles index page
	return c.Render(200, r.Auto(c, article))
}

// Destroy deletes a Article from the DB. This function is mapped
// to the path DELETE /articles/{article_id}
func (v ArticlesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Article
	article := &models.Article{}

	// To find the Article the parameter article_id is used.
	if err := tx.Find(article, c.Param("article_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(article); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Article was destroyed successfully")

	// Redirect to the articles index page
	return c.Render(200, r.Auto(c, article))
}
