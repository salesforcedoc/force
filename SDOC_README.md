# salesforcedoc customizations

### 10/22
- update export to include managed package metadata types with -p
- update export to extract only specific metadata types with -i
- update export exclude all __ChangeEvent objects for CustomObjects
- update export with additional metadata types for v50.0 (including NamedCredential)
- update export to include standard object StandardValueSets

### 10/28
- enhance export to include email templates under unfiled$public

### manifest
Similar to export but only generates the package.xml

### Build Notes
go get . && rm -rf src metadata && force export -i NetworkBranding

### build for macos-x64
env GOOS=darwin GOARCH=amd64 go build -o force-macos-x64 main.go

## build for windows-x64
env GOOS=windows GOARCH=amd64 go build -o force-windows-x64.exe main.go

### build for linux-x64
env GOOS=linux GOARCH=amd64 go build -o force-linux-x64 main.go
