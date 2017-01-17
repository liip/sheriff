/*
Package sheriff transforms structs into a map based on specific tags on the struct fields.
A typical use is an API which marshals structs into JSON and maintains different API versions.
Using sheriff, struct fields can be annotated with API version and group tags. By invoking
sheriff with specific options, those tags determine whether a field will be added to the output
map or not. It can then be marshalled using "encoding/json".

*/
package sheriff
