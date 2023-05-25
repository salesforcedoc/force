package command

import (
	"container/list"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ForceCLI/force/desktop"
	. "github.com/ForceCLI/force/error"
	. "github.com/ForceCLI/force/lib"
)

////////////////////////////////////////////////////////////////////////
// Parse the permissions for a given profile and return a Profile struct
////////////////////////////////////////////////////////////////////////

type XLS interface {
	addProperty(name string, value string)
	addToProfile(p Profile)
}

//--
type OLS struct {
	objectName       string
	allowCreate      string
	allowDelete      string
	allowEdit        string
	allowRead        string
	modifyAllRecords string
	viewAllRecords   string
}

func (o *OLS) addProperty(name string, value string) {
	switch name {
	case "object":
		o.objectName = value
	case "allowCreate":
		o.allowCreate = value
	case "allowDelete":
		o.allowDelete = value
	case "allowEdit":
		o.allowEdit = value
	case "allowRead":
		o.allowRead = value
	case "modifyAllRecords":
		o.modifyAllRecords = value
	case "viewAllRecords":
		o.viewAllRecords = value
	}

	//	fmt.Println("Object Property " + name + "=" + value)
}
func (o *OLS) getProperty(name string) string {
	switch name {
	case "Allow Create":
		return o.allowCreate
	case "Allow Delete":
		return o.allowDelete
	case "Allow Edit":
		return o.allowEdit
	case "Allow Read":
		return o.allowRead
	case "Modify All Records":
		return o.modifyAllRecords
	case "View All Records":
		return o.viewAllRecords
	}
	return ""
}
func (o *OLS) addToProfile(p Profile) {
	p.objectPermissions[o.objectName] = *o
}
//--
type FLS struct {
	field    string
	editable string
	readable string
}
func (f *FLS) addProperty(name string, value string) {
	switch name {
	case "field":
		f.field = value
	case "editable":
		f.editable = value
	case "readable":
		f.readable = value
	}
	//	fmt.Println("Field Property " + name + "=" + value)
}
func (f *FLS) addToProfile(p Profile) {
	p.fieldPermissions[f.field] = *f
}
//--
type UPERM struct {
	name    string
	enabled string
}
func (f *UPERM) addProperty(name string, value string) {
	switch name {
	case "name":
		f.name = value
	case "enabled":
		f.enabled = value
	}
}
func (f *UPERM) addToProfile(p Profile) {
	p.userPermissions[f.name] = *f
}
//-
type Profile struct {
	name              string
	fieldPermissions  map[string]FLS
	objectPermissions map[string]OLS
	userPermissions   map[string]UPERM
}

func parseProfileXML(profileName string, text string) Profile {
	p := new(Profile)
	p.name = profileName
	p.fieldPermissions = map[string]FLS{}
	p.objectPermissions = map[string]OLS{}
	p.userPermissions = map[string]UPERM{}
	var currentElement XLS

	r := strings.NewReader(text)
	parser := xml.NewDecoder(r)
	depth := 0

	eltType := ""
	propertyName := ""

	for {

		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			if depth == 1 {
				eltType = name
				if eltType == "objectPermissions" {
					currentElement = new(OLS)
				} else if eltType == "fieldPermissions" {
					currentElement = new(FLS)
				} else if eltType == "userPermissions" {
					currentElement = new(UPERM)
				} else {
					currentElement = nil
				}
			}
			if depth == 2 {
				propertyName = name
			}
			depth++
		case xml.EndElement:
			if depth == 2 && currentElement != nil {
				currentElement.addToProfile(*p)
			}
			depth--
		case xml.CharData:
			bytes := xml.CharData(t)
			if currentElement != nil && depth == 3 {
				currentElement.addProperty(propertyName, string(bytes))
			}
		default:
		}
	}

	//	fmt.Println(p)

	return *p
}

//////////////////////////////////////////////////////////////////////
// Read information about an SObject and returns a CustomObject struct
//////////////////////////////////////////////////////////////////////

type CustomObject struct {
	objectName	string
	fieldNames	[]string
	nbFields	int
}

