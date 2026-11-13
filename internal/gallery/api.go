package gallery

// apiField is one exported field of a library struct.
type apiField struct {
	Name, Type, Doc string
}

// apiType is a library struct as its source declares it.
type apiType struct {
	Doc    string
	Fields []apiField
}
