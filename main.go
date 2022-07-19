package main

import (
	"fmt"
	// "io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
	"strings"
	"syscall"

	"serve/serving"
)

func check(e error) {
    if e != nil {
        fmt.Println(e)
		os.Exit(3)
    }
}

func resolveArg(arg string) {
	switch arg {
	case "--help":
		help()
		break;
	case "?":
		help()
		break;
	case "--version":
		printVersion()
		break;
	case "-v":
		printVersion()
		break;
	default:
		help()
	}
}

func help() {
	help, err := os.ReadFile("./help.txt")
	check(err)
	fmt.Println(string(help))
	os.Exit(0)
}

func printVersion() {
	version, err := os.ReadFile("./version.txt")
	check(err)
	fmt.Println(string(version))
	os.Exit(0)
}

func main() {
	var args []string
	var Path []string

	valid_args := []string{"help", "?", "v", "version"}
	vals := os.Args

	// fmt.Println(vals)

	for i, arg := range vals {
		if arg[0] == '-' || arg[0] == '?' {
			args = append(args, arg)
		} else if i != 0 {
			Path = append(Path, arg)
		}
	}

	if len(Path) < 1 {
		ls, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		Path = append(Path, ls)
	}

	if len(Path) > 1 {
		panic("Only one Path is supported")
	}

	// fmt.Println(args)
	// fmt.Println(Path)
	// fmt.Println(valid_args)

	// fmt.Println(strings.Split("--help", "--")[1])

	for _, arg := range args {
		if !stringInSlice(strings.ReplaceAll(arg, "-", ""), valid_args) {
			fmt.Printf("Not a valid argument, '%s'", arg)
			help()
		} else {
			resolveArg(arg)
		}
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		err := serving.Serve(Path[0])
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print("Gracefully shutting down...")
		// Abruptly end if 'SIGTERM' if fired before full shut down
		go func() {
			sig := <-sigs
			fmt.Println()
			fmt.Println(sig)
			os.Exit(2)
		}()
		// Shutdown processes
		time.Sleep(1 * time.Second)
		fmt.Print("\rGracefully shut down       \n")
        done <- true
	}()

	<-done
	// Exiting program automatically
	// fmt.Println("exited")
}

// Helper Functions

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}