package ethgen

var eventsTemplate = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lacasian/ethwheels/ethgen"
)

const {{.Prefix}}ABI = "{{.InputABI}}"

type {{.Prefix}}Decoder struct {
	*ethgen.Decoder
}

func New{{.Prefix}}Decoder() (*{{.Prefix}}Decoder, error) {
	dec, err := ethgen.NewDecoder({{.Prefix}}ABI)
	return &{{.Prefix}}Decoder {
		dec,
	}, err
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
type {{ $typeName }} struct {
	{{- range .Inputs }}
	{{ gopherize .Name }} {{ if .Indexed }}{{ bindTopicType .Type $.Structs }}{{ else }}{{ bindType .Type $.Structs }}{{ end }}
	{{- end }}
	Raw types.Log
}

func (d *{{$.Prefix}}Decoder) {{ $typeName }}ID() common.Hash {
	return common.HexToHash("{{ .ID }}")
}

func (d *{{$.Prefix}}Decoder) Is{{ $typeName }}(log *types.Log) bool {
	if len(log.Topics) == 0 {
		return false
	}
	return log.Topics[0] == d.{{ $typeName }}ID()
}

func (d *{{$.Prefix}}Decoder) {{ $typeName }}(log types.Log) ({{ $typeName }}, error) {
	var out {{ $typeName }}
	if !d.Is{{ $typeName }}(&log) {
		return out, ethgen.ErrMismatchingEvent
	}
	err := d.UnpackLog(&out, "{{ $event.Name }}", log)
	out.Raw = log
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
