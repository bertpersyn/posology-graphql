package vmpparser

import (
	"bufio"
	"io"
	"strconv"

	"github.com/bertpersyn/posology-graphql/internal/sam/parser"
	xmlparser "github.com/tamerh/xml-stream-parser"
)

type Parser struct {
	p                 *xmlparser.XMLParser
	VmpPosologyNote   map[int]string
	VmpNoPosologyNote []int
}

func New(rd io.Reader) (*Parser, error) {
	br := bufio.NewReaderSize(rd, 65536)
	p := xmlparser.NewXMLParser(br, parser.Vmp).EnableXpath()
	return &Parser{
		p:                 p,
		VmpPosologyNote:   make(map[int]string),
		VmpNoPosologyNote: []int{},
	}, nil
}

func (p *Parser) Parse() error {
	for xml := range p.p.Stream() {
		vmpCode, err := strconv.Atoi(xml.Attrs[parser.ACode])
		if err != nil {
			return err
		}
		posologNoteElement, err := xml.SelectElement("Data[not(attribute::to) and attribute::from!='']//CommentedClassification//Data[not(attribute::to) and attribute::from!='']//PosologyNote")
		if err != nil {
			return err
		}
		if posologNoteElement == nil {
			p.VmpNoPosologyNote = append(p.VmpNoPosologyNote, vmpCode)
			continue
		}
		p.VmpPosologyNote[vmpCode] = posologNoteElement.Childs[parser.Nl][0].InnerText
	}
	return nil
}
