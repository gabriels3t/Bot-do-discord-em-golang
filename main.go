package main

import(
	"fmt"
	"os"
	"log"

	"github.com/bwmarrin/discordgo"
	
)
// Função para ler o arquivo token
func readTokenFromFile(filename string)(string, error){
	data, err := os.ReadFile(filename) // lendo o arquivo, nil == null 
	if err != nil { 
		return "", err
	}
	return string(data),nil
}

func messageCreate(s *discordgo.Session, message *discordgo.MessageCreate){
	// Ignorando mensagens feitas pelo proprio bot
	if message.Author.ID == s.State.User.ID{
		return
	}
	if message.Content == "!ping"{
		s.ChannelMessageSend(message.ChannelID,"Pong !")
	}
}

func main(){

	token, err := readTokenFromFile("data/token.txt")
	if err != nil{
		log.Fatal(err)
	}
	
	dg, err := discordgo.New("Bot "+token)
	if err != nil{
		fmt.Println("Erro ao criar o bot",err)
		return 
	}
	dg.AddHandler(messageCreate)

	err = dg.Open()

	if err != nil{
		fmt.Println("erro ao abrir conexão com o bot",err)
		return
	}
	defer dg.Close()

	fmt.Println("bot ativo")
	select {}
}