# salesforcedoc customizations

### 04/01/22
- update security.go to support inclusion/exclusion of profiles
- update security.go to allow for sorting of profile names in output
- update security.go to output file as <objectname>.html
- rebased with https://github.com/ForceCLI/force.git

### 03/16/22
- update folder.go to support email template folders
- rebased with https://github.com/ForceCLI/force.git

### 03/04/22
- add export_check.go to export each metadata type individually to find problematic metadata types
- update apiversion.go to 50.0

### 12/02/20
- update apiversion.go to 48.0
### 11/04/20
- update export.go to include standard object PersonAccount

### 10/22/20
- update export.go to include managed package metadata types with -p
- update export.go to export only specific metadata types with -i
- update export.go exclude all __ChangeEvent objects for CustomObjects
- update export.go with additional metadata types for v50.0 (including NamedCredential)
- update export.go to include standard object StandardValueSets

### 10/28/20
- enhance export.go to include email templates under unfiled$public

### manifest
Similar to export but only generates the package.xml

### Build Notes
go get . && rm -rf src metadata && force export -i NetworkBranding

### build for macos-arm64 future
#env GOOS=darwin GOARCH=arm64 go build -o force-macos-arm64 main.go

### build for macos-x64
env GOOS=darwin GOARCH=amd64 go build -o force-macos-x64 main.go
## build for windows-x64
env GOOS=windows GOARCH=amd64 go build -o force-windows-x64.exe main.go

### build for linux-x64
env GOOS=linux GOARCH=amd64 go build -o force-linux-x64 main.go

### install to macos bin path
cp -pR force-macos-x64 /usr/local/Cellar/go/1.13/bin/force
