package xsd

import (
	"fmt"
	"sort"
	"strings"
)

type Message struct {
	Package      string         `json:"package,omitempty"`
	Name         string         `json:"name,omitempty"`
	MessageItems []*MessageItem `json:"messageItems,omitempty"`
	Description  string         `json:"description,omitempty"`
	IsNamed      bool           `json:"isNamed,omitempty"`
	sequence     int
}

type MessageItem struct {
	Name              string   `json:"name,omitempty"`
	Type              string   `json:"type,omitempty"`
	Repeated          bool     `json:"repeated,omitempty"`
	MandatoryOptional string   `json:"mandatoryOptional,omitempty"` // M=Mandatory,O=Optional, blank=don't know
	MinOccurs         string   `json:"minOccurs,omitempty"`
	MaxOccurs         string   `json:"maxOccurs,omitempty"`
	Description       string   `json:"description,omitempty"`
	Values            []string `json:"values,omitempty"`
	Pattern           string   `json:"pattern,omitempty"`
	MinInclusive      string   `json:"minInclusive,omitempty"`
	MaxInclusive      string   `json:"maxInclusive,omitempty"`
}

// Messages returns messages and message items (protobuf style)
// Could also align to json schema some time
func (xsd *XSD) Messages() (messages []*Message, err error) {
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
				msg.IsNamed = true
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
				msg.IsNamed = true
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
				mi.setTypeOrMessage(t.Type, messageMap)
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
				currentMsg = &Message{sequence: len(messageMap), Name: t.Name, IsNamed: t.Name > ""}
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
				currentMsg.MessageItems[len(currentMsg.MessageItems)-1].Pattern = t.Value
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
		if messages[i].IsNamed != messages[j].IsNamed {
			return messages[i].IsNamed
		}
		return messages[i].sequence < messages[j].sequence
	})
	return messages, nil
}

func (mi *MessageItem) setTypeOrMessage(t string, messageMap map[string]*Message) {
	switch t {
	case "xs:string", "xs:normalizedString", "xs:token":
		mi.Type = "string"
	case "xs:long":
		mi.Type = "int64"
	case "xs:int", "xs:integer":
		mi.Type = "int64"
	case "xs:float":
		mi.Type = "float"
	case "xs:double":
		mi.Type = "double"
	case "xs:boolean":
		mi.Type = "bool"
	case "xs:date", "xs:datetime", "xs:time":
		mi.Type = "google.protobuf.Timestamp"
	case "duration":
		mi.Type = "google.protobuf.Duration"
	default: // This must be a message
		if t > "" {
			// Is this a Package:Message field
			if pkgMsg := strings.Split(t, ":"); len(pkgMsg) > 1 {
				mi.Type = strings.ReplaceAll(t, ":", ".")
				if len(pkgMsg) == 2 {
					if _, inMap := messageMap[mi.Type]; !inMap {
						messageMap[mi.Type] = &Message{
							Package:      pkgMsg[0],
							Name:         pkgMsg[1],
							MessageItems: nil,
							Description:  "",
							IsNamed:      false,
							sequence:     len(messageMap) + 1,
						}
					}
				}
			} else {
				mi.Type = t
			}
		} else {
			mi.Type = mi.Name
		}

	}

}

// This maps to protobuf types
// 			switch t.Type {
//			case "xs:string", "xs:normalizedString", "xs:token":
//				mi.Type = "string"
//			case "xs:long":
//				mi.Type = "sint64"
//			case "xs:int", "xs:integer":
//				mi.Type = "int64"
//			case "xs:float":
//				mi.Type = "float"
//			case "xs:double":
//				mi.Type = "double"
//			case "xs:boolean":
//				mi.Type = "bool"
//			case "xs:date", "xs:datetime", "xs:time":
//				mi.Type = "google.protobuf.Timestamp"
//			case "duration":
//				mi.Type = "google.protobuf.Duration"
//			default: // This must be a message
//				if t.Type > "" {
//					// Is this a Package:Message field
//					if strings.Contains(t.Type, ":") {
//						mi.Type = strings.ReplaceAll(t.Type, ":", ".")
//						pkgMsg := strings.Split(t.Type, ":")
//						if len(pkgMsg) == 2 {
//							if _, inMap := messageMap[mi.Type]; !inMap {
//								messageMap[mi.Type] = &Message{
//									Package:      pkgMsg[0],
//									Name:         pkgMsg[1],
//									MessageItems: nil,
//									Description:      "",
//									IsNamed:      false,
//									sequence:     len(messageMap) + 1,
//								}
//							}
//						}
//					} else {
//						mi.Type = t.Type
//					}
//				} else {
//					mi.Type = t.Name
//				}
//
//			}
