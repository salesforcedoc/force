package command

import (
	"container/list"
//	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ForceCLI/force/desktop"
	. "github.com/ForceCLI/force/error"
	. "github.com/ForceCLI/force/lib"
)

/////////////////////////////////////////////////////////

var cmdUserPermissions = &Command{
	Run:   runUserPermissions,
	Usage: "userPermissions",
	Short: "Displays the user permissions for Profiles",
	Long: `
Displays the user permissions for Profiles

Security Options
  -x, -exclude   # Exclude a profile
  -i, -include   # Include a profile

  Examples:

  force userPermissions -i Admin

  force userPermissions -x Admin
`,
	MaxExpectedArgs: -1,
}

// these names need to be unique across all cmd.Flag
var (
	userPermssionsShowWarnings   bool
	userPermssionsExcludeNames   stringList
	userPermssionsIncludeNames   stringList
)


func init() {
	cmdUserPermissions.Flag.Var(&userPermssionsExcludeNames, "x", "exclude profile name")
	cmdUserPermissions.Flag.Var(&userPermssionsExcludeNames, "exclude", "exclude profile name")
	cmdUserPermissions.Flag.Var(&userPermssionsIncludeNames, "i", "include only profile name")
	cmdUserPermissions.Flag.Var(&userPermssionsIncludeNames, "include", "include only profile name")
}

func runUserPermissions(cmd *Command, args []string) {
	// Get path from args if available
	wd, _ := os.Getwd()
	root := filepath.Join(wd, ".")

	sort.Strings(userPermssionsExcludeNames)
	sort.Strings(userPermssionsIncludeNames)
	
	var query ForceMetadataQuery
	query = ForceMetadataQuery{
		{Name: []string{"Profile"}, Members: []string{"*"}},
		{Name: []string{"CustomObject"}, Members: []string{"Account"}},
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

	// Step 2: go through the metadata and construct a list of Profile (profiles) and a ProfileObject (theObject)
	var profiles list.List
	theObject := ProfileObject{profileName: "profileName", nbPerms: 0, permNames: make([]string, 900, 900)}

	if (true) {
		// new way sorted
		var profileListKeys = make([]string, 0, len(files))
		var profileMap = make(map[string]string)
		for name, data := range files {
			if strings.HasPrefix(name, "profiles/") {
				profileName := strings.TrimSuffix(strings.TrimPrefix(name, "profiles/"), ".profile")
				if (len(userPermssionsIncludeNames)>0 && inList(profileName, userPermssionsIncludeNames)) {
					profileListKeys = append(profileListKeys, profileName)
					profileMap[profileName] = string(data)
					theObject = parseProfileObjectXML(profileName, string(data), theObject)
				} else if (len(userPermssionsExcludeNames)>0 && !inList(profileName, userPermssionsExcludeNames)) { 
					profileListKeys = append(profileListKeys, profileName)
					profileMap[profileName] = string(data)
					theObject = parseProfileObjectXML(profileName, string(data), theObject)
				} else if (len(userPermssionsIncludeNames)==0 && len(userPermssionsExcludeNames)==0) {
					profileListKeys = append(profileListKeys, profileName)
					profileMap[profileName] = string(data)
					theObject = parseProfileObjectXML(profileName, string(data), theObject)
				}
			}
		}		
		// sort the profiles
		sort.Strings(profileListKeys)
		for _, profileName := range profileListKeys {
			profiles.PushBack(parseProfileXML(profileName, string(profileMap[profileName])))
		}		
	} 
	
	// Step 3: group the profiles that have the exact same user permissions
	// for the desired object together
	allProfiles := map[string]list.List{}
	var p Profile
	var profileKeys list.List

	for e := profiles.Front(); e != nil; e = e.Next() {
		p = e.Value.(Profile)
		key := theObject.getFootprint(p)
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
	"<thead><tr><td></td>"

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

	sort.Strings(permNames)
	for idx := 0; idx < theObject.nbPerms; idx++ {
		permName := theObject.permNames[idx]
		permName = permNames[idx]

		//fmt.Printf("permName: %d %d %s, %s\n", theObject.nbPerms, idx, permNames[idx], permName)
		HTMLoutput += "<tr><td>" + permName + "</td>"
		for key := profileKeys.Front(); key != nil; key = key.Next() {
			val := allProfiles[key.Value.(string)]
			theProfile := val.Front().Value.(Profile)
			if theProfile.userPermissions[permName].enabled == "true" {
				HTMLoutput += "<td>True</td>"
			} else {
				HTMLoutput += "<td>-</td>"
			}
		}
		HTMLoutput += "</tr>"
	}
	creds, err := ActiveCredentials(true)
	if err != nil {
		ErrorAndExit(err.Error())
	}
	HTMLoutput += "</tbody></table>Source: " + creds.InstanceUrl + 
	"<br/>Created: " + time.Now().Format(time.RFC850) + "<br/></body></html>"

	// Last step: write the file on disk and display it inside a Web browser
	if err := ioutil.WriteFile(filepath.Join(root, "UserPermissions_" + creds.SessionName() + ".html"), []byte(HTMLoutput), 0644); err != nil {
		ErrorAndExit(err.Error())
	}

	desktop.Open(filepath.Join(root, "UserPermissions_" + creds.SessionName() + ".html"))
}
