package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/JeffreySmith/vmtools"
)

func main() {
	var OutputBuffer io.Writer = os.Stdout
	var InputBuffer io.Reader
	var header string
	stdin := os.Stdin
	f, err := stdin.Stat()

	ip := flag.String("ip", "", "Comma separated list of ip addresses.")
	output := flag.String("output", "", "Output file for generated yaml.")
	input := flag.String("input", "", "Input file for user names.")
	header_path := flag.String("header", "", "Path to a file containing your yaml file header (optional).")
	indentation_level := flag.Int("indent", 2, "Set the indentation level. Must be > 2")
	flag.Parse()
	if len(*ip) == 0 {
		fmt.Fprintf(os.Stderr, "You must supply at least 1 ip address\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	if f.Size() > 0 {
		InputBuffer = os.Stdin
	} else if len(*input) > 0 {
		var err error
		InputBuffer, err = os.Open(*input)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		InputBuffer, err = os.Open("users")
		if err != nil {
			fmt.Fprintf(os.Stderr, "You must suppy an input through stdin, a supplied file, or a 'users' file in this directory\n")
			os.Exit(1)
		}
	}

	ips := strings.Split(*ip, ",")

	if len(*output) > 0 {
		var err error
		OutputBuffer, err = os.Create(*output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if len(*header_path) > 0 {
		_, err := os.Stat(*header_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read file %v: %v\n", *header_path, err)
			os.Exit(1)
		}
		f, err := os.ReadFile(*header_path)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %v: \n", err)
			os.Exit(1)
		}
		header = string(f)
	} else {
		header = "---"
	}
	config := vmtools.NewConfig(vmtools.WithOutput(OutputBuffer),
		vmtools.WithInput(InputBuffer),
		vmtools.WithHeader(header),
		vmtools.SetIndent(*indentation_level),
	)

	config.GetUsers(ips)

	_, err = config.GenerateYaml()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.WriteYaml()
}