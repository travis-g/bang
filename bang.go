package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"text/template"
)

// SliceToMap returns a map of Bangs based on the input slice's Bangs' names.
func SliceToMap(slice []Bang) (bangs map[string]Bang) {
	bangs = map[string]Bang{}
	for _, bang := range slice {
		bangs[bang.Name] = bang
	}
	return
}

// MapToNamedBangs returns a slice of Bangs with their names set according to
// keys from an input map.
func MapToNamedBangs(bangs map[string]Bang) (slice []Bang) {
	slice = []Bang{}
	for key, bang := range bangs {
		bang.Name = key
		slice = append(slice, bang)
	}
	return
}

const symbol = "{{{s}}}"
const bangTemplate = `{{.Name}} - {{.Description}}
`

// URL returns the direct query URL for a Bang.
func (b *Bang) URL(q string) string {
	fmt.Println(b)
	b.EscapeMethod = Bang_PATH_ESCAPE
	str, _ := json.Marshal(b)
	fmt.Println(string(str))
	var s string
	switch b.GetEscapeMethod() {
	case Bang_QUERY_ESCAPE:
		s = url.QueryEscape(q)
	case Bang_PATH_ESCAPE:
		s = url.PathEscape(q)
	case Bang_PASS_THROUGH:
		s = q
	}
	return fmt.Sprint(strings.Replace(b.GetFormat(), symbol, s, 1))
}

// Bangs is the registration list of all Bangs. We use a map[string] here for
// faster lookups at runtime than digging through slices.
var Bangs = map[string]Bang{}

func listBangs(bangs map[string]Bang) string {
	list := MapToNamedBangs(bangs)
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	tmpl, _ := template.New("output").Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
	}).Parse(bangTemplate)
	var buf bytes.Buffer
	for _, bang := range list {
		tmpl.Execute(&buf, bang)
	}
	return strings.TrimSpace(buf.String())
}

// Filter filters a map of Bangs based on a predicate.
//
//    fmt.Println(Filter(Bangs, func(v string) bool {
//      return strings.Contains(strings.ToLower(v), strings.ToLower("Bing"))
//    }))
func Filter(vs map[string]Bang, f func(string) bool) map[string]Bang {
	vsf := make(map[string]Bang, 0)
	for i, v := range vs {
		if f(i) || f(v.Description) {
			vsf[i] = v
		}
	}
	return vsf
}
