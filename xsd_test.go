package xsd_test

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"os"
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
		if item.IsDir() {
			continue
		}
		// Only test with the test XSD
		//if !strings.HasPrefix(item.Name(), "test.xml") {
		//	continue
		//}
		t.Logf("Found file %s", item.Name())
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
		s := newXSD.ItemsString()
		t.Logf("%s", strings.Join(s, "\n"))
		messages := origXSD.Protobuf()
		pkg := ""
		for _, msg := range messages {
			if pkg != msg.Package {
				t.Logf("Package: %s", msg.Package)
				pkg = msg.Package
			}
			t.Logf("message %s {", msg.Name)
			for i, mi := range msg.MessageItems {
				repeated := ""
				if mi.Repeated {
					repeated = " repeated"
				}
				comment := mi.Comment
				if comment > "" {
					comment = " // " + comment
				}
				t.Logf("    %s %s %s = %d%s", repeated, mi.TypeOrMessage, mi.Name, i+1, comment)
			}
			t.Logf("} // End of %s", msg.Name)
		}
		//origXSD.Display()
	}
}

//		rules := []Rule{
//			{Id: "R_CD_CAM_BPD2021_DR_ZpjiKw", Type: "business", Code: "", Name: "Service Identifier", Description: "It is expected that the MSISDN field currently used in SNOW will accept any service identifier which can be an MSISDN or PSTN for Mobile or fixed service as defined in ITU E.164 or a service identifier for none-voice based services such as leased lines, FTTO, FTTH etc."},
//			{Id: "R_CD_CAM_BPD2021_DR_ZpkSFg", Type: "business", Code: "", Name: "Skip Customer Notification", Description: "If the case has the skip notification flag set then no notification to customer is required"},
//			{Id: "R_CD_CAM_BPD2021_DR_ZnZ8vg", Type: "business", Code: "", Name: "Skip approval", Description: "If the case has a Bulk Request ID or is pre-approved then no approval required"},
//			{Id: "R_CD_CAM_BPD2021_DR_ZpjiKw", Type: "business", Code: "", Name: "Service Identifier", Description: "It is expected that the MSISDN field currently used in SNOW will accept any service identifier which can be an MSISDN or PSTN for Mobile or fixed service as defined in ITU E.164 or a service identifier for none-voice based services such as leased lines, FTTO, FTTH etc."},
//			{Id: "R_CD_CAM_CAS_D3AP_ZlQztg", Type: "customerAttr", Code: "A-4-8", Name: "Billing :  Change Bill Bycle", Description: "To_BE confirmed the LOV of customerBillCycleId"},
//			{Id: "R_CD_CAM_BPD2021_DR_Zq-kPw", Type: "business", Code: "A-4-9", Name: "Subscriber : preferredLanguage", Description: "To_BE confirmed the LOV of preferredLanguage"},
//			{Id: "R_CD_NS_POS_MSR_AC_ZkCGeg", Type: "customerAttr", Code: "A-3-4", Name: "Other - Optional & Mandatory by Condition", Description: "Email, Contact Number, Location, etc"},
//			{Id: "R_CD_CAM_CAS_D3AP_ZlQzsA", Type: "customerAttr", Code: "A-4-2", Name: "Customer : Contact Info", Description: "Primary Contact ,Email Address ,Alternative Contact Number ,Office Number ,Home Number ,Preferred Contact Method ,Preferred Contact Time, Preferred Language"},
//			{Id: "R_CD_CAM_CAS_ZlctdA", Type: "customerAttr", Code: "A-4-5", Name: "Account : Billing Contact Info", Description: "Salution,Name,Address ,Email Address"},
//			{Id: "R_CD_CAM_CAS_D3AP_ZlQztQ", Type: "customerAttr", Code: "A-4-7", Name: "Billing :  Bill Medium : Paper,Email , e-Billing , CD,EMS_Bill, SMS", Description: "To_BE confirmed the LOV of Bill Media"},
//			{Id: "R_CD_CAM_CAS_D3AP_ZlQzsQ", Type: "customerAttr", Code: "A-4-3", Name: "Customer : Home Address", Description: "Address Type , Address"},
//			{Id: "R_CD_CAM_CAS_D3AP_ZlQztA", Type: "customerAttr", Code: "A-4-6", Name: " Account: Payment Mode", Description: "TBC the LOV of preferredPaymentMethod"},
//			{Id: "R_CD_CAM_CAS_D3AP_ZlQzrw", Type: "customerAttr", Code: "A-4-1", Name: "Customer : General Information", Description: "Customer Name,Nationality ,DOB, Race , IVR Language,CTC , Segment  ,OKUID and etc"},
//			{Id: "R_CD_CAM_CAS_D3AP_ZlQvTg", Type: "business", Code: "A-2-1", Name: "Customer Type Changes Allowed", Description: "Changes allowed configured by Customer Type & User Access Control"},
//			{Id: "R_CD_CAM_CAS_D3AP_ZlQvTw", Type: "business", Code: "A-3-2", Name: "Subscriber/Line Account Status", Description: "CRM Subscriber Status in Active"},
//		}
//
//		t.Logf("\nRules:")
//		assert.Len(t, origD.GetRules(), len(rules))
//		for i, r := range origD.GetRules() {
//			//t.Logf(`  {Id:"%s",Type:"%s",Code:"%s",Name:"%s",Description:"%s"},`, r.Id, r.Type, r.Code, r.Name, r.Description)
//			//_ = rules[i-i]
//			//t.Log(rules[i])
//			assert.Equal(t, r.Type, rules[i].Type)
//			assert.Equal(t, r.Code, rules[i].Code)
//			assert.Equal(t, r.Name, rules[i].Name)
//			assert.Equal(t, r.Description, rules[i].Description)
//		}
//		type testBE struct {
//			be     BaseElement
//			group  BaseElement
//			parent BaseElement
//		}
//		// Parent --> Child Hierarchy Rules
//		// Participant / Process are the same but Participany is the Base Element to use if found
//		// Collaboration is not part of the hierarchy
//		var baseElements = []testBE{
//			{be: &Task{Id: "Test", Name: "Can't find this", Documentation: "Not in BPMN"}, group: nil, parent: nil},
//			{be: &CategoryValue{Id: "CategoryValue_1jrdokb", Value: "Notification", Documentation: ""}, group: nil, parent: nil},
//			{
//				be: &EndEvent{Id: "Event_1nxge64", Name: "", Documentation: ""}, group: nil,
//				parent: &Lane{Id: "Lane_1w0jggb", Name: "Channel", Documentation: ""},
//			},
//			{be: &LaneSet{Id: "Id_f1213f36-2baf-44c0-a324-0e229d3a589e", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Id_11f0231c-aaca-4075-9816-701163f80c93", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be: &Lane{Id: "Lane_0tjvvj2", Name: "SNOW", Documentation: ""}, group: nil,
//				parent: &Lane{Id: "Lane_0jfx5j6", Name: "BSS", Documentation: ""},
//			},
//			{be: &SequenceFlow{Id: "Flow_0vkfvvb", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &Task{Id: "Id_10113598-8694-4f30-a238-1b3edbf567da", Name: "Review and Confirm Info", Documentation: "This will include any fraud, access restrictions"},
//				group:  &Group{Id: "Group_179pikd", Name: "Pre-validation Customer/Acct/Sub info", Documentation: ""},
//				parent: &Lane{Id: "Lane_1w0jggb", Name: "Channel", Documentation: ""},
//			},
//			{be: &SequenceFlow{Id: "Flow_1sf6eqh", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &CallActivity{Id: "Id_106ace90-cd83-4b10-87c3-e0d0c45a97aa", Name: "KYC Processing", Documentation: "Selection of Active Customer, Subscriber or Account made in sub-process"},
//				group:  &Group{Id: "Group_0xfiy9v", Name: "KYC Processing", Documentation: ""},
//				parent: &Lane{Id: "Lane_1w0jggb", Name: "Channel", Documentation: ""},
//			},
//			{be: &StartEvent{Id: "Event_0byk73z", Name: "", Documentation: ""}, group: nil,
//				parent: &Lane{Id: "Lane_1w0jggb", Name: "Channel", Documentation: ""},
//			},
//			{
//				be:     &Task{Id: "Activity_1qquwgi", Name: "Updated", Documentation: ""},
//				group:  &Group{Id: "Group_1o2px5j", Name: "Notification", Documentation: ""},
//				parent: &Lane{Id: "Lane_1w0jggb", Name: "Channel", Documentation: ""},
//			},
//			{be: &SequenceFlow{Id: "Flow_0c5mdcq", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be: &Lane{Id: "Lane_1w0jggb", Name: "Channel", Documentation: ""}, group: nil,
//				parent: &Lane{Id: "Lane_0jfx5j6", Name: "BSS", Documentation: ""},
//			},
//			{
//				be: &EndEvent{Id: "Event_1dutsc5", Name: "", Documentation: ""}, group: nil,
//				parent: &Participant{Id: "Id_6e3ea679-5a59-40f3-9376-452c13335907", Name: "Customer", Documentation: "<p>&nbsp;</p>"},
//			},
//			{be: &Collaboration{Id: "Id_5181efce-7e1e-4f72-bab6-c5f1b88bb26f", Name: "BP122-D3-CAS", Documentation: ""}, group: nil, parent: nil},
//			{be: &Group{Id: "Group_179pikd", Name: "Pre-validation Customer/Acct/Sub info", Documentation: ""}, group: nil, parent: nil},
//			{be: &CategoryValue{Id: "CategoryValue_ZlbiqQ", Value: "Customer/Acct/Sub Info Capturing & Processing", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &Group{Id: "Group_0xfiy9v", Name: "KYC Processing", Documentation: ""},
//				group:  nil,
//				parent: nil,
//			},
//			{
//				be:     &UserTask{Id: "Id_28feca90-5ee5-4f29-a1a3-19805274cc8c", Name: "Request change information", Documentation: ""},
//				group:  &Group{Id: "Group_0xfb3ag", Name: "Customer/Acct/Sub Info Capturing & Processing", Documentation: ""},
//				parent: &Participant{Id: "Id_6e3ea679-5a59-40f3-9376-452c13335907", Name: "Customer", Documentation: "<p>&nbsp;</p>"},
//			},
//			{be: &MessageFlow{Id: "Flow_1oywhls", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Flow_0xs4cbp", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Id_3e49d651-be54-48a1-a673-d0b958f3d370", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &ChildLaneSet{Id: "LaneSet_164o496", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Flow_1geyr5h", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &Group{Id: "Group_0xfb3ag", Name: "Customer/Acct/Sub Info Capturing & Processing", Documentation: ""}, group: nil, parent: nil},
//			{be: &CategoryValue{Id: "CategoryValue_0ij1yqt", Value: "Pre-validation Customer/Acct/Sub info", Documentation: ""}, group: nil, parent: nil},
//			{be: &LaneSet{Id: "Id_6e36b251-2fff-49bc-9888-ae7b2898cfee", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &CallActivity{Id: "Activity_0skg3vn", Name: "KYC Processing", Documentation: ""},
//				group:  &Group{Id: "Group_0xfiy9v", Name: "KYC Processing", Documentation: ""},
//				parent: &Lane{Id: "Lane_0fxgohc", Name: "DTE", Documentation: ""},
//			},
//			{
//				be: &Lane{Id: "Lane_0fxgohc", Name: "DTE", Documentation: ""}, group: nil,
//				parent: &Lane{Id: "Lane_0jfx5j6", Name: "BSS", Documentation: ""},
//			},
//			{be: &SequenceFlow{Id: "Flow_0hlspwl", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &Task{Id: "Activity_0ef3hwm", Name: "Submit Update Details", Documentation: "Submitting the update should return an ID for tracking the update"},
//				group:  &Group{Id: "Group_0xfb3ag", Name: "Customer/Acct/Sub Info Capturing & Processing", Documentation: ""},
//				parent: &Lane{Id: "Lane_1w0jggb", Name: "Channel", Documentation: ""},
//			},
//			{be: &MessageFlow{Id: "Flow_1ftan93", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &CategoryValue{Id: "CategoryValue_1nwb2n0", Value: "KYC Processing", Documentation: ""},
//				group:  nil,
//				parent: nil,
//			},
//			{
//				be:     &Task{Id: "Activity_1usn017", Name: "Update Details Request", Documentation: ""},
//				group:  &Group{Id: "Group_0xfb3ag", Name: "Customer/Acct/Sub Info Capturing & Processing", Documentation: ""},
//				parent: &Lane{Id: "Lane_0fxgohc", Name: "DTE", Documentation: ""},
//			},
//			{be: &Process{Id: "Id_4205379a-6bc7-4b9f-841a-5b8d73645f18", Name: "Customer Experience", Documentation: "<p>&nbsp;</p>"}, group: nil, parent: nil},
//			{be: &Process{Id: "Id_3759fac8-ca5f-427f-918d-43a5e77d2500", Name: "CelcomDigi", Documentation: ""}, group: nil, parent: nil},
//			{be: &ChildLaneSet{Id: "LaneSet_1x7wl8o", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Flow_0n6twdq", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Flow_050ip18", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &CallActivity{Id: "Id_ad9a2c97-3494-42a1-930b-1fec62c44fd2", Name: "KYC Processing", Documentation: "Selection of Active Customer, Subscriber or Account made in sub-process"},
//				group:  &Group{Id: "Group_0xfiy9v", Name: "KYC Processing", Documentation: ""},
//				parent: &Participant{Id: "Id_6e3ea679-5a59-40f3-9376-452c13335907", Name: "Customer", Documentation: "<p>&nbsp;</p>"},
//			},
//			{be: &Group{Id: "Group_1o2px5j", Name: "Notification", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Flow_0yavxvo", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Id_a0551418-91f7-4a5e-b81d-28b87b2d991e", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be: &Lane{Id: "Lane_0jfx5j6", Name: "BSS", Documentation: ""}, group: nil,
//				parent: &Participant{Id: "Id_e09dbcb2-02f4-46a3-8340-8ce672e9a13f", Name: "CelcomDigi", Documentation: ""},
//			},
//			{
//				be:     &CallActivity{Id: "Activity_0ltxm7x", Name: "Request Case Process", Documentation: ""},
//				group:  &Group{Id: "Group_0xfb3ag", Name: "Customer/Acct/Sub Info Capturing & Processing", Documentation: ""},
//				parent: &Lane{Id: "Lane_0tjvvj2", Name: "SNOW", Documentation: ""},
//			},
//			{be: &SequenceFlow{Id: "Flow_05vbvss", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &UserTask{Id: "Id_4755704f-fd95-42f9-acbd-f8c2fd3a7e09", Name: "Review and Confirm Info", Documentation: ""},
//				group:  &Group{Id: "Group_179pikd", Name: "Pre-validation Customer/Acct/Sub info", Documentation: ""},
//				parent: &Participant{Id: "Id_6e3ea679-5a59-40f3-9376-452c13335907", Name: "Customer", Documentation: "<p>&nbsp;</p>"},
//			},
//			{be: &Participant{Id: "Id_e09dbcb2-02f4-46a3-8340-8ce672e9a13f", Name: "CelcomDigi", Documentation: ""}, group: nil, parent: nil},
//			{be: &Participant{Id: "Id_6e3ea679-5a59-40f3-9376-452c13335907", Name: "Customer", Documentation: "<p>&nbsp;</p>"}, group: nil, parent: nil},
//			{be: &MessageFlow{Id: "Flow_110iwqb", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Flow_08bo7o1", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &CallActivity{Id: "Id_455ba7d3-8062-4e4c-a575-dd4085a20e11", Name: "Display 360 Info", Documentation: "User access may restrict the data returned to the channel"},
//				group:  &Group{Id: "Group_179pikd", Name: "Pre-validation Customer/Acct/Sub info", Documentation: ""},
//				parent: &Lane{Id: "Lane_0fxgohc", Name: "DTE", Documentation: ""},
//			},
//			{
//				be: &StartEvent{Id: "Id_4753722c-d8ae-43f0-8dd9-3d877d904a37", Name: "", Documentation: ""}, group: nil,
//				parent: &Participant{Id: "Id_6e3ea679-5a59-40f3-9376-452c13335907", Name: "Customer", Documentation: "<p>&nbsp;</p>"},
//			},
//			{be: &Category{Id: "Category_0mpk3vc", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{be: &SequenceFlow{Id: "Flow_14roaqe", Name: "", Documentation: ""}, group: nil, parent: nil},
//			{
//				be:     &Task{Id: "Id_67b8f4b1-f016-4e3d-a1b9-d273dcd7fd00", Name: "Capture change requirement", Documentation: ""},
//				group:  &Group{Id: "Group_0xfb3ag", Name: "Customer/Acct/Sub Info Capturing & Processing", Documentation: ""},
//				parent: &Lane{Id: "Lane_1w0jggb", Name: "Channel", Documentation: ""},
//			},
//		}
//
//		origD.BpmnIdBaseElementMap()
//		// BaseElementsMap
//		t.Logf("\nBaseElementsMap:%d", len(origD._BaseElementMap))
//		assert.Len(t, baseElements, len(origD._BaseElementMap)+1) // one -ve test element
//		for _, testBaseElement := range baseElements {
//			var be BaseElement
//			be = origD.FindBaseElementById(testBaseElement.be.GetId())
//			if testBaseElement.be.GetId() == "Test" {
//				assert.Nil(t, be, "found test element!")
//				continue
//			}
//			assert.NotNil(t, be, "could not find element!")
//			if be == nil {
//				continue
//			}
//			assert.Equal(t, testBaseElement.be.GetId(), be.GetId())
//			assert.Equal(t, testBaseElement.be.GetName(), be.GetName())
//			assert.Equal(t, testBaseElement.be.GetDocumentation(), be.GetDocumentation())
//
//			beGroup := origD.Group(be)
//			if testBaseElement.group == nil {
//				assert.Nil(t, beGroup, testBaseElement.be.ToString())
//			} else {
//				assert.NotNil(t, beGroup, testBaseElement.be.ToString())
//				if beGroup != nil {
//					assert.Equal(t, testBaseElement.group.GetName(), beGroup.GetName(), testBaseElement.be.ToString())
//					if testBaseElement.group.GetName() == beGroup.GetName() {
//						assert.Equal(t, testBaseElement.group.GetId(), beGroup.GetId(), testBaseElement.be.ToString())
//						assert.Equal(t, testBaseElement.group.GetDocumentation(), beGroup.GetDocumentation(), testBaseElement.be.ToString())
//					}
//				}
//			}
//
//			beParent := origD.Parent(be)
//			if testBaseElement.parent == nil {
//				assert.Nil(t, beParent, testBaseElement.be.ToString())
//			} else {
//				assert.NotNil(t, beParent, testBaseElement.be.ToString())
//				if beParent != nil {
//					assert.Equal(t, testBaseElement.parent.GetName(), beParent.GetName(), testBaseElement.be.ToString())
//					if testBaseElement.parent.GetName() == beParent.GetName() {
//						assert.Equal(t, testBaseElement.parent.GetId(), beParent.GetId(), testBaseElement.be.ToString())
//						assert.Equal(t, testBaseElement.parent.GetDocumentation(), beParent.GetDocumentation(), testBaseElement.be.ToString())
//					}
//				}
//			}
//		}
//
//		// Builds the test data
//		//for _, be := range origD._BaseElementMap {
//		//	t.Logf(`{be:&%s{Id:"%s",Name:"%s",Documentation:"%s"},group: nil,parent: nil},`,
//		//		strings.Title(be.GetXMLName().Local), be.GetId(), be.GetName(), be.GetDocumentation())
//		//	//t.Logf(".. %s:%v", bpmnID, be)
//		//}
//
//		origD.BpmnIdGroupMap()
//		t.Logf("\nGroupMap:%d", len(origD._BaseElementGroup))
//		for bpmnID, group := range origD._BaseElementGroup {
//			t.Logf(".. %s belongs to group %s\n", origD._BaseElementMap[bpmnID].ToString(), group.ToString())
//		}
//
//		origD.BpmnIdParentMap()
//		t.Logf("\nParentMap:%d", len(origD._BaseElementParent))
//		for bpmnID, parent := range origD._BaseElementParent {
//			t.Logf(".. '%s', %s has parent %s\n", bpmnID, origD._BaseElementMap[bpmnID].ToString(), parent.ToString())
//		}
//		//t.Log(origD._BaseElementParent["_6-61"].ToString())
//
//		//testStringGraph(t)
//
//		// Process Flow Hierarchy
//		//for _, p := range origD.Processes {
//		//	t.Logf("\nProcess FLow:%s", p.ToString())
//		//	if p.Name == "CelcomDigi" {
//		//		for _, id := range []string{"Activity_1qvc83h", "Event_1fdhkgw", "Lane_18l2dnr",
//		//			"Id_8235bf52-0e42-4df2-903f-6a34f380d702"} {
//		//			assert.Equal(t, id, p.FindBaseElementById(id).GetId(), "FindBaseElementById")
//		//		}
//		//	}
//		//	//t.Logf("Sequence Flows:%v", p.SequenceFlows)
//		//	for _, tbe := range p.TopologicalSort() {
//		//		var rules []string
//		//		for _, r := range tbe.BaseElement.GetExtensionElement().GetRules() {
//		//			rules = append(rules, fmt.Sprintf("(%s) %s %s", r.Type, r.Name, r.Description))
//		//		}
//		//		if tbe.BaseElement.GetType() == B2SequenceFlow {
//		//			t.Logf("  %s: %s: %s (%s) : %s : %s",
//		//				origD.GroupName(tbe.BaseElement),
//		//				tbe.BaseElement.GetType(),
//		//				tbe.BaseElement.GetName(),
//		//				tbe.BaseElement.GetId(),
//		//				strings.Join(rules, ", "),
//		//				origD.ParentName(tbe.BaseElement))
//		//		} else {
//		//			t.Logf("  %s: %s: %s (%s)  : %s : %s: %v",
//		//				origD.GroupName(tbe.BaseElement),
//		//				tbe.BaseElement.GetType(),
//		//				tbe.BaseElement.GetName(),
//		//				tbe.BaseElement.GetId(),
//		//				strings.Join(rules, ", "),
//		//				origD.Parent(tbe.BaseElement).GetId(),
//		//				origD.ShapeOfBaseElement(origD.Parent(tbe.BaseElement)).Id)
//		//		}
//		//	}
//		//}
//
//		// Check we can't find a baseElement
//		be := origD.FindBaseElementById("tim")
//		assert.Nil(t, be, "could not find base element should be NIL")
//
//		break
//	}
//}

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
