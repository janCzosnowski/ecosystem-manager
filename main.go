package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func loadSystems() (map[string]interface{}, error) {
	data, err := os.ReadFile("systems.json")
	if err != nil {
		return nil, err
	}

	var systems map[string]interface{}
	err = json.Unmarshal(data, &systems)
	if err != nil {
		return nil, err
	}

	return systems, nil
}

func runSystem(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify a system to run.")
		return
	}

	systemName := args[0]
	systemArgs := args[1:]
	systems, err := loadSystems()
	if err != nil {
		fmt.Println("Error loading systems:", err)
		return
	}
	if _, exists := systems[systemName]; !exists {
		fmt.Printf("System '%s' not found in systems.json\n", systemName)
		return
	}

	scriptPath := fmt.Sprintf("./systems/%s/main.sh", systemName)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		fmt.Printf("Error: The executable %s doesn't exist", scriptPath)
		return
	}

	cmdToRun := exec.Command(scriptPath, systemArgs...)
	cmdToRun.Stdout = os.Stdout
	cmdToRun.Stderr = os.Stderr

	fmt.Println("Running:", scriptPath, strings.Join(systemArgs, ""))
	err = cmdToRun.Run()
	if err != nil {
		fmt.Println("Error running system:", err)
	}
}

func main() {
	var rootCmd = &cobra.Command{Use: "eco"}

	var cmdRun = &cobra.Command{
		Use:   "run [system] args",
		Short: "Run a specified system",
		Long:  "Run a specified system from your list of systems defined in systems.json",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runSystem(cmd, args)
		},
	}

	rootCmd.AddCommand(cmdRun)
	rootCmd.Execute()
}
