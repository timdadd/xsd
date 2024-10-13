package xsd

import (
	"fmt"
	"sort"
	"strings"
)

type Message struct {
	Package       string         `json:"package,omitempty"`
	Name          string         `json:"name,omitempty"`
	MessageItems  []*MessageItem `json:"messageItems,omitempty"`
	Description   string         `json:"description,omitempty"`
	IsRootMessage bool           `json:"isNamed,omitempty"` // If set to true then this is a root level message and not a sub message
	sequence      int
}

type MessageItem struct {
	Name              string   `json:"name,omitempty"`
	Type              string   `json:"type,omitempty"`
	Format            string   `json:"format,omitempty"`
	Repeated          bool     `json:"repeated,omitempty"`
	MandatoryOptional string   `json:"mandatoryOptional,omitempty"` // M=Mandatory,O=Optional, blank=don't know
	MinOccurs         string   `json:"minOccurs,omitempty"`
	MaxOccurs         string   `json:"maxOccurs,omitempty"`
	Description       string   `json:"description,omitempty"`
	Values            []string `json:"values,omitempty"`
	MinInclusive      string   `json:"minInclusive,omitempty"`
	MaxInclusive      string   `json:"maxInclusive,omitempty"`
}

type tf struct {
	t            string // type
	f            string // format
	minInclusive string
	maxInclusive string
}

var xsdTransMap map[string]tf

