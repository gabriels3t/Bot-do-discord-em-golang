package main

// Importando libs 
import (
	"fmt"                   // Para operações de entrada/saída formatadas
	"os"                    // Para manipulação de arquivos
	"log"                   // Para registro de logs de erro
	"time"                  // Para operações relacionadas a tempo
	"github.com/bwmarrin/discordgo" // Biblioteca para interagir com a API do Discord
)

// Função para ler o arquivo de token e retornar seu conteúdo
func readTokenFromFile(filename string)(string, error){
	data, err := os.ReadFile(filename) // Lê o conteúdo do arquivo
	if err != nil { 
		return "", err // Se houver erro ao ler o arquivo, retorna o erro
	}
	return string(data), nil // Retorna o conteúdo do arquivo como string
}

// Função para conectar o bot a um canal de voz
func joinVoiceChannel(s *discordgo.Session, guildId string, channelID string) error {
	// Tenta conectar o bot ao canal de voz específico
	vc, err := s.ChannelVoiceJoin(guildId, channelID, false, false)
	if err != nil {
		return fmt.Errorf("Não foi possível se conectar ao canal de voz: %v", err)
	}

	// Assegura que o bot se desconecte do canal quando a função terminar
	defer vc.Disconnect()

	// Imprime no console que o bot entrou no canal de voz
	fmt.Println("Bot entrou no canal de voz", channelID)

	return nil
}

// Função que é chamada toda vez que uma mensagem é criada
func messageCreate(s *discordgo.Session, message *discordgo.MessageCreate){
	// Ignora mensagens do próprio bot
	if message.Author.ID == s.State.User.ID {
		return
	}

	// Responde ao comando "!ping"
	if message.Content == "!ping" {
		s.ChannelMessageSend(message.ChannelID, "Pong !")
	}

	// Responde ao comando "!time"
	if message.Content == "!time" {
		s.ChannelMessageSend(message.ChannelID, "Estou contando :)")

		guildID := message.GuildID // ID do servidor (guild)
		userID := message.Author.ID // ID do usuário que enviou a mensagem

		// Obtém o estado  do usuário (se ele está em um canal de voz)
		voiceState, err := s.State.VoiceState(guildID, userID)
		if err != nil {
			fmt.Println("Erro ao obter dados do membro:", err)
			return
		}

		// Verifica se o usuário está em um canal de voz
		if voiceState != nil && voiceState.ChannelID != "" {
			// Tenta conectar ao canal de voz do usuário
			err = joinVoiceChannel(s, guildID, voiceState.ChannelID)
			if err != nil {
				fmt.Println("Erro ao conectar ao canal de voz:", err)
				return
			}

			// Inicia uma goroutine para enviar uma mensagem depois de 1 minuto
			go func() {
				time.Sleep(1 * time.Minute)
				s.ChannelMessageSend(message.ChannelID, "Tempo acabou :(")
			}()

			// Envia uma mensagem indicando que o bot entrou no canal de voz
			s.ChannelMessageSend(message.ChannelID, "Entrei no canal de voz")
		} else {
			// Caso o usuário não esteja em um canal de voz
			s.ChannelMessageSend(message.ChannelID, "Por favor, entre em um canal de voz")
		}
	}
}

// Função principal do programa
func main() {
	// Lê o token do arquivo e verifica se houve erro
	token, err := readTokenFromFile("data/token.txt")
	if err != nil {
		log.Fatal(err) // Se ocorrer erro, o programa é finalizado
	}
	
	// Cria uma nova sessão do bot com o token lido
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Erro ao criar o bot:", err)
		return
	}
	
	// Adiciona o manipulador de mensagens para o bot (quando uma mensagem é criada, a função 'messageCreate' é chamada)
	dg.AddHandler(messageCreate)

	// Abre a conexão com o Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Erro ao abrir conexão com o bot:", err)
		return
	}

	// Garante que a conexão será fechada ao final do programa
	defer dg.Close()

	// Informa que o bot está ativo
	fmt.Println("Bot ativo")

	// Aguarda indefinidamente para manter o bot rodando
	select {}
}
