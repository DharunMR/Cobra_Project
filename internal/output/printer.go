package output

import "fmt"

func PrintResult(name string, passed bool, err error) {

	if passed {
		fmt.Println("✔", name)
	} else {
		if err != nil {
			fmt.Printf("✖ %s → %v\n", name, err)
		} else {
			fmt.Printf("✖ %s\n", name)
		}
	}
}

func PrintSummary(passed, total int) {
	fmt.Printf("\n%d/%d tests passed\n", passed, total)
}