// Messages returns messages and message items (protobuf style)
// Could also align to json schema some time
// fmtStd is the format standard can be "protobuf" or "json"
func (xsd *XSD) Messages(fmtStd string) (messages []*Message, err error) {
	buildXsdTransMap(fmtStd)
	messageMap := make(map[string]*Message)
	fDisplay := func(xe XsdElement, h interface{}) (interface{}, error) {
		var currentMsg *Message
		if h != nil {
			currentMsg = h.(*Message)
		}
		// In general,we're only interested in complex types and elements
		switch t := xe.(type) {
		case *ComplexType:
			msg := &Message{sequence: len(messageMap)}
			if t.Name == "" { // If Complex Type doesn't have a name then element Name is the type
				if currentMsg == nil || len(currentMsg.MessageItems) == 0 {
					return currentMsg, nil
				}
				msg.Name = currentMsg.MessageItems[len(currentMsg.MessageItems)-1].Name
			} else {
				msg.Name = t.Name
				msg.IsRootMessage = currentMsg == nil
			}
			if _, inMap := messageMap[msg.Name]; !inMap {
				messageMap[msg.Name] = msg
				currentMsg = msg // Now the elements we come across belong to this message
			}
		case *SimpleType:
			if t.Name == "" {
				return currentMsg, nil
			}
			msg := &Message{sequence: len(messageMap)}
			if t.Name == "" { // If Simple Type doesn't have a name then element Name is the type
				if currentMsg == nil || len(currentMsg.MessageItems) == 0 {
					return currentMsg, fmt.Errorf("simpleType but no current message")
				}
				if len(currentMsg.MessageItems) == 0 {
					return currentMsg, fmt.Errorf("simpleType but no current message item")
				}
				msg.Name = currentMsg.MessageItems[len(currentMsg.MessageItems)-1].Name
			} else {
				msg.Name = t.Name
				msg.IsRootMessage = currentMsg == nil
			}
			if _, inMap := messageMap[msg.Name]; !inMap {
				messageMap[msg.Name] = msg
				currentMsg = msg // Now the elements we come across belong to this message
			}
		case *Attribute: // Attributes are added to a message
			if currentMsg == nil {
				return currentMsg, fmt.Errorf("attribute but no current message")
			}
			mi := &MessageItem{Name: t.Name, Repeated: false, MandatoryOptional: t.IsMandatoryOptional()}
			if t.Ref > "" {
				mi.Name, mi.Type = t.Ref, t.Ref
			} else {
				mi.setTypeOrMessage(mi.Type, messageMap)
			}
			currentMsg.MessageItems = append(currentMsg.MessageItems, mi)

		case *Extension: // Extension is extending an existing ComplexType, base is the baseline for the extension
			if currentMsg == nil {
				return currentMsg, fmt.Errorf("extension but no current message")
			}
			mi := &MessageItem{Name: currentMsg.Name, Repeated: false}
			mi.setTypeOrMessage(t.Base, messageMap)
			currentMsg.MessageItems = append(currentMsg.MessageItems, mi)

		case *Restriction: // Restriction provides more information about the current message item
			if currentMsg == nil {
				return currentMsg, fmt.Errorf("restriction but no current message")
			}
			// Create a message item if we haven't already
			if len(currentMsg.MessageItems) == 0 {
				mi := &MessageItem{Name: currentMsg.Name}
				mi.setTypeOrMessage(t.Base, messageMap)
				currentMsg.MessageItems = append(currentMsg.MessageItems, mi)
			} else {
				currentMsg.MessageItems[len(currentMsg.MessageItems)-1].setTypeOrMessage(t.Base, messageMap)
			}

		case *Element:
			if currentMsg == nil {
				currentMsg = &Message{sequence: len(messageMap), Name: t.Name, IsRootMessage: t.Name > ""}
				if _, inMap := messageMap[currentMsg.Name]; !inMap {
					messageMap[currentMsg.Name] = currentMsg
				}
				if t.Type == "" {
					return currentMsg, nil // Let something else create the message item
				}
			}
			mi := &MessageItem{
				Name:              t.Name,
				Repeated:          t.IsRepeated(),
				MandatoryOptional: t.IsMandatoryOptional(),
				MinOccurs:         t.MinOccurs,
				MaxOccurs:         t.MaxOccurs,
			}
			if t.Ref > "" {
				mi.Name, mi.Type = t.Ref, t.Ref
			} else {
				mi.setTypeOrMessage(t.Type, messageMap)
			}
			currentMsg.MessageItems = append(currentMsg.MessageItems, mi)

		case *Annotation:
			if currentMsg == nil {
				return currentMsg, fmt.Errorf("annotation but no current message")
			}
			if len(currentMsg.MessageItems) == 0 {
				currentMsg.Description = t.Documentation
			} else {
				currentMsg.MessageItems[len(currentMsg.MessageItems)-1].Description = t.Documentation
			}

		case *Pattern:
			if currentMsg == nil {
				return currentMsg, fmt.Errorf("pattern but no current message")
			}
			if len(currentMsg.MessageItems) == 0 {
				return currentMsg, fmt.Errorf("pattern but no current message item")
			} else {
				currentMsg.MessageItems[len(currentMsg.MessageItems)-1].Format = t.Value
			}

		case *MinInclusive:
			if currentMsg == nil {
				return currentMsg, fmt.Errorf("minInclusive but no current message")
			}
			if len(currentMsg.MessageItems) == 0 {
				return currentMsg, fmt.Errorf("minInclusive but no current message item")
			} else {
				currentMsg.MessageItems[len(currentMsg.MessageItems)-1].MinInclusive = t.Value
			}

		case *MaxInclusive:
			if currentMsg == nil {
				return currentMsg, fmt.Errorf("maxInclusive but no current message")
			}
			if len(currentMsg.MessageItems) == 0 {
				return currentMsg, fmt.Errorf("maxInclusive but no current message item")
			} else {
				currentMsg.MessageItems[len(currentMsg.MessageItems)-1].MaxInclusive = t.Value
			}

		case *Enumeration:
			if currentMsg == nil {
				return currentMsg, fmt.Errorf("enumeration but no current message")
			}
			if i := len(currentMsg.MessageItems) - 1; i >= 0 {
				currentMsg.MessageItems[i].Values = append(currentMsg.MessageItems[i].Values, t.Value)
			}
		}
		return currentMsg, nil
	}
	_, _ = xsd.ApplyFunctionP(fDisplay)
	for _, m := range messageMap {
		messages = append(messages, m)
	}
	sort.Slice(messages, func(i, j int) bool {
		if messages[i].Package != messages[j].Package {
			return messages[i].Package < messages[j].Package
		}
		if messages[i].IsRootMessage != messages[j].IsRootMessage {
			return messages[i].IsRootMessage
		}
		return messages[i].sequence < messages[j].sequence
	})
	return messages, nil
}

