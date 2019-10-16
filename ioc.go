package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

func loginCompleter(in prompt.Document) []prompt.Suggest {
	return nil
}

func executor(t string) {
	switch t {
	case "login":
		s := prompt.Input("login: ", loginCompleter)
		fmt.Println(s)
		fmt.Print("password: ")
		q, err := terminal.ReadPassword(0)
		fmt.Println(string(q), err)

	case "list clusters":
		fmt.Println(Magenta("List of clusters"))
	case "bye":
		fallthrough
	case "exit":
		fallthrough
	case "quit":
		fmt.Println(Magenta("Quitting"))
		os.Exit(0)
	case "help":
		fmt.Println("HELP:\nexit\nquit")
	default:
		fmt.Println("Command not found")
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "login", Description: "provide login info"},
		{Text: "help", Description: "show help with all commands"},
		{Text: "exit", Description: "quit the application"},
		{Text: "quit", Description: "quit the application"},
		{Text: "bye", Description: "quit the application"},
		{Text: "list", Description: "list resources (clusters, configurations)"},
	}

	empty_s := []prompt.Suggest{}

	list_s := []prompt.Suggest{
		{Text: "clusters", Description: "show list of all clusters available"},
		{Text: "configurations", Description: "show list all configurations"},
	}

	blocks := strings.Split(in.TextBeforeCursor(), " ")

	if len(blocks) == 2 {
		switch blocks[0] {
		case "ls":
			fallthrough
		case "list":
			return prompt.FilterHasPrefix(list_s, blocks[1], true)
		default:
			return empty_s
		}
	}
	if in.GetWordBeforeCursor() == "" {
		return nil
	} else {
		return prompt.FilterHasPrefix(s, blocks[0], true)
	}
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	p := prompt.New(executor, completer)
	p.Run()
}
