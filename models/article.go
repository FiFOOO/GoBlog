package models

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
)

type Article struct {
	ID               uuid.UUID    `json:"id" db:"id"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" db:"updated_at"`
	Title            string       `json:"title" db:"title"`
	Content          string       `json:"content" db:"content"`
	User             User         `belongs_to:"user"`
	UserID           uuid.UUID    `json:"user_id" db:"user_id"`
	TitleImage       binding.File `db:"-" form:"TitleImage"`
	PathToTitleImage string       `json:"title_image" db:"title_image"`
}

// String is not required by pop and may be deleted
func (a Article) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Articles is not required by pop and may be deleted
type Articles []Article

// String is not required by pop and may be deleted
func (a Articles) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

func (a *Article) BeforeCreate(tx *pop.Connection) error {
	if !a.TitleImage.Valid() {
		return nil
	}
	dir := filepath.Join(".", "public/assets/media/article-image")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}
	name := strings.Split(a.TitleImage.Filename, ".")
	ext := name[len(name)-1]
	filename := slug.Make(name[0]) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + "." + ext
	f, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	_, err = io.Copy(f, a.TitleImage)
	a.PathToTitleImage = filepath.Join("media/article-image", filename)
	return err
}

func (a *Article) BeforeDestroy(tx *pop.Connection) error {
	dir := filepath.Join("public/assets/", a.PathToTitleImage)
	if err := os.Remove(dir); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Article) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Article) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Title, Name: "Title"},
		&validators.StringIsPresent{Field: a.Content, Name: "Content"},
		&validators.StringIsPresent{Field: a.TitleImage.Filename, Name: "Title Image"},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Article) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Title, Name: "Title"},
		&validators.StringIsPresent{Field: a.Content, Name: "Content"},
	), nil
}
