package ampparser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/bertpersyn/posology-graphql/internal/sam/parser"

	"github.com/bertpersyn/posology-graphql/internal/sam/types"
	xmlparser "github.com/tamerh/xml-stream-parser"
)

type Parser struct {
	p                       *xmlparser.XMLParser
	ActualMedicinalProducts map[string]*types.ActualMedicinalProduct
}

func New(rd io.Reader) (*Parser, error) {
	br := bufio.NewReaderSize(rd, 65536)
	p := xmlparser.NewXMLParser(br, parser.Amp).EnableXpath()
	return &Parser{
		p:                       p,
		ActualMedicinalProducts: make(map[string]*types.ActualMedicinalProduct),
	}, nil
}

func (p *Parser) Parse() error {
	for xml := range p.p.Stream() {
		vmpCode := 0
		var err error
		if strVmpCode, exists := xml.Attrs[parser.AVmpCode]; exists {
			vmpCode, err = strconv.Atoi(strVmpCode)
		}
		if err != nil {
			return err
		}
		amp := &types.ActualMedicinalProduct{
			Code:    xml.Attrs[parser.ACode],
			VmpCode: vmpCode,
		}
		err = p.parseAmpData(xml, amp)
		if err != nil {
			return err
		}
		err = p.parseAmpComponents(xml, amp)
		if err != nil {
			return err
		}
		err = p.parseAmpp(xml, amp)
		if err != nil {
			return err
		}
		p.ActualMedicinalProducts[amp.Code] = amp
	}
	return nil
}

func (p *Parser) parseAmpData(xml *xmlparser.XMLElement, amp *types.ActualMedicinalProduct) error {
	data := p.find(xml, parser.Data)
	from, err := time.Parse(parser.TimeLayout, data.Attrs[parser.AFrom])
	if err != nil {
		return fmt.Errorf("could not parse from for amp data: %v", err)
	}
	amp.From = from
	amp.OfficialName = data.Childs[parser.OfficialName][0].InnerText
	return nil
}

func (p *Parser) parseAmpComponents(xml *xmlparser.XMLElement, amp *types.ActualMedicinalProduct) error {
	ampComponentElement, err := xml.SelectElement(parser.AmpComponent)
	if err != nil {
		return err
	}
	vmpComponentCode := 0
	if strvmpComponentCode, exists := ampComponentElement.Attrs[parser.AVmpComponentCode]; exists {
		vmpComponentCode, err = strconv.Atoi(strvmpComponentCode)
	}
	if err != nil {
		return err
	}

	ampComponentDataElement := p.find(ampComponentElement, parser.Data)
	amp.AmpComponent = types.AmpComponent{
		VmpComponentCode:       vmpComponentCode,
		PharmaceuticalFormCode: ampComponentDataElement.Childs[parser.PharmaceuticalForm][0].Attrs[parser.ACode],
	}
	realActualIngredientElement, err := ampComponentElement.SelectElement(parser.RealActualIngredient)
	if err != nil {
		return err
	}
	realActualIngredientDataElement := p.find(realActualIngredientElement, parser.Data)
	from, err := time.Parse(parser.TimeLayout, realActualIngredientDataElement.Attrs[parser.AFrom])
	if err != nil {
		return fmt.Errorf("could not parse from, element real actual ingredient  %v", err)
	}
	amp.AmpComponent.RealActualIngredient = types.RealActualIngredient{

		From:          from,
		Type:          realActualIngredientDataElement.Childs[parser.Type][0].InnerText,
		SubstanceCode: realActualIngredientDataElement.Childs[parser.Substance][0].Attrs[parser.ACode],
	}
	if strength, exists := realActualIngredientDataElement.Childs[parser.Strength]; exists {
		strenghtV, err := strconv.ParseFloat(strength[0].InnerText, 32)
		if err != nil {
			return err
		}
		amp.AmpComponent.RealActualIngredient.Strength = types.Quantity{
			Value: strenghtV,
			Unit:  realActualIngredientDataElement.Childs[parser.Strength][0].Attrs[parser.AUnit],
		}
	}

	return nil
}

func (p *Parser) parseAmpp(xml *xmlparser.XMLElement, amp *types.ActualMedicinalProduct) error {
	amppElements, err := xml.SelectElements(parser.Ampp)
	if err != nil {
		return err
	}
	ampp := types.Ampp{
		PosologyNotes: []string{},
	}
	for _, amppElement := range amppElements {
		amppDataElement := p.find(amppElement, parser.Data)
		if posologyNote, exists := amppDataElement.Childs[parser.PosologyNote]; exists {
			ampp.PosologyNotes = append(ampp.PosologyNotes, posologyNote[0].Childs[parser.Nl][0].InnerText)
		}
	}
	amp.Ampp = ampp
	return nil
}

func (p *Parser) find(xml *xmlparser.XMLElement, elementName string) *xmlparser.XMLElement {
	var actualData xmlparser.XMLElement
	for _, data := range xml.Childs[elementName] {
		actualData = data
		if _, exists := data.Attrs[parser.ATo]; !exists {
			break
		}
	}
	return &actualData
}
