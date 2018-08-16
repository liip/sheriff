package sheriff_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/liip/sheriff"
)

type UserType struct {
	Username string   `json:"username" type:"api"`
	Email    string   `json:"email" type:"personal"`
	Name     string   `json:"name" type:"api"`
	Roles    []string `json:"roles" type:"api" since:"2"`
}

type UserTypeList []UserType

func MarshalUserTypes(version *version.Version, groups []string, users UserTypeList) ([]byte, error) {
	o := &sheriff.Options{
		Groups:     groups,
		GroupName:  "type",
		ApiVersion: version,
	}

	data, err := sheriff.Marshal(o, users)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(data, "", "  ")
}

func ExampleType() {
	users := UserTypeList{
		UserType{
			Username: "alice",
			Email:    "alice@example.org",
			Name:     "Alice",
			Roles:    []string{"user", "admin"},
		},
		UserType{
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

	output, err := MarshalUserTypes(v1, []string{"api"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 1 output:")
	fmt.Printf("%s\n\n", output)

	output, err = MarshalUserTypes(v2, []string{"api"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 2 output:")
	fmt.Printf("%s\n\n", output)

	output, err = MarshalUserTypes(v2, []string{"api", "personal"}, users)
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
