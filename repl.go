package repl

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Command[S any] interface {
	Name() string
	Description() string
	ArgHelp() string
	SetupFlags(*flag.FlagSet)
	Execute(S, *flag.FlagSet) error
}

func Run[S any](state S, comms []Command[S], promptPrefix func(S) string) error {
	r := bufio.NewReader(os.Stdin)
	sort.Slice(comms, func(i, j int) bool {
		return comms[i].Name() < comms[j].Name()
	})
	for {
		fmt.Printf("%s> ", promptPrefix(state))
		l, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		args, err := split(strings.TrimRight(l, "\r\n"))
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		help := func() {
			fmt.Printf("Commands:\n")
			fmt.Printf("\tq    - Quit.\n")
			fmt.Printf("\thelp - Display this help.\n")
			fmt.Printf("\n")
			var max int
			for _, c := range comms {
				if len(c.Name()) > max {
					max = len(c.Name())
				}
			}
			for _, c := range comms {
				fmt.Printf("\t%s ", c.Name())
				for i := 0; i < max-len(c.Name()); i++ {
					fmt.Printf(" ")
				}
				fmt.Printf("- %s\n", c.Description())
			}
			fmt.Printf("\n")
			fmt.Printf("Run <COMMAND> -help to see help specific to that command.\n")
		}
		switch args[0] {
		case "q":
			return nil
		case "help":
			help()
			continue
		}
		var found bool
		for _, c := range comms {
			if args[0] == c.Name() {
				found = true
				fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
				flagHelp := fs.Bool("help", false, "Display help for the command.")
				fs.SetOutput(os.Stdout)
				fs.Usage = func() {
					fmt.Printf("Usage: %s [FLAGS]", c.Name())
					if ah := c.ArgHelp(); ah != "" {
						fmt.Printf(" %s", ah)
					}
					fmt.Printf("\n")
					fmt.Printf("  %s\n", c.Description())
					fmt.Printf("Flags:\n")
					fs.PrintDefaults()
				}
				c.SetupFlags(fs)
				if err := fs.Parse(args[1:]); err != nil {
					fmt.Printf("error: %s\n", err)
					continue
				}
				if *flagHelp {
					fs.Usage()
					continue
				}
				if err := c.Execute(state, fs); err != nil {
					fmt.Printf("error: %s\n", err)
				}
			}
		}
		if !found {
			help()
		}
	}
}
