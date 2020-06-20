package samparser

import (
	"fmt"
	vmpparser "github.com/bertpersyn/posology-graphql/internal/sam/parser/vmp"
	"os"
	"sync"

	ampparser "github.com/bertpersyn/posology-graphql/internal/sam/parser/amp"

	refparser "github.com/bertpersyn/posology-graphql/internal/sam/parser/ref"
	types "github.com/bertpersyn/posology-graphql/internal/sam/types"
	"github.com/koding/multiconfig"
	"github.com/sirupsen/logrus"
)

type Service struct {
	cfg *types.Cfg

	refParser     *refparser.Parser
	ampParser     *ampparser.Parser
	vmpParser     *vmpparser.Parser
	refParserFile *os.File
	ampParserFile *os.File
	vmpParserFile *os.File
}

func New() (*Service, error) {
	m := multiconfig.NewWithPath("config.yaml")
	// Get an empty struct for your configuration
	cfg := new(types.Cfg)
	m.MustLoad(cfg) // Panic's if there is any error
	logrus.Infof("%+v", cfg)
	s := &Service{cfg: cfg}
	err := s.initRefParser()
	if err != nil {
		return s, err
	}
	err = s.initAmpParser()
	if err != nil {
		return s, err
	}
	err = s.initVmpParser()
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s *Service) initRefParser() error {
	xmlFile, err := os.Open(s.cfg.XML.RefPath)
	if err != nil {
		return fmt.Errorf("could not open xml ref path: %v", err)
	}
	s.refParser, err = refparser.New(xmlFile)
	if err != nil {
		return fmt.Errorf("could not create ref parser: %v", err)
	}
	s.refParserFile = xmlFile
	return nil
}

func (s *Service) initAmpParser() error {
	xmlFile, err := os.Open(s.cfg.XML.AmpPath)
	if err != nil {
		return fmt.Errorf("could not open amp ref path: %v", err)
	}
	s.ampParser, err = ampparser.New(xmlFile)
	if err != nil {
		return fmt.Errorf("could not create amp parser: %v", err)
	}
	s.ampParserFile = xmlFile
	return nil
}

func (s *Service) initVmpParser() error {
	xmlFile, err := os.Open(s.cfg.XML.VmpPath)
	if err != nil {
		return fmt.Errorf("could not open vmp ref path: %v", err)
	}
	s.vmpParser, err = vmpparser.New(xmlFile)
	if err != nil {
		return fmt.Errorf("could not create vmp parser: %v", err)
	}
	s.vmpParserFile = xmlFile
	return nil
}

func (s *Service) ParseAll() error {
	defer func() {
		err := s.ampParserFile.Close()
		if err != nil {
			logrus.Errorf("could not close amp xml file: %v", err)
		}
		err = s.refParserFile.Close()
		if err != nil {
			logrus.Errorf("could not close ref xml file: %v", err)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(3)
	errChan := make(chan error)
	wgDoneChan := make(chan bool)

	go func() {
		logrus.Debugf("start parsing refs")
		errChan <- s.refParser.Parse()
		wg.Done()
	}()

	go func() {
		logrus.Debugf("start parsing amps")
		errChan <- s.ampParser.Parse()
		wg.Done()
	}()

	go func() {
		logrus.Debugf("start parsing vmps")
		errChan <- s.vmpParser.Parse()
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(wgDoneChan)
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				return err
			}
		case <-wgDoneChan:
			logrus.Debugf("parsed all")
			return nil
		}
	}
}

func (s *Service) GetSubstances() map[string]*types.Substance {
	return s.refParser.Substances
}

func (s *Service) GetPharmaceuticalForms() map[string]*types.PharmaceuticalForm {
	return s.refParser.PharmaceuticalForms
}

func (s *Service) GetActualMedicinalProducts() map[string]*types.ActualMedicinalProduct {
	return s.ampParser.ActualMedicinalProducts
}

func (s *Service) GetVmpPosologyNotes() map[int]string {
	return s.vmpParser.VmpPosologyNote
}