func (co *CustomObject) addField(name string) {
	co.fieldNames[co.nbFields] = name
	co.nbFields++
}
func (co *CustomObject) getProfileFootprint(p Profile) string {
	key := "OLS:" + p.objectPermissions[co.objectName].allowCreate + "," +
		p.objectPermissions[co.objectName].allowDelete + "," +
		p.objectPermissions[co.objectName].allowEdit + "," +
		p.objectPermissions[co.objectName].allowRead + "," +
		p.objectPermissions[co.objectName].modifyAllRecords + "," +
		p.objectPermissions[co.objectName].viewAllRecords + ","

	for idx := 0; idx < co.nbFields; idx++ {
		f := co.fieldNames[idx]
		key += f + ":" + p.fieldPermissions[co.objectName+"."+f].editable + "," +
			p.fieldPermissions[co.objectName+"."+f].readable + ","
	}
	return key
}

func parseCustomObjectXML(objectName string, text string) CustomObject {
	obj := CustomObject{objectName: objectName, nbFields: 0, fieldNames: make([]string, 900, 900)}
	r := strings.NewReader(text)
	parser := xml.NewDecoder(r)
	depth := 0
	var firstLevel, secondLevel string

	for {

		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			if depth == 1 {
				firstLevel = name
			} else if depth == 2 {
				secondLevel = name
			}
			depth++
		case xml.EndElement:
			if depth == 3 {
				secondLevel = ""
			} else if depth == 2 {
				firstLevel = ""
			}
			depth--
		case xml.CharData:
			bytes := xml.CharData(t)
			if depth == 3 && firstLevel == "fields" && secondLevel == "fullName" {
				obj.addField(string(bytes))
			}
		default:
		}
	}

	//	fmt.Println(obj)
	return obj
}

//////////////////////////////////////////////////////////////////////
// Read information about an Profile and returns a ProfileObject struct
//////////////////////////////////////////////////////////////////////

type ProfileObject struct {
	profileName	string
	permNames	[]string
	nbPerms		int
}

func stringArrayContains(a []string, x string) bool {
	for _, n := range a {
			if x == n {
					return true
			}
	}
	return false
}
func (obj *ProfileObject) addPerm(name string) {
	if (stringArrayContains(obj.permNames, name)==false) {
		//fmt.Printf("addPerm: %d %s\n", obj.nbPerms, name)
		obj.permNames[obj.nbPerms] = name
		obj.nbPerms++
		permNames.Set(name);
		//sort.Strings(obj.permNames)
	}	
}
func (obj *ProfileObject) getFootprint(p Profile) string {
	key := "UPERM:"
	for idx := 0; idx < obj.nbPerms; idx++ {
		f := obj.permNames[idx]
		key += f + ":" + p.userPermissions[f].enabled + ","
	}
	return key
}
func parseProfileObjectXML(profileName string, text string, obj ProfileObject) ProfileObject {
	//obj := ProfileObject{profileName: profileName, nbPerms: 0, permNames: make([]string, 900, 900)}
	r := strings.NewReader(text)
	parser := xml.NewDecoder(r)
	depth := 0
	var firstLevel, secondLevel string

	for {

		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			if depth == 1 {
				firstLevel = name
			} else if depth == 2 {
				secondLevel = name
			}
			depth++
		case xml.EndElement:
			if depth == 3 {
				secondLevel = ""
			} else if depth == 2 {
				firstLevel = ""
			}
			depth--
		case xml.CharData:
			bytes := xml.CharData(t)
			if depth == 3 && firstLevel == "userPermissions" && secondLevel == "name" {
				obj.addPerm(string(bytes))
			}
		default:
		}
	}

	//fmt.Println(obj)
	return obj
}

/////////////////////////////////////////////////////////

var cmdSecurity = &Command{
	Run:   runSecurity,
	Usage: "security [SObject]",
	Short: "Displays the OLS and FLS for a given SObject",
	Long: `
Displays the OLS and FLS for a given SObject

Security Options
  -o, -object    # Object
  -x, -exclude   # Exclude a profile
  -i, -include   # Include a profile

  Examples:

  force security -o Case

  force security -o Case -x Admin
`,
	MaxExpectedArgs: -1,
}

type stringList []string

func (i *stringList) String() string {
	return fmt.Sprint(*i)
}

