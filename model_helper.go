package xsd

import (
	"encoding/xml"
	"fmt"
)

func NewXSD(xsdXML []byte) (xsd *XSD, err error) {
	if err = xml.Unmarshal(xsdXML, &xsd); err != nil {
		return nil, fmt.Errorf("could not unmarshal XML into XSD, got %v", err)
	}
	return
}

// ApplyFunction applies a function to the XSD and all children as long as function returns true
func (xsd *XSD) ApplyFunction(f func(XsdElement) error) (err error) {
	if xsd == nil {
		return
	}
	if err = f(xsd); err != nil {
		return
	}
	//if err = f(xsd.XMLName); err != nil {
	//	return
	//}
	if err = xsd.Import.applyFunction(f); err != nil {
		return
	}
	for _, ct := range xsd.ComplexTypes {
		if err = ct.applyFunction(f); err != nil {
			return
		}
	}
	for _, st := range xsd.SimpleTypes {
		if err = st.applyFunction(f); err != nil {
			return
		}
	}
	for _, e := range xsd.Elements {
		if err = e.applyFunction(f); err != nil {
			return
		}
	}
	return
}

func (i *Import) applyFunction(f func(XsdElement) error) (err error) {
	if i == nil {
		return nil
	}
	return f(i)
}

// applyFunction applies a function ComplexType and children as long as function returns true
func (ct *ComplexType) applyFunction(f func(XsdElement) error) (err error) {
	if ct == nil {
		return nil
	}
	if err = f(ct); err != nil {
		return
	}
	if err = ct.Sequence.applyFunction(f); err != nil {
		return
	}
	if err = ct.ComplexContent.applyFunction(f); err != nil {
		return
	}
	if err = ct.SimpleContent.applyFunction(f); err != nil {
		return
	}
	for _, a := range ct.Attributes {
		if err = a.applyFunction(f); err != nil {
			return
		}
	}
	if err = ct.Choice.applyFunction(f); err != nil {
		return
	}
	return
}

// applyFunction applies a function SimpleType and children as long as function returns true
func (st *SimpleType) applyFunction(f func(XsdElement) error) (err error) {
	if st == nil {
		return nil
	}
	if err = f(st); err != nil {
		return
	}
	if err = st.Restriction.applyFunction(f); err != nil {
		return
	}
	return
}

// applyFunction applies a function Restriction and children as long as function returns true
func (r *Restriction) applyFunction(f func(XsdElement) error) (err error) {
	if r == nil {
		return nil
	}
	if err = f(r); err != nil {
		return
	}
	for _, e := range r.Enumerations {
		if err = e.applyFunction(f); err != nil {
			return
		}
	}
	return
}

// applyFunction applies a function Enumeration and children as long as function returns true
func (e *Enumeration) applyFunction(f func(XsdElement) error) (err error) {
	if e == nil {
		return nil
	}
	if err = f(e); err != nil {
		return
	}
	return
}

func (s *Sequence) applyFunction(f func(XsdElement) error) (err error) {
	if s == nil {
		return nil
	}
	if err = f(s); err != nil {
		return
	}
	for _, e := range s.Elements {
		if err = e.applyFunction(f); err != nil {
			return
		}
	}
	if err = s.Choice.applyFunction(f); err != nil {
		return
	}
	return
}

func (e *Element) applyFunction(f func(XsdElement) error) (err error) {
	if e == nil {
		return nil
	}
	if err = f(e); err != nil {
		return
	}
	if err = e.ComplexType.applyFunction(f); err != nil {
		return
	}
	if err = e.SimpleType.applyFunction(f); err != nil {
		return
	}
	if err = e.Annotation.applyFunction(f); err != nil {
		return
	}
	return
}

func (cc *ComplexContent) applyFunction(f func(XsdElement) error) (err error) {
	if cc == nil {
		return nil
	}
	if err = f(cc); err != nil {
		return
	}
	err = cc.Extension.applyFunction(f)
	return
}

func (sc *SimpleContent) applyFunction(f func(XsdElement) error) (err error) {
	if sc == nil {
		return nil
	}
	if err = f(sc); err != nil {
		return
	}
	err = sc.Extension.applyFunction(f)
	return
}

func (ex *Extension) applyFunction(f func(XsdElement) error) (err error) {
	if ex == nil {
		return nil
	}
	if err = f(ex); err != nil {
		return
	}
	if err = ex.Sequence.applyFunction(f); err != nil {
		return
	}
	for _, a := range ex.Attributes {
		if err = a.applyFunction(f); err != nil {
			return
		}
	}
	return
}

func (a *Attribute) applyFunction(f func(XsdElement) error) (err error) {
	if a == nil {
		return nil
	}
	if err = f(a); err != nil {
		return
	}
	if err = a.ComplexType.applyFunction(f); err != nil {
		return
	}
	if err = a.SimpleType.applyFunction(f); err != nil {
		return
	}
	if err = a.Annotation.applyFunction(f); err != nil {
		return
	}
	return
}

func (ch *Choice) applyFunction(f func(XsdElement) error) (err error) {
	if ch == nil {
		return nil
	}
	if err = f(ch); err != nil {
		return
	}
	if err = ch.Sequence.applyFunction(f); err != nil {
		return
	}
	for _, e := range ch.Elements {
		if err = e.applyFunction(f); err != nil {
			return
		}
	}
	return
}

func (a *Annotation) applyFunction(f func(XsdElement) error) (err error) {
	if a == nil {
		return nil
	}
	if err = f(a); err != nil {
		return
	}
	return
}

func (xsd *XSD) ItemsString() (sSlice []string) {
	fDisplay := func(xe XsdElement) error {
		sSlice = append(sSlice, xe.ToString())
		return nil
	}
	_ = xsd.ApplyFunction(fDisplay)
	return
}
