package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	AppName = "pcstat"
)

var (
	GitSHA    string
	BuildTime string
)

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Version = GitSHA + "-" + BuildTime
	app.Usage = "page cache stat"
	app.Writer = os.Stdout
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  "pid",
			Value: 0,
			Usage: "show all open maps for the given pid",
		},
		&cli.BoolFlag{
			Name:  "terse",
			Value: false,
			Usage: "show terse output",
		},
		&cli.BoolFlag{
			Name:  "nohdr",
			Value: false,
			Usage: "omit the header from terse &text output",
		},
		&cli.BoolFlag{
			Name:  "json",
			Value: false,
			Usage: "return data in JSON format",
		},
		&cli.BoolFlag{
			Name:  "unicode",
			Value: false,
			Usage: "return data with unicode box characters",
		},
		&cli.BoolFlag{
			Name:  "plain",
			Value: false,
			Usage: "return data with no box characters",
		},
		&cli.BoolFlag{
			Name:  "pps",
			Value: false,
			Usage: "include the per-page status in JSON output",
		},
		&cli.BoolFlag{
			Name:  "histo",
			Value: false,
			Usage: "print a simple histogram instead of raw data",
		},
		&cli.BoolFlag{
			Name:  "bname",
			Value: false,
			Usage: "covert paths to basename to narrow the output",
		},
		&cli.BoolFlag{
			Name:  "sort",
			Value: false,
			Usage: "sort output by cached pages desc",
		},
	}
	app.Action = func(c *cli.Context) error {
		pidFlag := c.Int("pid")
		terseFlag := c.Bool("terse")
		nohdrFlag := c.Bool("nohdr")
		jsonFlag := c.Bool("json")
		unicodeFlag := c.Bool("unicode")
		plainFlag := c.Bool("plian")
		ppsFlag := c.Bool("pps")
		histoFlag := c.Bool("histo")
		bnameFlag := c.Bool("bname")
		sortFlag := c.Bool("sort")
		files := c.Args().Slice()
		if pidFlag != 0 {
			//todo get and append to files
		}
		
		if len(files) == 0 {
			cli.ShowAppHelpAndExit(c, 1)
		}
		stats := make(PcStatusList, len(files))
		for _, fname := range files {
			// todo get status

			if bnameFlag {}
		}
		if sortFlag {
			sort.Slice(stats, func(i, j int) bool{
				return stats[i].Cached > stats[j].Cached
			})
		}

		if jsonFlag {
			stats.formatJson(!ppsFlag)
		} else if terseFlag {
			stats.formatTerse()
		} else if histoFlag {
			stats.formatHistogram()
		} else if unicodeFlag {
			stats.formatUnicode()
		} else if plainFlag {
			stats.formatPlain()
		} else {
			stats.formatText()
		}
		return nil
	}
	app.Run(os.Args)
}

func getPidMaps(pid int) []string {
	fname := fmt.Sprintf("/proc/%d/maps", pid)

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("could not open '%s' for read: %v", fname, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// use a map to help avoid duplicates
	maps := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 6 && strings.HasPrefix(parts[5], "/") {
			// found something that looks like a file
			maps[parts[5]] = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("reading '%s' failed: %s", fname, err)
	}

	// convert back to a list
	out := make([]string, 0, len(maps))
	for key := range maps {
		out = append(out, key)
	}

	return out
}
