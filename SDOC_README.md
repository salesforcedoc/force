# salesforcedoc customizations

### export
- add ability to include managed package metadata types
- add ability to extract only specific metadata types
- exclude all __ChangeEvent objects for CustomObjects
- added NamedCredential as a metadata type

### manifest
Similar to export but only generates the package.xml

### Notes
go get . && rm -rf src metadata && force export -i NetworkBranding