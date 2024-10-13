package xsd

import (
	"fmt"
	"strings"
)

// ApplyFunctionP applies a function to the XSD and all children as long as function returns true
func (xsd *XSD) ApplyFunctionP(f func(XsdElement, interface{}) (interface{}, error)) (child interface{}, err error) {
	if xsd == nil {
		return
	}
	if child, err = f(xsd, nil); err != nil {
		return
	}
	if _, err = xsd.Import.applyFunctionP(f, child); err != nil {
		return
	}
	for _, ct := range xsd.ComplexTypes {
		if _, err = ct.applyFunctionP(f, child); err != nil {
			return
		}
	}
	for _, st := range xsd.SimpleTypes {
		if _, err = st.applyFunctionP(f, child); err != nil {
			return
		}
	}
	for _, elm := range xsd.Elements {
		if _, err = elm.applyFunctionP(f, child); err != nil {
			return
		}
	}

	return
}

// applyFunctionP applies a function to Import as long as function returns true
func (i *Import) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if i == nil {
		return true, nil
	}
	return f(i, parent)
}

// applyFunctionP applies a function to ComplexType and children as long as function returns true
func (ct *ComplexType) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if ct == nil {
		return true, nil
	}
	if child, err = f(ct, parent); err != nil {
		return
	}
	if _, err = ct.Sequence.applyFunctionP(f, child); err != nil {
		return
	}
	if _, err = ct.ComplexContent.applyFunctionP(f, child); err != nil {
		return
	}
	if _, err = ct.SimpleContent.applyFunctionP(f, child); err != nil {
		return
	}
	for _, a := range ct.Attributes {
		if _, err = a.applyFunctionP(f, child); err != nil {
			return
		}
	}
	if _, err = ct.Choice.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function to SimpleType and children as long as function returns true
func (st *SimpleType) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if st == nil {
		return true, nil
	}
	if child, err = f(st, parent); err != nil {
		return
	}
	if _, err = st.Restriction.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function Restriction and children as long as function returns true
func (r *Restriction) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if r == nil {
		return true, nil
	}
	if child, err = f(r, parent); err != nil {
		return
	}
	for _, e := range r.Enumerations {
		if _, err = e.applyFunctionP(f, child); err != nil {
			return
		}
	}
	if _, err = r.Pattern.applyFunctionP(f, child); err != nil {
		return
	}
	if _, err = r.MinInclusive.applyFunctionP(f, child); err != nil {
		return
	}
	if _, err = r.MaxInclusive.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function Pattern and children as long as function returns true
func (p *Pattern) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if p == nil {
		return true, nil
	}
	if child, err = f(p, parent); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function MinInclusive and children as long as function returns true
func (mi *MinInclusive) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if mi == nil {
		return true, nil
	}
	if child, err = f(mi, parent); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function MaxInclusive and children as long as function returns true
func (mi *MaxInclusive) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if mi == nil {
		return true, nil
	}
	if child, err = f(mi, parent); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function Enumeration and children as long as function returns true
func (e *Enumeration) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if e == nil {
		return true, nil
	}
	return f(e, parent)
}

// applyFunctionP applies a function to Sequence and children as long as function returns true
func (s *Sequence) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if s == nil {
		return true, nil
	}
	if child, err = f(s, parent); err != nil {
		return
	}
	for _, elm := range s.Elements {
		if _, err = elm.applyFunctionP(f, child); err != nil {
			return
		}
	}
	if _, err = s.Choice.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function to Choice and children as long as function returns true
func (c *Choice) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if c == nil {
		return true, nil
	}
	if child, err = f(c, parent); err != nil {
		return
	}
	for _, elm := range c.Elements {
		if _, err = elm.applyFunctionP(f, child); err != nil {
			return
		}
	}
	if _, err = c.Sequence.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function to Element and children as long as function returns true
func (e *Element) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if e == nil {
		return true, nil
	}
	if child, err = f(e, parent); err != nil {
		return
	}
	if _, err = e.ComplexType.applyFunctionP(f, child); err != nil {
		return
	}
	if _, err = e.SimpleType.applyFunctionP(f, child); err != nil {
		return
	}
	if _, err = e.Annotation.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function to ComplexContent and children as long as function returns true
func (cc *ComplexContent) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if cc == nil {
		return true, nil
	}
	if child, err = f(cc, parent); err != nil {
		return
	}
	if _, err = cc.Extension.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function to SimpleContent and children as long as function returns true
func (sc *SimpleContent) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if sc == nil {
		return true, nil
	}
	if child, err = f(sc, parent); err != nil {
		return
	}
	if _, err = sc.Extension.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function to Extension and children as long as function returns true
func (ex *Extension) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if ex == nil {
		return true, nil
	}
	if child, err = f(ex, parent); err != nil {
		return
	}
	if _, err = ex.Sequence.applyFunctionP(f, child); err != nil {
		return
	}
	for _, a := range ex.Attributes {
		if child, err = a.applyFunctionP(f, child); err != nil {
			return
		}
	}
	return
}

// applyFunctionP applies a function to Attribute as long as function returns true
func (a *Attribute) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if a == nil {
		return true, nil
	}
	if child, err = f(a, parent); err != nil {
		return
	}
	if _, err = a.ComplexType.applyFunctionP(f, child); err != nil {
		return
	}
	if _, err = a.SimpleType.applyFunctionP(f, child); err != nil {
		return
	}
	if _, err = a.Annotation.applyFunctionP(f, child); err != nil {
		return
	}
	return
}

// applyFunctionP applies a function to Annotation and children as long as function returns true
func (a *Annotation) applyFunctionP(f func(XsdElement, interface{}) (interface{}, error), parent interface{}) (child interface{}, err error) {
	if a == nil {
		return true, nil
	}
	return f(a, parent)
}

// ToStringP is a tree view with knowledge of parent
func (xsd *XSD) ToStringP() string {
	var sSlice []string
	fDisplay := func(xe XsdElement, h interface{}) (interface{}, error) {
		var level int
		if h != nil {
			level = h.(int) + 1
		}
		sSlice = append(sSlice, strings.Repeat("  ", level)+fmt.Sprintf("%d) %s", level, xe.ToString()))
		return level, nil
	}
	_, _ = xsd.ApplyFunctionP(fDisplay)
	return strings.Join(sSlice, "\n")
}
