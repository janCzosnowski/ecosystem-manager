package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type Binding struct {
	Event   string `json:"event"`
	From    string `json:"from"`
	To      string `json:"to"`
	Handler string `json:"handler"`
}
var bindingsListPath string
var systemsListPath string



func addBinding(Event, From, To, Handler string) error {
	bindings, err := loadBindings()
	if err != nil {
		fmt.Println("Error loading bindings", err)
		return err
	}
	bindings = append(bindings, Binding{Event: Event, From: From, To: To, Handler: Handler})

	newFile, err := json.Marshal(bindings)
	os.WriteFile(bindingsListPath, newFile, 0644)
	return nil
}

func saveBindings(path string, bindings []Binding) error {
	data, err := json.MarshalIndent(bindings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func removeBindingInteractive(event, from, to, handler string) error {
	// Filter based on partial fields
	bindings, err := loadBindings()
	if err != nil {
		fmt.Println("error loading bindings:", err)
		return err
	}
	var matches []Binding
	var matchIndexes []int
	for i, b := range bindings {
		if (from == "0" || b.From == from) &&
			(to == "0" || b.To == to) &&
			(handler == "0" || b.Handler == handler) &&
			(event == "0" || b.Event == event) {
			matches = append(matches, b)
			matchIndexes = append(matchIndexes, i)
		}
	}

	if len(matches) == 0 {
		fmt.Println("No matching bindings found.")
		return nil
	}

	// Show matched bindings
	fmt.Printf("Found %d matching bindings:\n\n", len(matches))
	for i, b := range matches {
		fmt.Printf("[%d] Event: %s | From: %s | To: %s | Handler: %s\n", i+1, b.Event, b.From, b.To, b.Handler)
	}
	fmt.Print("\nRemove all [a], select by index [1 2 ...], or cancel [c]: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "c" || input == "" {
		fmt.Println("Cancelled.")
		return nil
	}

	var indexesToRemove map[int]bool = make(map[int]bool)

	if input == "a" {
		for _, i := range matchIndexes {
			indexesToRemove[i] = true
		}
	} else {
		parts := strings.Split(input, " ")
		for _, part := range parts {
			var idx int
			_, err := fmt.Sscanf(part, "%d", &idx)
			if err == nil && idx > 0 && idx <= len(matches) {
				indexesToRemove[matchIndexes[idx-1]] = true
			}
		}
	}

	var updated []Binding
	for i, b := range bindings {
		if !indexesToRemove[i] {
			updated = append(updated, b)
		}
	}
	fmt.Println(updated)
	saveBindings(bindingsListPath, updated)
	return nil
}

func loadBindings() ([]Binding, error) {
	data, err := os.ReadFile(bindingsListPath)
	if err != nil {
		return nil, err
	}

	var bindings []Binding
	err = json.Unmarshal(data, &bindings)
	if err != nil {
		return nil, err
	}
	return bindings, nil
}

func loadSystems() ([]string, error) {
	data, err := os.ReadFile(systemsListPath)
	if err != nil {
		return nil, err
	}

	var systems []string
	err = json.Unmarshal(data, &systems)
	if err != nil {
		return nil, err
	}

	return systems, nil
}

func addSystem(systemName string) error {
	systems, err := loadSystems()
	if err != nil {
		fmt.Println("Failed loading systems:", err)
		return err
	}
	systems = append(systems, systemName)
	newFile, err := json.MarshalIndent(systems, "", " ")
	if err != nil {
		fmt.Println("Code error, failed to Marshall", err)
		return err
	}
	err = os.WriteFile(systemsListPath, newFile, 0644)
	if err != nil {
		fmt.Println("Error saving /.config/ecosystem-manager/systems.json:", err)
		return err
	}
	return nil
}

func removeSystem(systemName string) error {
	systems, err := loadSystems()
	if err != nil {
		fmt.Println("Failed loading systems:", err)
		return err
	}
    for i, v := range systems {
        if v == systemName {
            systems = append(systems[:i], systems[i+1:]...)
        }
    }

	newFile, err := json.MarshalIndent(systems, "", " ")
	if err != nil {
		fmt.Println("Code error, failed to Marshall", err)
		return err
	}
	err = os.WriteFile(systemsListPath, newFile, 0644)
	if err != nil {
		fmt.Println("Error saving /.config/ecosystem-manager/systems.json:", err)
		return err
	}
	return nil
}

func runSystem(systemName string, systemArgs []string) {

	systems, err := loadSystems()
	if err != nil {
		fmt.Println("Error loading systems:", err)
		return
	}
	var contains = false
    for _, v := range systems {
        if v == systemName {
            contains = true
        }
    }
    

	if !contains {
		fmt.Printf("Requested system '%s' not added to /.config/ecosystem-manager/systems.json\n", systemName)
		return
	}

	

	cmdToRun := exec.Command(systemName, systemArgs...)
	cmdToRun.Stdout = os.Stdout
	cmdToRun.Stderr = os.Stderr

	err = cmdToRun.Run()
	if err != nil {
		fmt.Println("Error running system:", err)
	}
}

func main() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Couldn't load config directory")
	}
	bindingsListPath = configDir + "/ecosystem-manager/bindings.json"
	systemsListPath = configDir + "/ecosystem-manager/systems.json"
	os.MkdirAll(configDir + "/ecosystem-manager", 0755)
	_, err = os.Stat(bindingsListPath)
	if os.IsNotExist(err) {
		f, err := os.Create(bindingsListPath)
		if(err != nil){
			fmt.Println("Error creating a bindings list file at", bindingsListPath + ":", err)
		}
		defer f.Close()
		_, err = f.WriteString("[]")
		if err != nil {
			fmt.Println("Error creating a bindings list file at", bindingsListPath + ":", err)
		}
	} else if err != nil {
		fmt.Println("Error creating a bindings list file at", bindingsListPath + ":", err)
	}
	_, err = os.Stat(systemsListPath)
	if os.IsNotExist(err) {
		f, err := os.Create(systemsListPath)
		if(err != nil){
			fmt.Println("Error creating a systems list file at", systemsListPath + ":", err)
		}
		defer f.Close()
		_, err = f.WriteString("[]")
		if err != nil {
			fmt.Println("Error creating a systems list file at", systemsListPath + ":", err)
		}
	} else if err != nil {
		fmt.Println("Error creating a systems list file at", systemsListPath + ":", err)
	}
	
	var rootCmd = &cobra.Command{Use: "eco"}

	var cmdRun = &cobra.Command{
		Use:   "run [system] args",
		Short: "Run a specified system",
		Long:  "Run a specified system from your list of systems defined in /.config/ecosystem-manager/systems.json",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Please specify a system to run.")
				return
			}

			systemName := args[0]
			systemArgs := args[1:]
			runSystem(systemName, systemArgs)
		},
	}

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "List added systems",
		Long:  "List systems from the list of systems defined in /.config/ecosystem-manager/systems.json",
		Run: func(cmd *cobra.Command, args []string) {
			systems, err := loadSystems()
			if err != nil {
				fmt.Println("error loading systems:", err)
				return
			}
			systemList := ""
			for _, key := range systems {
				systemList = systemList + "\n" + key
			}
			fmt.Printf("Systems:%v", systemList)
		},
	}
	var cmdRemoveSystem = &cobra.Command{
		Use:   "remove [system name]",
		Short: "Removes a system",
		Long:  "Removes a system from $HOME/.config/ecosystem-manager/systems.json",
		Run: func(cmd *cobra.Command, args []string) {
			systemName := args[0]
			err := removeSystem(systemName)
			if err == nil {
				fmt.Println("System removed successfully")
			}
		},
	}
	var cmdAddSystem = &cobra.Command{
		Use:   "add [system name] [system main executable path]",
		Short: "Adds a system",
		Long:  "Adds a system to $HOME/.config/ecosystem-manager/systems.json",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Println("Please insert 2 arguments: name and path")
			}

			systemName := args[0]
			if err != nil {
				fmt.Println("Incorrect filepath format: ", err)
				return
			}

			err = addSystem(systemName)
			if err == nil {
				fmt.Println("System added successfully")
			}
		},
	}
	var cmdRemoveBinding = &cobra.Command{
		Use:   "unbind <event|0> <from|0> <to|0> <handler|0>",
		Short: "Removes bindings with specified filter(0 if none)",
		Long:  "Removes bindings with specified filter(0 if none) from $HOME/.config/ecosystem-manager/bindings.json",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 4 {
				fmt.Println("Please insert 4 arguments")
				return
			}
			err := removeBindingInteractive(args[0], args[1], args[2], args[3])
			if err == nil {
				fmt.Println("Success")
			}
		},
	}
	var cmdAddBinding = &cobra.Command{
		Use:   "bind <event> <from> <to> <handler>",
		Short: "Creates a Bind",
		Long:  "Creates a Bind, they are stored in bindings.json and allow you to automatise interaction between different systems",
		Run: func(cmd *cobra.Command, args []string) {
			err := addBinding(args[0], args[1], args[2], args[3])
			if err == nil {
				fmt.Println("Binding added successfully")
			}
		},
	}
	var cmdEmit = &cobra.Command{
		Use:   "emit <event> <from> [args]",
		Short: "Emits an event",

		Run: func(cmd *cobra.Command, args []string) {
			// find correct handler, run it with args.
			bindings, err := loadBindings()
			if err != nil {
				fmt.Println("Error loading bindings:", err)
				return
			}
			var handlers, recievers []string

			for i := range bindings {
				if bindings[i].Event == args[0] && bindings[i].From == args[1] {
					handlers = append(handlers, bindings[i].Handler)
					recievers = append(recievers, bindings[i].To)
				}
			}
			if handlers == nil || recievers == nil {
				fmt.Println("No matching binds were found")
				return
			}
			if len(handlers) != len(recievers) {
				fmt.Println("Error, handlers and recievers don't match")
				return
			}
			for i := range handlers {
				runSystem(recievers[i], append([]string{handlers[i]}, args[2:]...))
			}

		},
	}
	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdAddSystem)
	rootCmd.AddCommand(cmdRemoveSystem)
	rootCmd.AddCommand(cmdAddBinding)
	rootCmd.AddCommand(cmdRemoveBinding)
	rootCmd.AddCommand(cmdEmit)
	rootCmd.AddCommand(cmdList)
	rootCmd.Execute()
}
