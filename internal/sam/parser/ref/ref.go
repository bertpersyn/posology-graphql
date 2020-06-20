package refparser

import (
	"bufio"
	"github.com/bertpersyn/posology-graphql/internal/sam/parser"
	"io"

	types "github.com/bertpersyn/posology-graphql/internal/sam/types"
	xmlparser "github.com/tamerh/xml-stream-parser"
)

type Parser struct {
	p                   *xmlparser.XMLParser
	Substances          map[string]*types.Substance
	PharmaceuticalForms map[string]*types.PharmaceuticalForm
}

func New(rd io.Reader) (*Parser, error) {
	br := bufio.NewReaderSize(rd, 65536)
	p := xmlparser.NewXMLParser(br, parser.Substance, parser.PharmaceuticalForm)

	return &Parser{
		p:                   p,
		Substances:          make(map[string]*types.Substance),
		PharmaceuticalForms: map[string]*types.PharmaceuticalForm{},
	}, nil
}

func (ref *Parser) Parse() error {
	for xml := range ref.p.Stream() {
		switch xml.Name {
		case parser.Substance:
			if name, exists := xml.Childs[parser.Name]; exists {
				s := &types.Substance{
					Code: xml.Attrs[parser.ACode],
					Name: name[0].Childs[parser.Nl][0].InnerText,
				}
				ref.Substances[s.Code] = s
			}
		case parser.PharmaceuticalForm:
			if name, exists := xml.Childs[parser.Name]; exists {
				s := &types.PharmaceuticalForm{
					Code: xml.Attrs[parser.ACode],
					Name: name[0].Childs[parser.Nl][0].InnerText,
				}
				ref.PharmaceuticalForms[s.Code] = s
			}
		}
	}
	return nil
}
