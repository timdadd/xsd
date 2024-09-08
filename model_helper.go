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
func (xsd *XSD) ApplyFunction(f func(xe XsdElement) bool) (x bool) {
	if xsd == nil {
		return
	}
	if x = f(xsd); !x {
		return
	}
	//if x = f(xsd.XMLName); !x {
	//	return
	//}
	if x = f(xsd.Import); !x {
		return
	}
	for _, complexType := range xsd.ComplexTypes {
		if x = complexType.applyFunction(f); !x {
			return
		}
	}
	return
}

// applyFunction applies a function ComplexType and children as long as function returns true
func (ct *ComplexType) applyFunction(f func(xe XsdElement) bool) (x bool) {
	if ct == nil {
		return true
	}
	if x = f(ct); !x {
		return
	}
	if x = ct.Sequence.applyFunction(f); !x {
		return
	}
	if x = ct.ComplexContent.applyFunction(f); !x {
		return
	}
	if x = ct.Choice.applyFunction(f); !x {
		return
	}
	return
}

func (s *Sequence) applyFunction(f func(xe XsdElement) bool) (x bool) {
	if s == nil {
		return true
	}
	if x = f(s); !x {
		return
	}
	for _, e := range s.Elements {
		if x = e.applyFunction(f); !x {
			return
		}
	}
	if x = s.Choice.applyFunction(f); !x {
		return
	}
	return
}

func (e *Element) applyFunction(f func(xe XsdElement) bool) (x bool) {
	if e == nil {
		return true
	}
	if x = f(e); !x {
		return
	}
	if x = e.ComplexType.applyFunction(f); !x {
		return
	}
	if x = e.Annotation.applyFunction(f); !x {
		return
	}
	return
}

func (cc *ComplexContent) applyFunction(f func(xe XsdElement) bool) (x bool) {
	if cc == nil {
		return true
	}
	if x = f(cc); !x {
		return
	}
	x = cc.Extension.applyFunction(f)
	return
}

func (ex *Extension) applyFunction(f func(xe XsdElement) bool) (x bool) {
	if ex == nil {
		return true
	}
	if x = f(ex); !x {
		return
	}
	x = ex.Sequence.applyFunction(f)
	return
}

func (ch *Choice) applyFunction(f func(xe XsdElement) bool) (x bool) {
	if ch == nil {
		return true
	}
	if x = f(ch); !x {
		return
	}
	x = ch.Sequence.applyFunction(f)
	for _, e := range ch.Elements {
		if x = e.applyFunction(f); !x {
			return
		}
	}
	return
}

func (a *Annotation) applyFunction(f func(xe XsdElement) bool) (x bool) {
	if a == nil {
		return true
	}
	if x = f(a); !x {
		return
	}
	return
}

func (xsd *XSD) ItemsString() (sSlice []string) {
	fDisplay := func(xe XsdElement) bool {
		sSlice = append(sSlice, xe.ToString())
		return true
	}
	xsd.ApplyFunction(fDisplay)
	return
}
