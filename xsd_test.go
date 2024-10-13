package xsd_test

import (
	"encoding/json"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"strings"
	"testing"
	"xsd"
)

// TestXML reads the XSD files and unmarshalls them, then marshalls them and finally checks some stuff
func TestXML(t *testing.T) {
	t.Log("TestUnmarshalXML")
	dir := "./xsd"
	items, _ := os.ReadDir(dir)
	for _, item := range items {
		ext := path.Ext(item.Name())
		if item.IsDir() || ext == ".json" {
			continue
		}
		// Uncomment to limit what gets tested
		//if !strings.HasPrefix(item.Name(), "test") {
		//	continue
		//}
		jsonFile := path.Join(dir, strings.ReplaceAll(item.Name(), ext, ".json"))
		t.Logf("XSD file %s, json file:%s", item.Name(), jsonFile)
		// Read the file
		var xsdXML []byte
		var err error
		if xsdXML, err = os.ReadFile(dir + "/" + item.Name()); err != nil {
			t.Fatalf("could not read the XML file, got %v", err)
		}

		// Now unmarshall into structs
		var origXSD *xsd.XSD
		if origXSD, err = xsd.NewXSD(xsdXML); err != nil {
			t.Fatalf("could not unmarshal XML into XSD, got %v", err)
		}

		//t.Logf("%v", origD)
		// Now marshall back again into XML
		if xsdXML, err = xml.Marshal(origXSD); err != nil {
			t.Fatalf("could not marshal XML from standard map, got %v", err)
		}

		// Now unmarshall the newly marshalled XML
		var newXSD *xsd.XSD
		if newXSD, err = xsd.NewXSD(xsdXML); err != nil {
			t.Fatalf("could not unmarshal XML into new definition, got %v", err)
		}

		// Now see if they are the same
		assert.True(t, compareDefinitions(t, origXSD, newXSD), "Something changed to marshal / unmarshal")
		//t.Logf("%s", origXSD.ToString())

		// Load the expected result
		var expect []*xsd.Message
		var expectByte []byte
		if expectByte, err = os.ReadFile(jsonFile); err != nil {
			if !os.IsNotExist(err) {
				t.Fatalf("could not read the XML file, got %v", err)
			}
		} else {
			if err = json.Unmarshal(expectByte, &expect); err != nil {
				t.Fatalf("could not unmarshal json into message map, got %v", err)
			}
		}
		noJson := len(expect) == 0
		//noJson = true // Description this out if not checking manually
		// Uncomment to see each element being visited and the order of the visit
		//s := newXSD.ItemsString()
		//t.Logf("%s", strings.Join(s, "\n"))
		//}

		var messages []*xsd.Message
		messages, err = origXSD.Messages("protobuf")
		if !assert.NoError(t, err, "converting xsd to messages") {
			continue
		}
		pkg := ""
		if noJson {
			t.Logf("%d messages", len(messages))
		} else {
			if !assert.Len(t, messages, len(expect)) {
				continue
			}
		}

		for m, msg := range messages {
			if pkg != msg.Package {
				if noJson {
					t.Logf("Package: %s", msg.Package)
				}
				pkg = msg.Package
			}
			var expectMsg *xsd.Message
			if noJson {
				t.Logf("message %s {", msg.Name)
			} else {
				expectMsg = expect[m]
				if !assert.Equal(t, expectMsg.Name, msg.Name) {
					continue
				}
			}
			for i, mi := range msg.MessageItems {
				if noJson {
					var attributes []string
					if mi.Repeated {
						attributes = append(attributes, "repeated")
					}
					switch mi.MandatoryOptional {
					case "M":
						attributes = append(attributes, "mandatory")
					case "O":
						attributes = append(attributes, "optional")
					}
					attributes = append(attributes, mi.Type)
					var comments []string
					if mi.Description != "" {
						comments = append(comments, mi.Description)
					}
					if len(mi.Values) > 0 {
						comments = append(comments, "Values:"+strings.Join(mi.Values, ", "))
					}
					if mi.Format > "" {
						comments = append(comments, "Format:"+mi.Format)
					}
					if mi.MinInclusive > "" {
						comments = append(comments, "MinInclusive:"+mi.MinInclusive)
					}
					if mi.MaxInclusive > "" {
						comments = append(comments, "MaxInclusive:"+mi.MaxInclusive)
					}
					var comment string
					if len(comments) > 0 {
						comment = " // " + strings.Join(comments, ",")
					}

					t.Logf("    %s %s = %d%s", strings.Join(attributes, " "), mi.Name, i+1, comment)
				} else {
					if !assert.Equal(t, mi, expectMsg.MessageItems[i]) {
						continue
					}
				}
			}
			if noJson {
				t.Logf("} // End of %s", msg.Name)
			}
		}
		// If there is noJson then write a file
		if noJson {
			if expectByte, err = json.MarshalIndent(messages, "", "  "); err != nil {
				t.Logf("unable to marshall result, %v", err)
			} else {
				if err = os.WriteFile(jsonFile, expectByte, 0644); err != nil {
					t.Logf("unable to save result, %v", err)
				}
			}
		}

		//origXSD.Display()
	}
}

