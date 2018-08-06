package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:      "couple",
			Usage:     "couple machine",
			ArgsUsage: "<personA> <personB>",
			Flags: []cli.Flag{
				cli.StringSliceFlag{Name: "kind", Value: &cli.StringSlice{"amour", "amitie", "travail"}},
			},
			Action: func(c *cli.Context) error {
				a := c.Args()[0]
				b := c.Args()[1]
				for _, kind := range c.StringSlice("kind") {
					fmt.Printf("%s: %d\n", kind, CompatNamesByKind(a, b, kind))
				}
				fmt.Printf("\ntotal: %d\n", CompatNames(a, b))
				return nil
			},
		}, {
			Name:  "group",
			Usage: "group machine",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "kind"},
			},
			ArgsUsage: "<personA> [personB [personC...]]",
			Action: func(c *cli.Context) error {
				names := c.Args()
				results := map[int][][]string{}
				for idx, a := range names {
					bRange := names
					if c.String("kind") != "" {
						bRange = names[idx:]
					}
					for _, b := range bRange {
						var score int
						if kind := c.String("kind"); kind != "" {
							score = CompatNamesByKind(a, b, kind)
						} else {
							score = CompatNames(a, b)
						}
						if _, found := results[score]; !found {
							results[score] = [][]string{}
						}
						results[score] = append(results[score], []string{a, b})
					}
				}
				uniqScores := []int{}
				for score := range results {
					uniqScores = append(uniqScores, score)
				}
				sort.Ints(uniqScores)
				for idx := len(uniqScores) - 1; idx >= 0; idx-- {
					score := uniqScores[idx]
					couples := results[score]
					fmt.Printf("# %d\n", score)
					for _, couple := range couples {
						fmt.Printf("  - %s %s\n", couple[0], couple[1])
					}
					fmt.Println()
				}
				return nil
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func CompatNames(a, b string) int {
	return CompatNamesByKind(a, "", b)
}

func CompatNamesByKind(a, b string, kind string) int {
	score := []int{}
	for _, letter := range kind {
		letterScore := 0
		for _, al := range a {
			if al == letter {
				letterScore++
			}
		}
		for _, bl := range b {
			if bl == letter {
				letterScore++
			}
		}
		score = append(score, letterScore)
	}
	for len(score) > 2 {
		if len(score) == 3 && score[0] == 1 && score[1] == 0 && score[2] == 0 {
			return 100
		}
		tmpScore := []int{}
		for i := 0; i < len(score)-1; i++ {
			sum := score[i] + score[i+1]
			if sum > 9 {
				sum = int(sum/10) + int(sum%10)
			}
			tmpScore = append(tmpScore, sum)
		}
		score = tmpScore
	}
	return score[0]*10 + score[1]
}
