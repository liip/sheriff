# sheriff
--
    import "github.com/mweibel/sheriff"

Package sheriff transforms structs into a map based on specific tags on the
struct fields. A typical use is an API which marshals structs into JSON and
maintains different API versions. Using sheriff, struct fields can be annotated
with API version and group tags. By invoking sheriff with specific options,
those tags determine whether a field will be added to the output map or not. It
can then be marshalled using "encoding/json".

Example:
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/hashicorp/go-version"
    "github.com/mweibel/sheriff"
)

type User struct {
    Username string   `json:"username" groups:"api"`
    Email    string   `json:"email" groups:"personal"`
    Name     string   `json:"name" groups:"api"`
    Roles    []string `json:"roles" groups:"api" since:"2"`
}

func (u User) Marshal(options *sheriff.Options) (interface{}, error) {
    return sheriff.Marshal(options, u)
}

type UserList []User

func (l UserList) Marshal(options *sheriff.Options) (interface{}, error) {
    list := make([]interface{}, len(l))
    for i, item := range l {
        target, err := item.Marshal(options)
        if err != nil {
            return nil, err
        }
        list[i] = target
    }
    return list, nil
}

func MarshalUsers(version *version.Version, groups []string, users UserList) ([]byte, error) {
    o := &sheriff.Options{
       Groups:     groups,
       ApiVersion: version,
    }

    data, err := users.Marshal(o)
    if err != nil {
        return nil, err
    }

    return json.MarshalIndent(data, "", "  ")
}

func main() {
    users := UserList{
        User{
            Username: "alice",
            Email:    "alice@example.org",
            Name:     "Alice",
            Roles:    []string{"user", "admin"},
        },
        User{
            Username: "bob",
            Email:    "bob@example.org",
            Name:     "Bob",
            Roles:    []string{"user"},
        },
    }

    v1, err := version.NewVersion("1.0.0")
    if err != nil {
        log.Panic(err)
    }
    v2, err := version.NewVersion("2.0.0")

    output, err := MarshalUsers(v1, []string{"api"}, users)
    if err != nil {
        log.Panic(err)
    }
    fmt.Println("Version 1 output:")
    fmt.Printf("%s\n\n", output)

    output, err = MarshalUsers(v2, []string{"api"}, users)
    if err != nil {
        log.Panic(err)
    }
    fmt.Println("Version 2 output:")
    fmt.Printf("%s\n\n", output)

    output, err = MarshalUsers(v2, []string{"api", "personal"}, users)
    if err != nil {
        log.Panic(err)
    }
    fmt.Println("Version 2 output with personal group too:")
    fmt.Printf("%s\n\n", output)

    // Output:
    // Version 1 output:
    // [
    //   {
    //     "name": "Alice",
    //     "username": "alice"
    //   },
    //   {
    //     "name": "Bob",
    //     "username": "bob"
    //   }
    // ]
    //
    // Version 2 output:
    // [
    //   {
    //     "name": "Alice",
    //     "roles": [
    //       "user",
    //       "admin"
    //     ],
    //     "username": "alice"
    //   },
    //   {
    //     "name": "Bob",
    //     "roles": [
    //       "user"
    //     ],
    //     "username": "bob"
    //   }
    // ]
    //
    // Version 2 output with personal group too:
    // [
    //   {
    //     "email": "alice@example.org",
    //     "name": "Alice",
    //     "roles": [
    //       "user",
    //       "admin"
    //     ],
    //     "username": "alice"
    //   },
    //   {
    //     "email": "bob@example.org",
    //     "name": "Bob",
    //     "roles": [
    //       "user"
    //     ],
    //     "username": "bob"
    //   }
    // ]
}
```

## Benchmarks

There's a simple benchmark in `bench_test.go` which compares running sheriff -> JSON versus just marshalling into JSON 
and runs on every build. Just marshalling JSON itself takes usually between 3 and 5 times less nanoseconds per operation
compared to running sheriff and JSON.

Want to make sheriff faster? Please send us your pull request or open an issue discussing a possible improvement 🚀!

## Usage

#### func  Marshal

```go
func Marshal(options *Options, data interface{}) (map[string]interface{}, error)
```
Marshal encodes the passed data into a map which can be used to pass to
json.Marshal().

The argument `data` needs to be a struct. If no struct is passed in, a
MarshalInvalidTypeError is returned.

#### type MarshalInvalidTypeError

```go
type MarshalInvalidTypeError struct {
}
```

MarshalInvalidTypeError is an error returned to indicate the wrong type has been
passed to Marshal.

#### func (MarshalInvalidTypeError) Error

```go
func (e MarshalInvalidTypeError) Error() string
```

#### type Marshaller

```go
type Marshaller interface {
	Marshal(options *Options) (interface{}, error)
}
```

Marshaller is the interface models have to implement in order to conform to
marshalling.

#### type Options

```go
type Options struct {
	// Groups determine which fields are getting marshalled based on the groups tag.
	// A field with multiple groups (comma-separated) will result in marshalling of that
	// field if one of their groups is specified.
	Groups []string
	// ApiVersion sets the API version to use when marshalling.
	// The tags `since` and `until` use the API version setting.
	// Specifying the API version as "1.0.0" and having an until setting of "2"
	// will result in the field being marshalled.
	// Specifying a since setting of "2" with the same API version specified,
	// will not marshal the field.
	ApiVersion *version.Version
}
```

Options determine which struct fields are being added to the output map.