func compareDefinitions(t *testing.T, xsd1 *xsd.XSD, xsd2 *xsd.XSD) bool {
	if assert.Equal(t, xsd1.Import, xsd2.Import) {
		if assert.Equal(t, xsd1.ComplexTypes, xsd2.ComplexTypes) {
			return true
		}
	}
	return false
}

//func compareBaseElements(t *testing.T, b1 BaseElement, b2 BaseElement) {
//	assert.Equal(t, b1.GetId(), b2.GetId())
//	assert.Equal(t, b1.GetName(), b2.GetName())
//	assert.Equal(t, b1.GetOutgoingAssociations(), b2.GetOutgoingAssociations())
//	assert.Equal(t, b1.GetIncomingAssociations(), b2.GetIncomingAssociations())
//	//var ok = "OK"
//	//if b1.GetId() != b2.GetId() || b1.GetName() != b2.GetName() {
//	//	ok = "NOT OK"
//	//}
//	//for i, oa := range b1.GetOutgoingAssociations() {
//	//	if oa != b2.GetOutgoingAssociations()[i] {
//	//		ok = "NOT OK"
//	//	}
//	//}
//	//for i, ia := range b1.GetIncomingAssociations() {
//	//	if ia != b2.GetIncomingAssociations()[i] {
//	//		ok = "NOT OK"
//	//	}
//	//}
//	//t.Logf("%s --> %s", b1.ToString(), ok)
//
//	for i, r1 := range b1.GetExtensionElement().GetRules() {
//		r2 := b2.GetExtensionElement().GetRules()[i]
//		assert.Equal(t, r1.Id, r2.Id)
//		assert.Equal(t, r1.Name, r2.Name)
//		assert.Equal(t, r1.Type, r2.Type)
//		assert.Equal(t, r1.Description, r2.Description)
//		//ok = "OK"
//		//if r1.Id != r2.Id || r1.Name != r2.Name || r1.Type != r2.Type || r1.Description != r2.Description {
//		//	ok = "NOT OK"
//		//}
//		//t.Logf(".. %s rule %s (%s) --> %s", r1.Type, r1.Id, r1.Name, ok)
//	}
//
//}
//
////
////func displayDefinitions(t *testing.T, d *Definition) {
////	for _, p := range d.Processes {
////		t.Logf("Process definition %s", p.Id)
////		for _, task := range p.Tasks {
////			t.Logf("..Task definition %s (%s)", task.Id, task.Name)
////		}
////		for _, sp := range p.SubProcesses {
////			t.Logf(".. Sub Process definition %s (%s)", sp.Id, sp.Name)
////			for _, task := range sp.Tasks {
////				t.Logf("   ..Task definition %s (%s)", task.Id, task.Name)
////			}
////		}
////	}
////}
