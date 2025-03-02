package decl

import (
	"go/ast"
)

type Comment struct {
	lines []string
}

func NewComment(commentMap ast.CommentMap) *Comment {
	comments := &Comment{}
	for _, comment := range commentMap.Comments() {
		if comment == nil {
			continue
		}
		for _, item := range comment.List {
			comments.lines = append(comments.lines, item.Text)
		}
	}
	return comments
}

func (c *Comment) Lines(filter func(string) bool) []string {
	if c == nil {
		return nil
	}
	if filter == nil {
		return c.lines
	}
	lines := []string{}
	for i := range c.lines {
		if filter(c.lines[i]) {
			lines = append(lines, c.lines[i])
		}
	}
	return lines
}