func (i *stringList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func inList(name string, list stringList) bool {
	index := sort.SearchStrings(list, name)

	return index < len(list) && list[index] == name
}

// these names need to be unique across all cmd.Flag
var (
	securityShowWarnings   bool
	securityObjects		   stringList
	securityExcludeNames   stringList
	securityIncludeNames   stringList
	permNames					 stringList
)

func init() {
	cmdSecurity.Flag.Var(&securityObjects, "o", "object")
	cmdSecurity.Flag.Var(&securityObjects, "object", "object")
	cmdSecurity.Flag.Var(&securityExcludeNames, "x", "exclude profile name")
	cmdSecurity.Flag.Var(&securityExcludeNames, "exclude", "exclude profile name")
	cmdSecurity.Flag.Var(&securityIncludeNames, "i", "include only profile name")
	cmdSecurity.Flag.Var(&securityIncludeNames, "include", "include only profile name")
}

func runSecurity(cmd *Command, args []string) {
	// Get path from args if available
	wd, _ := os.Getwd()
	root := filepath.Join(wd, ".")

	sort.Strings(securityObjects)
	sort.Strings(securityExcludeNames)
	sort.Strings(securityIncludeNames)
	
	var query ForceMetadataQuery
	var sobjectName string

	if len(securityObjects) > 0 {
		sobjectName = securityObjects[0]
		query = ForceMetadataQuery{
			{Name: []string{"Profile"}, Members: []string{"*"}},
			{Name: []string{"CustomObject"}, Members: []string{sobjectName}},
		}
	} else {
		fmt.Println("specify an SObject name")
		return
	}

	force, _ := ActiveForce()

	// Step 1: retrieve the desired metadata
	files, problems, err := force.Metadata.Retrieve(query)
	if err != nil {
		ErrorAndExit(err.Error())
	}
	for _, problem := range problems {
		fmt.Fprintln(os.Stderr, problem)
	}

	// Step 2: go through the metadata and construct a list of Profile (profiles) and a CustomObject (theObject)
	var profiles list.List
	var theObject CustomObject

	if (true) {
		// new way sorted
		var profileListKeys = make([]string, 0, len(files))
		var profileMap = make(map[string]string)
		for name, data := range files {
			if strings.HasPrefix(name, "profiles/") {
				profileName := strings.TrimSuffix(strings.TrimPrefix(name, "profiles/"), ".profile")
				if (len(securityIncludeNames)>0 && inList(profileName, securityIncludeNames)) {
					profileListKeys = append(profileListKeys, profileName)
					profileMap[profileName] = string(data)
				} else if (len(securityExcludeNames)>0 && !inList(profileName, securityExcludeNames)) { 
					profileListKeys = append(profileListKeys, profileName)
					profileMap[profileName] = string(data)
				} else if (len(securityIncludeNames)==0 && len(securityExcludeNames)==0) {
					profileListKeys = append(profileListKeys, profileName)
					profileMap[profileName] = string(data)
				}
			} else if strings.HasPrefix(name, "objects/") {
				objectName := strings.TrimSuffix(strings.TrimPrefix(name, "objects/"), ".object")
				if objectName == sobjectName {
					theObject = parseCustomObjectXML(objectName, string(data))
				}
			}
		}		
		// sort the profiles
		sort.Strings(profileListKeys)
		for _, profileName := range profileListKeys {
			profiles.PushBack(parseProfileXML(profileName, string(profileMap[profileName])))
		}		
	} else {
		// old way non-sorted
		for name, data := range files {
			if strings.HasPrefix(name, "profiles/") {
				profileName := strings.TrimSuffix(strings.TrimPrefix(name, "profiles/"), ".profile")
				if (len(securityIncludeNames)>0 && inList(profileName, securityIncludeNames)) {
					profiles.PushBack(parseProfileXML(profileName, string(data)))
				} else if (len(securityExcludeNames)>0 && !inList(profileName, securityExcludeNames)) { 
					profiles.PushBack(parseProfileXML(profileName, string(data)))
				} else if (len(securityIncludeNames)==0 && len(securityExcludeNames)==0) {
					profiles.PushBack(parseProfileXML(profileName, string(data)))
				}
			} else if strings.HasPrefix(name, "objects/") {
				objectName := strings.TrimSuffix(strings.TrimPrefix(name, "objects/"), ".object")
				if objectName == sobjectName {
					theObject = parseCustomObjectXML(objectName, string(data))
				}
			}
		}
	}
	
	// Step 3: group the profiles that have the exact same OLS and FLS
	// for the desired object together
	allProfiles := map[string]list.List{}
	var p Profile
	var profileKeys list.List

	for e := profiles.Front(); e != nil; e = e.Next() {
		p = e.Value.(Profile)
		key := theObject.getProfileFootprint(p)
		tmpList, OK := allProfiles[key]
		if !OK {
			var tmpList2 list.List
			tmpList2.PushBack(p)
			allProfiles[key] = tmpList2
			profileKeys.PushBack(key)
		} else {
			tmpList.PushBack(p)
		}
	}
	//fmt.Printf("profileKeys: %s\n", profileKeys)
	// Step 4: generate an HTML file that shows the various groups of profiles
	// as well as their OLS and FLS
	HTMLoutput := "<html>" +
	"<head>" +
	"<link rel=\"stylesheet\" href=\"https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.9.0/styles/github.min.css\" />" +
	"<style>" +
	"*{font-family:sans-serif}.content-table{border-collapse:collapse;margin:25px 0;font-size:.9em;min-width:400px;border-radius:5px 5px 0 0;overflow:hidden;box-shadow:0 0 20px rgba(0,0,0,.15)}.content-table thead tr{background-color:#009879;color:#fff;text-align:left;font-weight:700}.content-table td,.content-table th{padding:12px 15px}.content-table tbody tr{border-bottom:1px solid #ddd}.content-table tbody tr:nth-of-type(even){background-color:#f3f3f3}.content-table tbody tr:last-of-type{border-bottom:2px solid #009879}.content-table tbody tr.active-row{font-weight:700;color:#009879}" +
	"</style>" +
	"</head>" +
	"<body style=\"text-align: center; font-family: 'Source Sans Pro', sans-serif\">" +
	"<table class=\"content-table\" border=\"1\" style=\"border-collapse:collapse;\">" +
	"<thead><tr><td>" + sobjectName + "</td>"

	for key := profileKeys.Front(); key != nil; key = key.Next() {
		val := allProfiles[key.Value.(string)]
		profileNames := ""
		for v := val.Front(); v != nil; v = v.Next() {
			if v.Value == nil {
				continue
			}
			theProfile := v.Value.(Profile)
			profileNames += theProfile.name + "<br/>"
		}
		HTMLoutput += "<td>" + strings.Replace(profileNames, " ", "&nbsp;", -1) + "</td>"
	}
	HTMLoutput += "</tr></thead><tbody>"

	OLSproperties := []string{"Allow Create", "Allow Read", "Allow Edit", "Allow Delete", "View All Records", "Modify All Records" }

	for _, OLSproperty := range OLSproperties {
		HTMLoutput += "  <tr><td>[Object] " + OLSproperty + "</td>"

		for key := profileKeys.Front(); key != nil; key = key.Next() {
			val := allProfiles[key.Value.(string)]
			theProfile := val.Front().Value.(Profile)
			theOLS := theProfile.objectPermissions[theObject.objectName]
			HTMLoutput += "<td>" + strings.Title(theOLS.getProperty(OLSproperty)) + "</td>"
		}
		HTMLoutput += "</tr>"
	}

	for idx := 0; idx < theObject.nbFields; idx++ {
		fieldName := theObject.fieldNames[idx]
		HTMLoutput += "<tr><td>" + fieldName + "</td>"
		for key := profileKeys.Front(); key != nil; key = key.Next() {
			val := allProfiles[key.Value.(string)]
			theProfile := val.Front().Value.(Profile)
			if theProfile.fieldPermissions[theObject.objectName+"."+fieldName].editable == "true" {
				HTMLoutput += "<td>Editable</td>"
			} else if theProfile.fieldPermissions[theObject.objectName+"."+fieldName].readable == "true" {
				HTMLoutput += "<td>Readable</td>"
			} else {
				HTMLoutput += "<td>-</td>"
			}
		}
		HTMLoutput += "</tr>"
	}

	HTMLoutput += "</tbody></table></body></html>"

	// Last step: write the file on disk and display it inside a Web browser
	if err := ioutil.WriteFile(filepath.Join(root, sobjectName + ".html"), []byte(HTMLoutput), 0644); err != nil {
		ErrorAndExit(err.Error())
	}

	desktop.Open(filepath.Join(root, sobjectName + ".html"))
}
