package main

import (
	"io"
	"log"
	"net/url"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/hairyhenderson/gomplate/aws"
)

func (g *Gomplate) createChildTemplate(filename string) *template.Template {
	if filename == "" {
		filename = "anonymous_template"
	}
	return g.parent.New(filename)
}

// Gomplate -
type Gomplate struct {
	parent *template.Template
}

// RunTemplate -
func (g *Gomplate) RunTemplate(input Input, out io.Writer) {
	context := &Context{}
	tmpl, err := g.createChildTemplate(input.filename).Parse(input.text)
	if err != nil {
		log.Fatalf("Line %q: %v\n", input.text, err)
	}

	if input.partial {
		return // Do not execute partials
	}

	if err := tmpl.Execute(out, context); err != nil {
		panic(err)
	}
}

// NewGomplate -
func NewGomplate(data *Data, leftDelim, rightDelim string) *Gomplate {
	env := &Env{}
	typeconv := &TypeConv{}
	stringfunc := &stringFunc{}
	ec2meta := aws.NewEc2Meta()
	ec2info := aws.NewEc2Info()

	funcMap := template.FuncMap{
		"getenv":           env.Getenv,
		"fromStrings":      typeconv.fromStrings,
		"bool":             typeconv.Bool,
		"has":              typeconv.Has,
		"json":             typeconv.JSON,
		"jsonArray":        typeconv.JSONArray,
		"yaml":             typeconv.YAML,
		"yamlArray":        typeconv.YAMLArray,
		"toml":             typeconv.TOML,
		"csv":              typeconv.CSV,
		"csvByRow":         typeconv.CSVByRow,
		"csvByColumn":      typeconv.CSVByColumn,
		"slice":            typeconv.Slice,
		"indent":           typeconv.indent,
		"join":             typeconv.Join,
		"toJSON":           typeconv.ToJSON,
		"toJSONPretty":     typeconv.toJSONPretty,
		"toYAML":           typeconv.ToYAML,
		"toTOML":           typeconv.ToTOML,
		"toCSV":            typeconv.ToCSV,
		"ec2meta":          ec2meta.Meta,
		"ec2dynamic":       ec2meta.Dynamic,
		"ec2tag":           ec2info.Tag,
		"ec2region":        ec2meta.Region,
		"contains":         strings.Contains,
		"hasPrefix":        strings.HasPrefix,
		"hasSuffix":        strings.HasSuffix,
		"replaceAll":       stringfunc.replaceAll,
		"split":            strings.Split,
		"splitN":           strings.SplitN,
		"title":            strings.Title,
		"toUpper":          strings.ToUpper,
		"toLower":          strings.ToLower,
		"trim":             strings.Trim,
		"trimSpace":        strings.TrimSpace,
		"urlParse":         url.Parse,
		"datasource":       data.Datasource,
		"ds":               data.Datasource,
		"datasourceExists": data.DatasourceExists,
		"include":          data.include,
	}
	parentTpl := template.New("parent").
		Funcs(funcMap).Funcs(sprig.TxtFuncMap()).Delims(leftDelim, rightDelim)

	return &Gomplate{
		parent: parentTpl,
	}
}

func runTemplate(o *GomplateOpts) error {
	defer runCleanupHooks()
	data := NewData(o.dataSources, o.dataSourceHeaders)

	g := NewGomplate(data, o.lDelim, o.rDelim)

	if o.inputDir != "" {
		return processInputDir(o.inputDir, o.outputDir, g)
	}

	return processInputFiles(o.input, o.inputFiles, o.outputFiles, g)
}

// Called from process.go ...
func renderTemplate(g *Gomplate, input Input, outPath string) error {
	outFile, err := openOutFile(outPath)
	if err != nil {
		return err
	}
	// nolint: errcheck
	defer outFile.Close()
	g.RunTemplate(input, outFile)
	return nil
}
