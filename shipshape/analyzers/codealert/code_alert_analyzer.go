package codealert

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"code.google.com/p/goprotobuf/proto"
	notepb "shipshape/proto/note_proto"
	ctxpb "shipshape/proto/shipshape_context_proto"
)

var (
	// TODO(emso): Consider how to enable/disable code alerts
	alerts = []*CodeAlert{{
		Name:        "DoNotSubmitTxtTest",
		File:        "*",
		Description: "Do not submit test text check",
		Regexp:      regexp.MustCompile(".*do not submit.*"),
	}}
)

// A CodeAlert represents a regexp being matched on files included for the alert.
type CodeAlert struct {
	Name        string
	File        string
	Description string
	Regexp      *regexp.Regexp
}

type CodeAlertAnalyzer struct {
}

func (CodeAlertAnalyzer) Category() string { return "CodeAlert" }

// TODO(emso): Use file filter in code alert
func (a CodeAlertAnalyzer) Analyze(ctx *ctxpb.ShipshapeContext) ([]*notepb.Note, error) {
	var notes []*notepb.Note
	for _, path := range ctx.FilePath {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		notes = append(notes, a.FindMatches(string(content))...)
	}
	return notes, nil
}

// FindMatches returns an array of notes for each match in content.
func (a CodeAlertAnalyzer) FindMatches(content string) []*notepb.Note {
	var notes []*notepb.Note
	// Line number will start from zero and should be padded with a one if returned
	for lineNumber, line := range strings.Split(content, "\n") {
		for _, alert := range alerts {
			match := alert.Regexp.FindString(line)
			if match == "" {
				continue
			}
			log.Printf("Found match (%v) on line %v for %v code alert", match, lineNumber+1, alert.Name)
			// TODO(emso): Add location to note (filename, lineNumber + 1)
			notes = append(notes, &notepb.Note{
				Category:    proto.String(a.Category()),
				Subcategory: proto.String(alert.Name),
				Description: proto.String(alert.Description),
			})
		}
	}
	return notes
}
