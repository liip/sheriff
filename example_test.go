package sheriff_test

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
