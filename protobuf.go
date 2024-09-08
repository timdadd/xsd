package xsd

import (
	"fmt"
	"sort"
	"strings"
)

type Message struct {
	Package      string
	Name         string
	MessageItems []*MessageItem
	Comment      string
	IsNamed      bool
	sequence     int
}

type MessageItem struct {
	TypeOrMessage string
	Name          string
	Repeated      bool
	Comment       string
}

// Protobuf returns protobuf type message definitions
func (xsd *XSD) Protobuf() (messages []*Message) {
	messageMap := make(map[string]*Message)
	fDisplay := func(xe XsdElement, h interface{}) (bool, interface{}) {
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
					panic("Where to get complex type name from?")
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
		case *Element:
			if currentMsg == nil {
				panic(fmt.Sprintf("What message to add element %s to?", t.ToString()))
			}
			mi := &MessageItem{Name: t.Name, Repeated: t.MaxOccurs == "unbounded"}
			switch t.Type {
			case "xs:string", "xs:normalizedString", "xs:token":
				mi.TypeOrMessage = "string"
			case "xs:long":
				mi.TypeOrMessage = "sint64"
			case "xs:int", "xs:integer":
				mi.TypeOrMessage = "int64"
			case "xs:float":
				mi.TypeOrMessage = "float"
			case "xs:double":
				mi.TypeOrMessage = "double"
			case "xs:boolean":
				mi.TypeOrMessage = "bool"
			case "xs:date", "xs:datetime", "xs:time":
				mi.TypeOrMessage = "google.protobuf.Timestamp"
			case "duration":
				mi.TypeOrMessage = "google.protobuf.Duration"
			default: // This must be a message
				if t.Type > "" {
					// Is this a Package:Message field
					if strings.Contains(t.Type, ":") {
						mi.TypeOrMessage = strings.ReplaceAll(t.Type, ":", ".")
						pkgMsg := strings.Split(t.Type, ":")
						if len(pkgMsg) == 2 {
							if _, inMap := messageMap[mi.TypeOrMessage]; !inMap {
								messageMap[mi.TypeOrMessage] = &Message{
									Package:      pkgMsg[0],
									Name:         pkgMsg[1],
									MessageItems: nil,
									Comment:      "",
									IsNamed:      false,
									sequence:     len(messageMap) + 1,
								}
							}
						}
					} else {
						mi.TypeOrMessage = t.Type
					}
				} else {
					mi.TypeOrMessage = t.Name
				}

			}
			currentMsg.MessageItems = append(currentMsg.MessageItems, mi)
		case *Annotation:
			if currentMsg == nil {
				panic(fmt.Sprintf("What message to add annotation %s to?", t.ToString()))
			}
			if len(currentMsg.MessageItems) == 0 {
				currentMsg.Comment = t.Documentation
			} else {
				currentMsg.MessageItems[len(currentMsg.MessageItems)-1].Comment = t.Documentation
			}

		}
		return true, currentMsg
	}
	xsd.ApplyFunctionP(fDisplay)
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
	return messages
}

// This maps to protobuf types
// 			switch t.Type {
//			case "xs:string", "xs:normalizedString", "xs:token":
//				mi.TypeOrMessage = "string"
//			case "xs:long":
//				mi.TypeOrMessage = "sint64"
//			case "xs:int", "xs:integer":
//				mi.TypeOrMessage = "int64"
//			case "xs:float":
//				mi.TypeOrMessage = "float"
//			case "xs:double":
//				mi.TypeOrMessage = "double"
//			case "xs:boolean":
//				mi.TypeOrMessage = "bool"
//			case "xs:date", "xs:datetime", "xs:time":
//				mi.TypeOrMessage = "google.protobuf.Timestamp"
//			case "duration":
//				mi.TypeOrMessage = "google.protobuf.Duration"
//			default: // This must be a message
//				if t.Type > "" {
//					// Is this a Package:Message field
//					if strings.Contains(t.Type, ":") {
//						mi.TypeOrMessage = strings.ReplaceAll(t.Type, ":", ".")
//						pkgMsg := strings.Split(t.Type, ":")
//						if len(pkgMsg) == 2 {
//							if _, inMap := messageMap[mi.TypeOrMessage]; !inMap {
//								messageMap[mi.TypeOrMessage] = &Message{
//									Package:      pkgMsg[0],
//									Name:         pkgMsg[1],
//									MessageItems: nil,
//									Comment:      "",
//									IsNamed:      false,
//									sequence:     len(messageMap) + 1,
//								}
//							}
//						}
//					} else {
//						mi.TypeOrMessage = t.Type
//					}
//				} else {
//					mi.TypeOrMessage = t.Name
//				}
//
//			}
