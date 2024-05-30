package sheriff_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/liip/sheriff/v2"
)

type User struct {
	Username string   `json:"username" groups:"api"`
	Email    string   `json:"email" groups:"personal"`
	Name     string   `json:"name" groups:"api"`
	Roles    []string `json:"roles" groups:"api" since:"2"`
}

type UserList []User

func MarshalUsers(version *version.Version, groups []string, users UserList) ([]byte, error) {
	o := &sheriff.Options{
		Groups:     groups,
		ApiVersion: version,
	}

	data, err := sheriff.Marshal(o, users)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(data, "", "  ")
}

func Example() {
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
	if err != nil {
		log.Panic(err)
	}

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
	//     "username": "alice",
	//     "name": "Alice"
	//   },
	//   {
	//     "username": "bob",
	//     "name": "Bob"
	//   }
	// ]
	//
	// Version 2 output:
	// [
	//   {
	//     "username": "alice",
	//     "name": "Alice",
	//     "roles": [
	//       "user",
	//       "admin"
	//     ]
	//   },
	//   {
	//     "username": "bob",
	//     "name": "Bob",
	//     "roles": [
	//       "user"
	//     ]
	//   }
	// ]
	//
	// Version 2 output with personal group too:
	// [
	//   {
	//     "username": "alice",
	//     "email": "alice@example.org",
	//     "name": "Alice",
	//     "roles": [
	//       "user",
	//       "admin"
	//     ]
	//   },
	//   {
	//     "username": "bob",
	//     "email": "bob@example.org",
	//     "name": "Bob",
	//     "roles": [
	//       "user"
	//     ]
	//   }
	// ]
}
