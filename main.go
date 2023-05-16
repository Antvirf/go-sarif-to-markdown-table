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
			// "CVE ID",
			"Severity",
			// "Detailed explanation",
			"CVE Details & affected versions",
		})
		for _, result := range run.Results {
			ruleIndex := getOrCreateRule(
				*result.RuleID,
				sarif.Runs[0].Tool.Driver,
			)

			// Parse CVE ID (contained in the first set of square brackets in a string)
			// cve_id := regexp.MustCompile(`\[(.*?)\]`).FindString(*result.Message.Text)

			tw.AppendRow(table.Row{
				*result.RuleID,
				// cve_id,
				run.Tool.Driver.Rules[ruleIndex].Properties["security-severity"],
				// *run.Tool.Driver.Rules[rule_index].Help.Markdown,
				*result.Message.Text,
			})
		}
		// Sort by severity
		tw.SortBy([]table.SortBy{
			{Name: "Severity", Mode: table.DscNumeric},
		})

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

// Find a particular rule by ID
func getOrCreateRule(id string, driver *sarif.ToolComponent) uint {
	for i, r := range driver.Rules {
		if r.ID == id {
			return uint(i)
		}
	}
	return uint(len(driver.Rules) - 1)
}
