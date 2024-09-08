package xsd

import (
	"fmt"
	"strings"
)

// ApplyFunctionP applies a function to the XSD and all children as long as function returns true
func (xsd *XSD) ApplyFunctionP(f func(xe XsdElement, p interface{}) (bool, interface{})) (x bool, c interface{}) {
	if xsd == nil {
		return
	}
	if x, c = f(xsd, nil); !x {
		return
	}
	//if x, _ = f(xsd.XMLName, c); !x {
	//	return
	//}
	if x, _ = f(xsd.Import, c); !x {
		return
	}
	for _, complexType := range xsd.ComplexTypes {
		if x, _ = complexType.applyFunctionP(f, c); !x {
			return
		}
	}
	return
}

// applyFunctionP applies a function to ComplexType and children as long as function returns true
func (ct *ComplexType) applyFunctionP(f func(xe XsdElement, p interface{}) (bool, interface{}), p interface{}) (x bool, c interface{}) {
	if ct == nil {
		return true, nil
	}
	if x, c = f(ct, p); !x {
		return
	}
	if x, _ = ct.Sequence.applyFunctionP(f, c); !x {
		return
	}
	if x, _ = ct.ComplexContent.applyFunctionP(f, c); !x {
		return
	}
	if x, _ = ct.Choice.applyFunctionP(f, c); !x {
		return
	}
	return
}

// applyFunctionP applies a function to Sequence and children as long as function returns true
func (s *Sequence) applyFunctionP(f func(xe XsdElement, p interface{}) (bool, interface{}), p interface{}) (x bool, c interface{}) {
	if s == nil {
		return true, nil
	}
	if x, c = f(s, p); !x {
		return
	}
	for _, e := range s.Elements {
		if x, _ = e.applyFunctionP(f, c); !x {
			return
		}
	}
	if x, _ = s.Choice.applyFunctionP(f, c); !x {
		return
	}
	return
}

// applyFunctionP applies a function to Choice and children as long as function returns true
func (ch *Choice) applyFunctionP(f func(xe XsdElement, p interface{}) (bool, interface{}), p interface{}) (x bool, c interface{}) {
	if ch == nil {
		return true, nil
	}
	if x, c = f(ch, p); !x {
		return
	}
	for _, e := range ch.Elements {
		if x, _ = e.applyFunctionP(f, c); !x {
			return
		}
	}
	if x, _ = ch.Sequence.applyFunctionP(f, c); !x {
		return
	}
	return
}

// applyFunctionP applies a function to Element and children as long as function returns true
func (e *Element) applyFunctionP(f func(xe XsdElement, p interface{}) (bool, interface{}), p interface{}) (x bool, c interface{}) {
	if e == nil {
		return true, nil
	}
	if x, c = f(e, p); !x {
		return
	}
	if x, _ = e.ComplexType.applyFunctionP(f, c); !x {
		return
	}
	if x, _ = e.Annotation.applyFunctionP(f, c); !x {
		return
	}
	return
}

// applyFunctionP applies a function to ComplexContent and children as long as function returns true
func (cc *ComplexContent) applyFunctionP(f func(xe XsdElement, p interface{}) (bool, interface{}), p interface{}) (x bool, c interface{}) {
	if cc == nil {
		return true, nil
	}
	if x, c = f(cc, p); !x {
		return
	}
	if x, _ = cc.Extension.applyFunctionP(f, c); !x {
		return
	}
	return
}

// applyFunctionP applies a function to Extension and children as long as function returns true
func (ex *Extension) applyFunctionP(f func(xe XsdElement, p interface{}) (bool, interface{}), p interface{}) (x bool, c interface{}) {
	if ex == nil {
		return true, nil
	}
	if x, c = f(ex, p); !x {
		return
	}
	if x, _ = ex.Sequence.applyFunctionP(f, c); !x {
		return
	}
	return
}

// applyFunctionP applies a function to Extension and children as long as function returns true
func (a *Annotation) applyFunctionP(f func(xe XsdElement, p interface{}) (bool, interface{}), p interface{}) (x bool, c interface{}) {
	if a == nil {
		return true, nil
	}
	if x, c = f(a, p); !x {
		return
	}
	return
}

// ToStringP is a tree view with knowledge of parent
func (xsd *XSD) ToStringP() string {
	var sSlice []string
	fDisplay := func(xe XsdElement, h interface{}) (bool, interface{}) {
		var level int
		if h != nil {
			level = h.(int) + 1
		}
		sSlice = append(sSlice, strings.Repeat("  ", level)+fmt.Sprintf("%d) %s", level, xe.ToString()))
		return true, level
	}
	xsd.ApplyFunctionP(fDisplay)
	return strings.Join(sSlice, "\n")
}
