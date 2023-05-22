package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

func main() {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	sarif, err := sarif.FromBytes(content)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("## Security Vulnerabilities\n\n")

	// Loop through all the runs in the SARIF input
	for runIndex, run := range sarif.Runs {
		tw := table.NewWriter()
		tw.AppendHeader(table.Row{
			"Rule ID",
			"Severity",
			"Rule details",
			"Location",
		})
		for _, result := range run.Results {
			ruleIndex, ruleIdString := getOrCreateRule(
				*result.RuleID,
				sarif.Runs[0].Tool.Driver,
			)

			tw.AppendRow(table.Row{
				ruleIdString,
				nilToEmptyStringFilter(run.Tool.Driver.Rules[ruleIndex].Properties["security-severity"]),
				*result.Message.Text,
				nilToEmptyStringFilter(*result.Locations[0].PhysicalLocation.ArtifactLocation.URI),
			})
		}
		// Sort by severity
		tw.SortBy([]table.SortBy{
			{Name: "Severity", Mode: table.DscNumeric},
		})

		// Drop empty columns
		tw.SuppressEmptyColumns()

		// Loop through each run
		fmt.Printf("*Run %d: Scanned with [%s](%s)*\n\n", runIndex+1, run.Tool.Driver.Name, *run.Tool.Driver.InformationURI)
		if len(run.Results) > 0 {
			fmt.Print(tw.RenderMarkdown())
			fmt.Print("\n---\n")
		} else {
			fmt.Print("No security vulnerabilities found.\n")
		}
	}

}

// Find a particular rule by ID - return index and the ID, if helpUri exists, the ID is markdown link format
func getOrCreateRule(id string, driver *sarif.ToolComponent) (uint, string) {
	for i, r := range driver.Rules {
		if r.ID == id {
			if r.HelpURI != nil {
				return uint(i), fmt.Sprintf("[%s](%s)", id, *r.HelpURI)
			} else {
				return uint(i), fmt.Sprintf("%s", id)
			}
		}
	}
	return uint(len(driver.Rules) - 1), fmt.Sprintf("%s", id)
}

func nilToEmptyStringFilter(input interface{}) interface{} {
	if input == nil {
		return ""
	} else {
		return input
	}
}