// setTypeOrMessages tries to convert internal xsd Types to other message types
// This relies upon xsdTransMap being set
func (mi *MessageItem) setTypeOrMessage(t string, messageMap map[string]*Message) {
	if t == "" {
		mi.Type = mi.Name
	} else {
		var inMap bool
		var typeFmt tf
		if typeFmt, inMap = xsdTransMap[t]; inMap {
			if typeFmt.t > "" {
				mi.Type = typeFmt.t
			}
			if typeFmt.f > "" {
				mi.Format = typeFmt.f
			}
			if typeFmt.minInclusive > "" {
				mi.MinInclusive = typeFmt.minInclusive
			}
			if typeFmt.maxInclusive > "" {
				mi.MaxInclusive = typeFmt.maxInclusive
			}
		} else {
			// Is this a Package:Message field
			if pkgMsg := strings.Split(t, ":"); len(pkgMsg) > 1 {
				mi.Type = strings.ReplaceAll(t, ":", ".")
				if len(pkgMsg) == 2 {
					if _, inMap := messageMap[mi.Type]; !inMap {
						messageMap[mi.Type] = &Message{
							Package:       pkgMsg[0],
							Name:          pkgMsg[1],
							MessageItems:  nil,
							Description:   "",
							IsRootMessage: false,
							sequence:      len(messageMap) + 1,
						}
					}
				}
			} else {
				mi.Type = t
			}
		}
	}
}

func buildXsdTransMap(formatStandard string) {
	switch formatStandard {
	case "protobuf":
		xsdTransMap = map[string]tf{
			"xs:string":           {t: "string"},
			"xs:normalizedString": {t: "string"},
			"xs:token":            {t: "string"},
			"xs:long":             {t: "int64"},
			"xs:int":              {t: "int64"},
			"xs:integer":          {t: "int64"},
			"xs:positiveInteger":  {t: "int64", minInclusive: "0"},
			"xs:float":            {t: "float"},
			"xs:decimal":          {t: "float"},
			"xs:double":           {t: "double"},
			"xs:boolean":          {t: "bool"},
			"xs:date":             {t: "google.protobuf.Timestamp"},
			"xs:datetime":         {t: "google.protobuf.Timestamp"},
			"xs:time":             {t: "google.protobuf.Timestamp"},
			"duration":            {t: "google.protobuf.Duration"},
		}
	case "json":
		xsdTransMap = map[string]tf{
			"xs:string":           {t: "string"},
			"xs:normalizedString": {t: "string"},
			"xs:token":            {t: "string"},
			"xs:long":             {t: "integer"},
			"xs:int":              {t: "integer"},
			"xs:integer":          {t: "integer"},
			"xs:positiveInteger":  {t: "integer", minInclusive: "0"},
			"xs:float":            {t: "number"},
			"xs:decimal":          {t: "number"},
			"xs:double":           {t: "number"},
			"xs:boolean":          {t: "boolean"},
			"xs:date":             {t: "string", f: "2018-11-13"},                //New in draft 7 Date
			"xs:datetime":         {t: "string", f: "2018-11-13T20:20:39+00:00"}, // Date and time together
			"xs:time":             {t: "string", f: "20:20:39+00:00"},            // New in draft 7 Time0
			"duration":            {t: "string"},                                 //New in draft 2019-09 A duration as defined by the ISO 8601 ABNF for "duration". For example, P3D expresses a duration of 3 days
		}
	default:
		xsdTransMap = map[string]tf{} // No translation
	}

}
