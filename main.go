package main

import(
	"fmt"
	"io/ioutil"
	"log"
	
)

func readTokenFromFile(filename string)(string, error){
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data),nil
}

func main(){

	token, err := readTokenFromFile("data/token.txt")

	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(token)


}