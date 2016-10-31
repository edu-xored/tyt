package main

import (
	"github.com/satori/go.uuid"
	"time"
)

type IEntity interface {
	GetID() string
	Created(by string)
	Updated(by string)
}

type Entity struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
	CreatedBy string    `json:"created_by,omitempty"`
	UpdatedBy string    `json:"updated_by,omitempty"`
}

func (e *Entity) GetID() string {
	return e.ID
}

func (e *Entity) Created(by string) {
	e.ID = uuid.NewV4().String()
	e.CreatedAt = time.Now().UTC()
	e.CreatedBy = by
}

func (e *Entity) Updated(by string) {
	e.UpdatedAt = time.Now().UTC()
	e.UpdatedBy = by
}

type User struct {
	Entity
	TeamID      string `json:"team_id,omitempty"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Login       string `json:"login"`
	Description string `json:"description,omitempty"`
	Comment     string `json:"comment,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	GravatarID  string `json:"gravatar_id,omitempty"`
	Course      int32  `json:"course"`
	Faculty     string `json:"faculty,omitempty"`
	Group       int32  `json:"group,omitempty"`
	Role        string `json:"role,omitempty"` // student, mentor
	Github      string `json:"github,omitempty"`
	Skype       string `json:"skype,omitempty"`
	Twitter     string `json:"twitter,omitempty"`
	Telegram    string `json:"telegram,omitempty"`
	WebURL      string `json:"web_url,omitempty"` // url to user website
}

type Team struct {
	Entity
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Slug        string   `json:"slug,omitempty"`
	Github      string   `json:"github,omitempty"`   // github URL
	Telegram    string   `json:"telegram,omitempty"` // team chat
	Members     []string `json:"members,omitempty"`  // member ids
}

type Event struct {
	Entity
	UserID  string    `json:"user_id"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	// allow to track just time spent in hours
	Duration int32 `json:"duration"`
}
