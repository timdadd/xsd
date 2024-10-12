package xsd

import (
	"encoding/xml"
	"fmt"
)

type XsdElement interface {
	ToString() string
}

type XSD struct {
	XMLName      xml.Name       `xml:"schema"`
	Import       *Import        `xml:"import,omitempty"`
	SimpleTypes  []*SimpleType  `xml:"simpleType,omitempty"`
	ComplexTypes []*ComplexType `xml:"complexType,omitempty"`
}

type Import struct {
	Namespace      string `xml:"namespace,attr"`
	SchemaLocation string `xml:"schemaLocation,attr"`
}

type ComplexType struct {
	Name           string          `xml:"name,attr"`
	Sequence       *Sequence       `xml:"sequence,omitempty"`
	ComplexContent *ComplexContent `xml:"complexContent,omitempty"`
	SimpleContent  *SimpleContent  `xml:"simpleContent,omitempty"`
	Choice         *Choice         `xml:"choice,omitempty"`
}

type SimpleType struct {
	Name        string       `xml:"name,attr"`
	Restriction *Restriction `xml:"restriction,omitempty"`
}

type Enumeration struct {
	Value string `xml:"value,attr"`
}

type Restriction struct {
	Base         string         `xml:"base,attr"`
	Enumerations []*Enumeration `xml:"enumeration,omitempty"`
	MinInclusive *MinInclusive `xml:"minInclusive,omitempty"`
	MaxInclusive *MaxInclusive `xml:"maxInclusive,omitempty"`
	Pattern *Pattern `xml:"pattern,omitempty"`
}

type MinInclusive struct {
	Value string `xml:"value,attr"`
}

type MaxInclusive struct {
	Value string `xml:"value,attr"`
}


type Pattern struct {
	Value string `xml:"value,attr"`
}

type Sequence struct {
	Name      string     `xml:"name,attr"`
	MinOccurs string     `xml:"minOccurs,attr"`
	MaxOccurs string     `xml:"maxOccurs,attr"`
	Elements  []*Element `xml:"element,omitempty"`
	Choice    *Choice    `xml:"choice,omitempty"`
}

type Element struct {
	Name        string       `xml:"name,attr"`
	Type        string       `xml:"type,attr"`
	MinOccurs   string       `xml:"minOccurs,attr"`
	MaxOccurs   string       `xml:"maxOccurs,attr"`
	ComplexType *ComplexType `xml:"complexType,omitempty"`
	SimpleType  *SimpleType  `xml:"simpleType,omitempty"`
	Annotation  *Annotation  `xml:"annotation,omitempty"`
}

type Annotation struct {
	Documentation string `xml:"documentation"`
}

type Choice struct {
	Name      string     `xml:"name,attr"`
	MinOccurs string     `xml:"minOccurs,attr"`
	MaxOccurs string     `xml:"maxOccurs,attr"`
	Elements  []*Element `xml:"element,omitempty"`
	Sequence  *Sequence  `xml:"sequence,omitempty"`
}

type SimpleContent struct {
	Extension *Extension `xml:"extension,omitempty"`
}

type ComplexContent struct {
	Extension *Extension `xml:"extension,omitempty"`
}

// This
type Extension struct {
	ID     string    `xml:"id,attr"`
	Base     string    `xml:"base,attr"`
	Sequence *Sequence `xml:"sequence,omitempty"`
	Attributes []*Attribute `xml:"attribute,omitempty"`
}

type Attribute struct {
	Name      string     `xml:"name,attr"`

}

func (xsd *XSD) ToString() string {
	return "XSD"
}

func (i *Import) ToString() string {
	return fmt.Sprintf("Import: Namespace:%s, Location:%s", i.Namespace, i.SchemaLocation)
}

func (ct *ComplexType) ToString() string {
	return fmt.Sprintf("ComplexType: %s", ct.Name)
}
func (st *SimpleType) ToString() string {
	return fmt.Sprintf("SimpleType: %s", st.Name)
}
func (r *Restriction) ToString() string {
	return fmt.Sprintf("Base: %s", r.Base)
}
func (s *Sequence) ToString() string {
	return fmt.Sprintf("Sequence: %s %s", s.Name, occurs(s.MinOccurs, s.MaxOccurs))
}
func (e *Element) ToString() string {
	return fmt.Sprintf("Element: %s (%s) %s", e.Name, e.Type, occurs(e.MinOccurs, e.MaxOccurs))
}
func (ch *Choice) ToString() string {
	return fmt.Sprintf("Choice: %s %s", ch.Name, occurs(ch.MinOccurs, ch.MaxOccurs))
}
func (a *Annotation) ToString() string {
	return fmt.Sprintf("Annotation: %s", a.Documentation)
}
func (cc *ComplexContent) ToString() string {
	return "ComplexContent"
}
func (ex *Extension) ToString() string {
	return fmt.Sprintf("Extension: Base:%s", ex.Base)
}

func occurs(minOccurs, maxOccurs string) string {
	if minOccurs == "" && maxOccurs == "" {
		return ""
	}
	if minOccurs == "0" && maxOccurs == "" {
		return "optional"
	}
	if minOccurs == "1" && maxOccurs == "1" {
		return "1:1"
	}
	if minOccurs == "0" && maxOccurs == "unbounded" {
		return "0:N"
	}
	if minOccurs == "1" && maxOccurs == "unbounded" {
		return "1:N"
	}
	return fmt.Sprintf("%s:%s", minOccurs, maxOccurs)
}
