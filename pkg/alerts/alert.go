package alerts

import (
	"time"

	"github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"
)

type Name string

func (n Name) String() string {
	return string(n)
}

type Code string

func (n Code) String() string {
	return string(n)
}

type Alert interface {
	GetDate() time.Time
	GetSeverity() Severity
	GetName() Name
	GetCode() Code
	GetSummary() string
	GetDescription() string
	GetAdditionalLabels() collections.BasicMapCollection[string]
	GetAdditionalAnnotations() collections.BasicMapCollection[string]
	AddAdditionalLabels(label string, description string) Alert
	AddAdditionalAnnotations(annotation string, description string) Alert
	Clone() Alert
}

func NewAlert(
	severity Severity,
	name Name,
	summary string,
	description string,
) Alert {
	return &BasicAlert{
		date:        time.Now().UTC(),
		name:        name,
		summary:     summary,
		description: description,
	}
}

type BasicAlert struct {
	date                  time.Time
	severity              Severity
	name                  Name
	code                  Code
	summary               string
	description           string
	additionalLabels      collections.BasicMapCollection[string]
	additionalAnnotations collections.BasicMapCollection[string]
}

func (a *BasicAlert) GetDate() time.Time {
	return a.date
}

func (a *BasicAlert) GetSeverity() Severity {
	return a.severity
}

func (a *BasicAlert) GetName() Name {
	return a.name
}

func (a *BasicAlert) GetCode() Code {
	return a.code
}

func (a *BasicAlert) GetSummary() string {
	return a.summary
}

func (a *BasicAlert) GetDescription() string {
	return a.description
}

func (a *BasicAlert) GetAdditionalLabels() collections.BasicMapCollection[string] {
	return a.additionalLabels
}

func (a *BasicAlert) GetAdditionalAnnotations() collections.BasicMapCollection[string] {
	return a.additionalAnnotations
}

func (a *BasicAlert) SetCode(code Code) Alert {
	a.code = code
	return a
}

func (a *BasicAlert) AddAdditionalLabels(label string, description string) Alert {
	a.additionalLabels.Set(label, description)
	return a
}

func (a *BasicAlert) AddAdditionalAnnotations(annotation string, description string) Alert {
	a.additionalAnnotations.Set(annotation, description)
	return a
}

func (a *BasicAlert) Clone() Alert {
	b := &BasicAlert{
		date:        time.Now().UTC(),
		name:        a.name,
		code:        a.code,
		summary:     a.summary,
		description: a.description,
	}

	for k, v := range a.additionalAnnotations {
		b.AddAdditionalAnnotations(k, v)
	}

	for k, v := range a.additionalLabels {
		b.AddAdditionalLabels(k, v)
	}

	return b
}
