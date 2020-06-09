package vmpparser

import (
	"bufio"
	"io"

	"github.com/bertpersyn/posology-graphql/internal/sam/parser"
	xmlparser "github.com/tamerh/xml-stream-parser"
)

type Parser struct {
	p               *xmlparser.XMLParser
	VmpPosologyNote map[string]string
}

func New(rd io.Reader) (*Parser, error) {
	br := bufio.NewReaderSize(rd, 65536)
	p := xmlparser.NewXMLParser(br, parser.Vmp).EnableXpath()
	return &Parser{
		p:               p,
		VmpPosologyNote: make(map[string]string),
	}, nil
}

func (p *Parser) Parse() error {

}
