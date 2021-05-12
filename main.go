package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"path/filepath"
)

type User struct {
	Id					string	`csv:"id"`
	UserPrincipalName	string	`csv:"userPrincipalName"`
	DisplayName			string	`csv:"displayName"`
	ObjectType			string	`csv:"objectType"`
	UserType			string	`csv:"userType"`
	IsUser				bool	`csv:"isUser"`
	IsGroup				bool	`csv:"isGroup"`
	IsGuest				bool	`csv:"isGuest"`
}

type OutPutData struct {
	UserPrincipalName	string
	Group				string
}

var output []OutPutData = []OutPutData{}
var userstofind []User = []User{}
var userstocompare []User = []User{}
var group chan map[string][]User = make(chan map[string][]User)

func main() {
	GetUsersToFind("C:\\Users\\rachid.moyse\\OneDrive - Tivit\\Documentos\\VPN Groups\\Grupos a Revisar")
	GetUserstoCompare("C:\\Users\\rachid.moyse\\OneDrive - Tivit\\Documentos\\VPN Groups")

	groups, ok := <- group
	if !ok {
		log.Println("Ocurrio un error leyendo el canal!")
		return
	}

	for k, v := range groups {
		for _, user := range v {
			output = append(output, OutPutData{UserPrincipalName: user.UserPrincipalName,Group:k})
		}
	}

	cvsContent, err := gocsv.MarshalString(&output)
	if err != nil {
		log.Println("Error: ", err)
	}
	file, err := os.OpenFile("./outputdata.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Println("Error: ", err)
	}
	file.WriteString(cvsContent)
}

func GetUsersToFind(path string) {
	//Finding the CSV files in the path given
	path_file := path + "\\*.csv"
	files, err := filepath.Glob(path_file)
	if err != nil {
		log.Println("Error: ", err)
	}

	// Setting the Output File for logging
	f, err := os.OpenFile("./Log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	log.SetOutput(f)

	//Map for groups
	g := make(map[string][]User)

	//Parsing CSV Files to Slice of User Struct
	for _, file := range files {
		log.Println("Leyendo el archivo ", filepath.Base(file))
		clientfile, err := os.OpenFile(string(file), os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		defer clientfile.Close()
		if err := gocsv.UnmarshalFile(clientfile, &userstofind); err != nil {
			panic(err)
		}

		//Adding the group to the users
		g[filepath.Base(file)] = userstofind
	}
	go func() {
		group <- g
	}()
}

func GetUserstoCompare(path string) {
	//Finding the CSV files in the path given
	path_file := path + "\\*.csv"
	files, err := filepath.Glob(path_file)
	if err != nil {
		log.Println("Error: ", err)
	}

	// Setting the Output File for logging
	f, err := os.OpenFile("./Log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	log.SetOutput(f)

	//Adding Map
	g := make(map[string][]User)

	//Parsing CSV Files to Slice of User Struct
	for _, file := range files {
		log.Println("Leyendo el archivo ", filepath.Base(file))
		clientfile, err := os.OpenFile(string(file), os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		defer clientfile.Close()
		if err := gocsv.UnmarshalFile(clientfile, &userstocompare); err != nil {
			panic(err)
		}
		g[filepath.Base(file)] = userstofind
	}
	go func() {
		group <- g
	}()
}
