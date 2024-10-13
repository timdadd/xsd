package xsd

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type XsdElement interface {
	ToString() string
}

type XsdItem interface {
	XsdElement
	IsMandatoryOptional() string
	IsRepeated() bool
}

type XSD struct {
	XMLName      xml.Name       `xml:"schema"`
	Import       *Import        `xml:"import,omitempty"`
	SimpleTypes  []*SimpleType  `xml:"simpleType,omitempty"`
	ComplexTypes []*ComplexType `xml:"complexType,omitempty"`
	Elements     []*Element     `xml:"element,omitempty"`
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
	Attributes     []*Attribute    `xml:"attribute,omitempty"`
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
	MinInclusive *MinInclusive  `xml:"minInclusive,omitempty"`
	MaxInclusive *MaxInclusive  `xml:"maxInclusive,omitempty"`
	Pattern      *Pattern       `xml:"pattern,omitempty"`
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
	Ref         string       `xml:"ref,attr"`
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
	ID         string       `xml:"id,attr"`
	Base       string       `xml:"base,attr"`
	Sequence   *Sequence    `xml:"sequence,omitempty"`
	Attributes []*Attribute `xml:"attribute,omitempty"`
}

type Attribute struct {
	Name        string       `xml:"name,attr"`
	Type        string       `xml:"type,attr"`
	Use         string       `xml:"use,attr,omitempty"`
	Ref         string       `xml:"ref,attr"`
	ComplexType *ComplexType `xml:"complexType,omitempty"`
	SimpleType  *SimpleType  `xml:"simpleType,omitempty"`
	Annotation  *Annotation  `xml:"annotation,omitempty"`
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
	return fmt.Sprintf("Restriction: %s", r.Base)
}
func (e *Enumeration) ToString() string {
	return fmt.Sprintf("Enumeration: %s", e.Value)
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
func (sc *SimpleContent) ToString() string {
	return "SimpleContent"
}
func (ex *Extension) ToString() string {
	return fmt.Sprintf("Extension: Base:%s", ex.Base)
}
func (a *Attribute) ToString() string {
	return fmt.Sprintf("Attribute: %s", a.Name)
}
func (p *Pattern) ToString() string       { return fmt.Sprintf("Pattern: %s", p.Value) }
func (mi *MinInclusive) ToString() string { return fmt.Sprintf("MinInclusive: %s", mi.Value) }
func (mi *MaxInclusive) ToString() string { return fmt.Sprintf("MaxInclusive: %s", mi.Value) }

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

// xsdItem methods
func (a *Attribute) IsMandatoryOptional() string {
	if a.Use == "required" {
		return "M"
	}
	return ""
}

func (a *Attribute) IsRepeated() bool {
	return false
}

func (e *Element) IsMandatoryOptional() string {
	//if e.Use == "required" {return "M"}
	if e.MinOccurs == "0" {
		return "O"
	}
	if e.MinOccurs > "" {
		return "M"
	}
	return ""
}

func (e *Element) IsRepeated() bool {
	if e.MaxOccurs == "" {
		return false
	}
	if e.MaxOccurs == "unbounded" {
		return true
	}
	if i, err := strconv.Atoi(e.MaxOccurs); err == nil && i > 1 {
		return true
	}
	return false
}
