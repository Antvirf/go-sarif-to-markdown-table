package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

func main() {
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	sarif, err := sarif.FromBytes(content)
	if err != nil {
		log.Fatal(err)
	}

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Rule ID", "CVE ID", "Severity", "Detailed explanation", "Fixed version or mitigation"})

	// loop through runs in the SARIF input, if multiple exist
	for _, run := range sarif.Runs {
		for _, result := range run.Results {
			rule_index := getOrCreateRule(
				*result.RuleID,
				sarif.Runs[0].Tool.Driver,
			)

			// Parse CVE ID (contained in the first set of square brackets in a string)
			cve_id := regexp.MustCompile(`\[(.*?)\]`).FindString(*result.Message.Text)

			tw.AppendRow(table.Row{
				*result.RuleID,
				cve_id,
				run.Tool.Driver.Rules[rule_index].Properties["security-severity"],
				*run.Tool.Driver.Rules[rule_index].Help.Markdown,
				*result.Message.Text,
			})
		}
	}

	// Sort by severity
	tw.SortBy([]table.SortBy{
		{Name: "Severity", Mode: table.Asc},
	})

	// Output
	fmt.Print(tw.RenderMarkdown())

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
