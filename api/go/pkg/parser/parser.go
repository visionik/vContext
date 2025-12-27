// Package parser provides interfaces and implementations for parsing vAgenda documents.
package parser

import (
	"io"

	"github.com/visionik/vAgenda/api/go/pkg/core"
)

// Parser handles document parsing from various formats.
type Parser interface {
	// Parse reads a document from a reader.
	Parse(r io.Reader) (*core.Document, error)

	// ParseBytes parses a document from a byte slice.
	ParseBytes(data []byte) (*core.Document, error)

	// ParseString parses a document from a string.
	ParseString(s string) (*core.Document, error)
}

// Format represents the document format type.
type Format string

const (
	// FormatJSON represents JSON format.
	FormatJSON Format = "json"
	// FormatTRON represents TRON format.
	FormatTRON Format = "tron"
	// FormatAuto automatically detects the format.
	FormatAuto Format = "auto"
)

// New creates a new parser for the specified format.
func New(format Format) Parser {
	switch format {
	case FormatJSON:
		return NewJSONParser()
	case FormatTRON:
		return NewTRONParser()
	case FormatAuto:
		return NewAutoParser()
	default:
		return NewAutoParser()
	}
}
