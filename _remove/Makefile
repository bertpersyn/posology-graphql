SHELL = /bin/bash
PACKAGE ?= sam

#Output purposes
OUTPUT_DIR = $(CURDIR)/output
BIN_OUTPUT_DIR      = $(OUTPUT_DIR)/bin
TEST_OUTPUT_DIR     = $(OUTPUT_DIR)/test
DIRS=$(BIN_OUTPUT_DIR) $(TEST_OUTPUT_DIR)
$(shell mkdir -p $(DIRS))

#Build flags
LDFLAGS ?= "-X 'main.version=$(VERSION)' -X 'main.gitCommit=$(GIT_COMMIT)' -X 'main.application=$(PACKAGE)'"
LINT_FLAGS ?= -E golint -E errcheck --out-format checkstyle
BUILD_FLAGS ?= "./cmd/app"
TEST_FLAGS ?= "-tags=unit"
PACKAGE_EXTENSION ?=  $(shell if [ "$(GOOS)" == windows ]; then echo .exe; fi)
GOPROXY=https://dev-athens.be-mobile.biz
GONOSUMDB=bitbucket.org/be-mobile
GO111MODULE=on
CGO_ENABLED=0
TARGET ?= final

.SILENT: ; # no need for @
.ONESHELL: ; # recipes execute in same shell
.NOTPARALLEL: ; # wait for this target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell
.PHONY: version build docs test scripts api cmd configs examples

deps:
	go get -u github.com/a-h/generate/...

xml-to-json:
	docker build -f scripts/xmltojson/Dockerfile -t xml-to-json:latest .
	docker run --rm -v $(CURDIR)/pkg/:/tmp xml-to-json:latest

structs:
	schema-generate -p types -o $(CURDIR)/pkg/types/sam_gen.go pkg/SAM.jsonschema

build: ## Build the app
	go build -o $(BIN_OUTPUT_DIR)/$(PACKAGE)$(PACKAGE_EXTENSION) --ldflags=$(LDFLAGS) $(BUILD_FLAGS)


quicktype-json:
	quicktype $(CURDIR)/data/data.json --lang go -o pkg/types/sam_gen.go --debug all --package types
	sed -i s/Ns1//g $(CURDIR)/pkg/types/sam_gen.go
	sed -i s/Ns2//g $(CURDIR)/pkg/types/sam_gen.go
	sed -i s/Ns3//g $(CURDIR)/pkg/types/sam_gen.go
	sed -i s/Ns4//g $(CURDIR)/pkg/types/sam_gen.go
	cat $(CURDIR)/pkg/types/sam_gen.go | awk '{ r = gensub(/`json:"@(.*)"`/, "`json:\"@\\1\" graphql:\"\\1\"`", "g"); print r > "$(CURDIR)/pkg/types/sam_gen.go";}'
	cat $(CURDIR)/pkg/types/sam_gen.go | awk '{ r = gensub(/`json:"ns(.*):(.*)"`/, "`json:\"ns\\1:\\2\" graphql:\"!\\2\"`", "g"); print r > "$(CURDIR)/pkg/types/sam_gen.go";}' 
	
quicktype-schema:
	quicktype -s schema $(CURDIR)/data/schemas/json/schema.json -o pkg/types/sam_gen.go --debug all
	#--debug print-graph,print-reconstitution,print-gather-names,print-transformations,print-schema-resolving,print-times,provenance > out
	#Comma separated debug options: print-graph, print-reconstitution, print-gather-names, print-transformations, print-schema-resolving, print-times, provenance 


xsd-to-json-schema: #does not work
	docker build -f scripts/xsdtojsonschema/Dockerfile -t xsd-to-json-schema:latest .
	sudo docker run --rm -v $(CURDIR)/data/schemes/xsd:/tmp xsd-to-json-schema:latest --debug-info /tmp/common/SamCommonItems.xsd

xsd-to-json-schema3: #does not work
	docker build -f scripts/xsdtojsonschema3/Dockerfile -t xsd-to-json-schema3:latest .

jsonix-schema-compiler: #nope
	docker build -f scripts/jsonix-schema-compiler/Dockerfile -t jsonix-schema-compiler:latest .

sam-xsd-convert: #nope
	docker run --rm -v $(CURDIR)/data/schemas/xsd:/data -v $(CURDIR):/tmp jsonix-schema-compiler:latest -logLevel TRACE -generateJsonSchema -p SAM -d /tmp export/SamExport.xsd

sam-remove-ns:
	sed -i s/ns.://g /media/bertp/Data/SAM/2636-2020-06-03-13-21-08/AMP.xml
	sed -i s/ns.://g /media/bertp/Data/SAM/2636-2020-06-03-13-21-08/REF.xml
	sed -i s/ns.://g /media/bertp/Data/SAM/2636-2020-06-03-13-21-08/VMP.xml

#search regex [^func (t].*\*[^xsdT]+
#SubmitRecursiveLegalReferenceType	
#RootObjectFullDataType
#FindRecursiveLegalReferenceType