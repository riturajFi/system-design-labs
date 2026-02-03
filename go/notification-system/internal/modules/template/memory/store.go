package memory

import (
	"context"
	"errors"
	"strings"

	"notification-system/internal/core/model"
	"notification-system/internal/core/ports"
)

type Template struct {
	Title string
	Body  string
}

type Store struct {
	templates map[string]map[model.Channel]Template
}

func New(
	templates map[string]map[model.Channel]Template,
) *Store {
	return &Store{templates: templates}
}

func (s *Store) Render(
	_ context.Context,
	templateID string,
	channel model.Channel,
	params map[string]string,
) (ports.Rendered, error) {

	chMap, ok := s.templates[templateID]
	if !ok {
		return ports.Rendered{}, errors.New("template not found")
	}

	tpl, ok := chMap[channel]
	if !ok {
		return ports.Rendered{}, errors.New("template not defined for channel")
	}

	render := func(text string) string {
		out := text
		for k, v := range params {
			out = strings.ReplaceAll(out, "{{"+k+"}}", v)
		}
		return out
	}

	return ports.Rendered{
		Channel: channel,
		Title:   render(tpl.Title),
		Body:    render(tpl.Body),
	}, nil
}
