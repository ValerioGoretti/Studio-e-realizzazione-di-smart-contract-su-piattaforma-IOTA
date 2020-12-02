package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"wasp/tools/cluster"
)

func check(err error) {
	if err != nil {
		fmt.Printf("[cluster] Waspt error: %s. Exit...\n", err)
		os.Exit(1)
	}
}

func usage(globalFlags *flag.FlagSet) {
	fmt.Printf("Usage: %s [options] [init|start|gendksets]\n", os.Args[0])
	globalFlags.PrintDefaults()
	os.Exit(1)
}

func main() {
	globalFlags := flag.NewFlagSet("", flag.ExitOnError)
	configPath := globalFlags.String("config", ".", "Config path")
	dataPath := globalFlags.String("data", "cluster-data", "Data path")
	err := globalFlags.Parse(os.Args[1:])
	check(err)

	wasps, err := cluster.New(*configPath, *dataPath)
	check(err)

	if globalFlags.NArg() < 1 {
		usage(globalFlags)
	}

	switch globalFlags.Arg(0) {

	case "init":
		initFlags := flag.NewFlagSet("init", flag.ExitOnError)
		resetDataPath := initFlags.Bool("r", false, "Reset data path if it exists")
		err = initFlags.Parse(globalFlags.Args()[1:])
		check(err)
		err = wasps.Init(*resetDataPath, "init")
		check(err)

	case "start":
		err = wasps.Start()
		check(err)
		fmt.Printf("-----------------------------------------------------------------\n")
		fmt.Printf("           The cluster started\n")
		fmt.Printf("-----------------------------------------------------------------\n")

		waitCtrlC()
		wasps.Wait()

	case "gendksets":
		err = wasps.Start()
		check(err)
		fmt.Printf("-----------------------------------------------------------------\n")
		fmt.Printf("           Generate DKSets\n")
		fmt.Printf("-----------------------------------------------------------------\n")
		err = wasps.GenerateDKSetsToFile()
		check(err)
		wasps.Stop()

	default:
		usage(globalFlags)
	}
}

func waitCtrlC() {
	fmt.Printf("[waspt] Press CTRL-C to stop\n")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
