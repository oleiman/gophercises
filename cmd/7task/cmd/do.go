/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gophercises/cmd/7task/pkg/database"
	"strconv"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Do a task",
	Long: `Mark the given task as complete.

usage: 7task do <INDEX>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listIdx, err := strconv.Atoi(args[0])
		listIdx--
		if err != nil {
			return err
		}
		tasks, err := database.GetTasks(func(t database.Task) bool { return !t.Done })
		if len(tasks) <= listIdx {
			return errors.New("Invalid task ID")
		}
		tasks[listIdx].Done = true
		if err := database.UpdateTask(&tasks[listIdx]); err != nil {
			return err
		}
		fmt.Printf("You have completed the \"%s\" task. Lazy bones.\n", tasks[listIdx].Content)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
