package ethgen

var eventsTemplate = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"math/big"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lacasian/ethwheels/ethgen"
)

// Reference imports to suppress errors
var (
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = web3types.Log{}
)

const {{.Prefix}}ABI = "{{.InputABI}}"

var {{.Prefix}} = New{{.Prefix}}Decoder()

type {{.Prefix}}Decoder struct {
	*ethgen.Decoder
}

func New{{.Prefix}}Decoder() *{{.Prefix}}Decoder {
	dec := ethgen.NewDecoder({{.Prefix}}ABI)
	return &{{.Prefix}}Decoder {
		dec,
	}
}

{{range .Structs}}
	type {{.Name}} struct {
	{{range $field := .Fields}}
	{{$field.Name}} {{$field.Type}}{{end}}
	}
{{end}}

{{ range $key, $event := .Defs }}
{{ $typePrefix := namedType $.Prefix $event.Name }}
{{ $typeName := (printf "%s%s" $typePrefix "Event") }}
{{ $typeShortName := (printf "%s%s" $event.Name "Event") }}
type {{ $typeName }} struct {
	{{- range .Inputs }}
	{{ gopherize .Name }} {{ if .Indexed }}{{ bindTopicType .Type $.Structs }}{{ else }}{{ bindType .Type $.Structs }}{{ end }}
	{{- end }}
	Raw types.Log
}

func (d *{{$.Prefix}}Decoder) {{ $typeShortName }}ID() common.Hash {
	return common.HexToHash("{{ .ID }}")
}

func (d *{{$.Prefix}}Decoder) Is{{ $typeShortName }}(log *types.Log) bool {
	if len(log.Topics) == 0 {
		return false
	}
	return log.Topics[0] == d.{{ $typeShortName }}ID()
}

func (d *{{$.Prefix}}Decoder) Is{{ $typeShortName }}W3(log *web3types.Log) bool {
	if len(log.Topics) == 0 {
		return false
	}
	return log.Topics[0] == d.{{ $typeShortName }}ID().String()
}

func (d *{{$.Prefix}}Decoder) {{ $typeShortName }}W3(w3l web3types.Log) ({{ $typeName }}, error) {
	l, err := ethgen.W3LogToLog(w3l)
	if err != nil {
		return  {{ $typeName }}{}, err
	}

	return d.{{ $typeShortName }}(l)
}

func (d *{{$.Prefix}}Decoder) {{ $typeShortName }}(l types.Log) ({{ $typeName }}, error) {
	var out {{ $typeName }}
	if !d.Is{{ $typeShortName }}(&l) {
		return out, ethgen.ErrMismatchingEvent
	}
	err := d.UnpackLog(&out, "{{ $event.Name }}", l)
	out.Raw = l
	return out, err
}

{{ end }}
`

/*
"strings"

ethereum "github.com/ethereum/go-ethereum"
"github.com/ethereum/go-ethereum/accounts/abi"

"github.com/ethereum/go-ethereum/accounts/abi/bind"
"github.com/ethereum/go-ethereum/event"

// Reference imports to suppress errors if they are not otherwise used.
// var (
// 	_ = big.NewInt
// 	_ = strings.NewReader
// 	_ = ethereum.NotFound
// 	_ = bind.Bind
// 	_ = common.Big1
// 	_ = types.BloomLookup
// 	_ = event.NewSubscription
// )

// {{$contract.Type}}{{.Normalized.Name}} represents a {{.Normalized.Name}} event raised by the {{$contract.Type}} contract.
type {{$contract.Type}}{{.Normalized.Name}} struct { {{range .Normalized.Inputs}}
{{capitalise .Name}} {{if .Indexed}}{{bindtopictype .Type $structs}}{{else}}{{bindtype .Type $structs}}{{end}}; {{end}}
Raw types.Log // Blockchain specific contextual infos
}
*/
