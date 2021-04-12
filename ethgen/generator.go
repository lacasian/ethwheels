package ethgen

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var ErrMismatchingEvent = errors.New("event id does not match given log event")

type EventDefinition struct {
	source string
	event  abi.Event
}

type ContractData struct {
	InputABI string
	Defs     map[string]abi.Event
	Structs  map[string]*tmplStruct
}

type TemplateData struct {
	ContractData
	Package string
	Prefix  string
}

func NewFromABIs(abisDir string, packagePath string) error {
	pkg := filepath.Base(packagePath)
	contracts, err := ReadABIs(abisDir)
	if err != nil {
		return errors.Wrap(err, "read abis")
	}

	for prefix, contract := range contracts {
		code, err := GenerateCode(prefix, pkg, contract)
		if err != nil {
			return errors.Wrap(err, "generating code")
		}

		outFN := path.Join(packagePath, strings.ToLower(prefix)+".go")
		err = ioutil.WriteFile(outFN, code, 0644)
		if err != nil {
			return errors.Wrapf(err, "writing %s", outFN)
		}
		log.Infof("generated %s", outFN)
	}

	return nil
}

func ReadABIs(dir string) (map[string]ContractData, error) {
	definitions := make(map[string]ContractData)
	fc, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "read abi json dir %s", dir)
	}

	for _, f := range fc {
		if path.Ext(f.Name()) == ".json" {
			fn := f.Name()
			n := gopherize(strings.TrimSuffix(fn, ".json"))
			cd, err := ProcessFile(dir, fn)
			if err != nil {
				return nil, errors.Wrapf(err, "process file %s", fn)
			}
			if cd != nil {
				definitions[n] = *cd
			}
		}
	}
	return definitions, nil
}

func ProcessFile(dir string, fileName string) (*ContractData, error) {
	fn := path.Join(dir, fileName)
	log.Infof("reading %s", fn)
	jsonABI, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, errors.Wrapf(err, "read abi json %s", fn)
	}

	fileDefs := make(map[string]abi.Event)
	// structs is the map of all redeclared structs shared by passed contracts.
	structs := make(map[string]*tmplStruct)

	unpackedABI, err := abi.JSON(bytes.NewReader(jsonABI))
	if err != nil {
		return nil, errors.Wrap(err, "unpack abi json %s")
	}

	for _, e := range unpackedABI.Events {
		if e.Anonymous {
			log.Trace("skipping anonymous event")
		}

		event := e
		for j, input := range event.Inputs {
			if input.Name == "" {
				event.Inputs[j].Name = fmt.Sprintf("arg%d", j)
			}
			if hasStruct(input.Type) {
				bindStructType(input.Type, structs)
			}
		}
		id := event.ID.String()
		ee, ok := fileDefs[id]
		if ok {
			log.Infof("we already have and event for %s, old: %s, new: %s", id, ee, event)
			continue
		}

		fileDefs[id] = event
	}
	if len(fileDefs) > 0 {
		return &ContractData{
			InputABI: embeddableABI(jsonABI),
			Defs:     fileDefs,
			Structs:  structs,
		}, nil
	}

	return nil, nil
}

func GenerateCode(prefix string, pkg string, contract ContractData) ([]byte, error) {
	buffer := new(bytes.Buffer)

	tmplFuncs := map[string]interface{}{
		"bindType":      bindType,
		"bindTopicType": bindTopicType,
		"namedType":     TypeName,
		"gopherize":     gopherize,
	}

	data := TemplateData{
		ContractData: contract,
		Package:      pkg,
		Prefix:       prefix,
	}
	tmpl := template.Must(template.New("").Funcs(tmplFuncs).Parse(eventsTemplate))
	err := tmpl.Execute(buffer, data)
	if err != nil {
		return nil, errors.Wrap(err, "execute template")
	}

	return format.Source(buffer.Bytes())
}
