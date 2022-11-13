package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	pcstat "pcstat/pkg"
	"strings"
)

type PcStatusList []pcstat.PcStatus

func (stats PcStatusList) formatUnicode() {
	maxName := stats.maxNameLen()

	// create horizontal grid line
	pad := strings.Repeat("─", maxName+2)
	top := fmt.Sprintf("┌%s┬────────────────┬────────────┬───────────┬─────────┐", pad)
	hr := fmt.Sprintf("├%s┼────────────────┼────────────┼───────────┼─────────┤", pad)
	bot := fmt.Sprintf("└%s┴────────────────┴────────────┴───────────┴─────────┘", pad)

	fmt.Println(top)

	// -nohdr may be chosen to save 2 lines of precious vertical space
	if !nohdrFlag {
		pad = strings.Repeat(" ", maxName-4)
		fmt.Printf("│ Name%s │ Size (bytes)   │ Pages      │ Cached    │ Percent │\n", pad)
		fmt.Println(hr)
	}

	for _, pcs := range stats {
		pad = strings.Repeat(" ", maxName-len(pcs.Name))

		// %07.3f was chosen to make it easy to scan the percentages vertically
		// I tried a few different formats only this one kept the decimals aligned
		fmt.Printf("│ %s%s │ %-15d│ %-11d│ %-10d│ %07.3f │\n",
			pcs.Name, pad, pcs.Size, pcs.Pages, pcs.Cached, pcs.Percent)
	}

	fmt.Println(bot)
}

func (stats PcStatusList) formatText() {
	maxName := stats.maxNameLen()

	// create horizontal grid line
	pad := strings.Repeat("-", maxName+2)
	top := fmt.Sprintf("+%s+----------------+------------+-----------+---------+", pad)
	hr := fmt.Sprintf("|%s+----------------+------------+-----------+---------|", pad)
	bot := fmt.Sprintf("+%s+----------------+------------+-----------+---------+", pad)

	fmt.Println(top)

	// -nohdr may be chosen to save 2 lines of precious vertical space
	if !nohdrFlag {
		pad = strings.Repeat(" ", maxName-4)
		fmt.Printf("| Name%s | Size (bytes)   | Pages      | Cached    | Percent |\n", pad)
		fmt.Println(hr)
	}

	for _, pcs := range stats {
		pad = strings.Repeat(" ", maxName-len(pcs.Name))

		// %07.3f was chosen to make it easy to scan the percentages vertically
		// I tried a few different formats only this one kept the decimals aligned
		fmt.Printf("| %s%s | %-15d| %-11d| %-10d| %07.3f |\n",
			pcs.Name, pad, pcs.Size, pcs.Pages, pcs.Cached, pcs.Percent)
	}

	fmt.Println(bot)
}

func (stats PcStatusList) formatPlain() {
	maxName := stats.maxNameLen()

	// -nohdr may be chosen to save 2 lines of precious vertical space
	if !nohdrFlag {
		pad := strings.Repeat(" ", maxName-4)
		fmt.Printf("Name%s  Size (bytes)    Pages       Cached     Percent\n", pad)
	}

	for _, pcs := range stats {
		pad := strings.Repeat(" ", maxName-len(pcs.Name))

		// %07.3f was chosen to make it easy to scan the percentages vertically
		// I tried a few different formats only this one kept the decimals aligned
		fmt.Printf("%s%s  %-15d %-11d %-10d %07.3f\n",
			pcs.Name, pad, pcs.Size, pcs.Pages, pcs.Cached, pcs.Percent)
	}
}

func (stats PcStatusList) formatTerse() {
	if !nohdrFlag {
		fmt.Println("name,size,timestamp,mtime,pages,cached,percent")
	}
	for _, pcs := range stats {
		time := pcs.Timestamp.Unix()
		mtime := pcs.Mtime.Unix()
		fmt.Printf("%s,%d,%d,%d,%d,%d,%g\n",
			pcs.Name, pcs.Size, time, mtime, pcs.Pages, pcs.Cached, pcs.Percent)
	}
}

func (stats PcStatusList) formatJson(clearpps bool) {
	// clear the per-page status when requested
	// emits an empty "status": [] field in the JSON when disabled, but NBD.
	if clearpps {
		for i := range stats {
			stats[i].PPStat = nil
		}
	}

	b, err := json.Marshal(stats)
	if err != nil {
		log.Fatalf("JSON formatting failed: %s\n", err)
	}
	os.Stdout.Write(b)
	fmt.Println("")
}

// references:
// http://www.unicode.org/charts/PDF/U2580.pdf
// https://github.com/puppetlabs/mcollective-puppet-agent/blob/master/application/puppet.rb#L143
// https://github.com/holman/spark
func (stats PcStatusList) formatHistogram() {
	ws := getwinsize()
	maxName := stats.maxNameLen()
	// fmt.Printf("wc %+v", ws)

	// block elements are wider than characters, so only use 1/2 the available columns
	buckets := (int(ws.ws_col)-maxName)/2 - 10

	for _, pcs := range stats {
		pad := strings.Repeat(" ", maxName-len(pcs.Name))
		fmt.Printf("%s%s %8d ", pcs.Name, pad, pcs.Pages)

		// when there is enough room display on/off for every page
		if buckets > pcs.Pages {
			for _, v := range pcs.PPStat {
				if v {
					fmt.Print("\u2588") // full block = 100%
					// fmt.Print("a")
				} else {
					fmt.Print("\u2581") // lower 1/8 block
				}
			}
		} else {
			// maybe 400/10
			bsz := pcs.Pages / buckets
			fbsz := float64(bsz)
			total := 0.0
			for i, v := range pcs.PPStat {
				if v {
					total++
				}

				// ignore some data to show, only show buckets number's data
				if (i+1)%bsz == 0 {
					// maybe total=20, fbsz = 40
					avg := total / fbsz
					if total == 0 {
						fmt.Print("\u2581") // lower 1/8 block = 0
					} else if avg < 0.16 {
						fmt.Print("\u2582") // lower 2/8 block for width
					} else if avg < 0.33 {
						fmt.Print("\u2583") // lower 3/8 block for width
					} else if avg < 0.50 {
						fmt.Print("\u2584") // lower 4/8 block for width
					} else if avg < 0.66 {
						fmt.Print("\u2585") // lower 5/8 block for width
					} else if avg < 0.83 {
						fmt.Print("\u2586") // lower 6/8 block for width
					} else if avg < 1.00 {
						fmt.Print("\u2587") // lower 7/8 block for width
					} else {
						fmt.Print("\u2588") // full block = 100% for width
					}

					total = 0
				}
			}
		}
		fmt.Println("")
	}
}

// maxNameLen returns the len of longest filename in the stat list
// if the bnameFlag is set, this will return the max basename len
func (stats PcStatusList) maxNameLen() int {
	var maxName int
	for _, pcs := range stats {
		if len(pcs.Name) > maxName {
			maxName = len(pcs.Name)
		}
	}

	if maxName < 5 {
		maxName = 5
	}
	return maxName
}
