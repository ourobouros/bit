// Package main is complete tool for the go command line
package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/chriswalz/bit/cmd"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"github.com/thoas/go-funk"
	"os"
	"strconv"
	"strings"
)

func main() {
	compLine := os.Getenv("COMP_LINE")
	compPoint := os.Getenv("COMP_POINT")
	i, err := strconv.Atoi(compPoint)
	if err != nil {
		fmt.Println("COMP_LINE", compLine, "COMP_POINT", compPoint, "err:", err)
		return
	}
	if i > len(compLine) {
		err := os.Setenv("COMP_POINT", strconv.Itoa(len(compLine)))
		if err != nil {
			fmt.Println(err)
		}
	}

	branchCompletion := &complete.Command{
		Args: complete.PredictFunc(func(prefix string) []string {
			branches := cmd.BranchListSuggestions()
			completion := make([]string, len(branches))
			for i, v := range branches {
				completion[i] = v.Text
			}
			return completion
		}),
	}

	cmds := cmd.AllBitAndGitSubCommands(cmd.ShellCmd)
	completionSubCmdMap := map[string]*complete.Command{}
	for _, v := range cmds {
		flagSuggestions := append(cmd.FlagSuggestionsForCommand(v.Name(), "--"), cmd.FlagSuggestionsForCommand(v.Name(), "-")...)
		flags := funk.Map(flagSuggestions, func(x prompt.Suggest) (string, complete.Predictor) {
			if strings.HasPrefix(x.Text, "--") {
				return x.Text[2:], predict.Nothing
			} else if strings.HasPrefix(x.Text, "-") {
				return x.Text[1:2], predict.Nothing
			} else {
				return "", predict.Nothing
			}
		}).(map[string]complete.Predictor)
		completionSubCmdMap[v.Name()] = &complete.Command{
			Flags: flags,
		}
		if v.Name() == "checkout" || v.Name() == "co" || v.Name() == "switch" || v.Name() == "pull" || v.Name() == "merge" {
			branchCompletion.Flags = flags
			completionSubCmdMap[v.Name()] = branchCompletion
		}
		if v.Name() == "release" {
			completionSubCmdMap[v.Name()].Sub = map[string]*complete.Command{
				"bump": {},
			}
		}
	}

	gogo := &complete.Command{
		Sub: completionSubCmdMap,
	}

	gogo.Complete("bit")
}
