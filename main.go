// mubi2letterboxd is a simple utility for user data migration from MUBI to letterboxd.
// With the utility, you can create a .csv file suitable for manual import to Letterboxd.
//
// inspired by the reddit entry by jcunews1
// https://www.reddit.com/r/learnjavascript/comments/auwynr/export_mubi_data/ehcx2zf/

package main

import (
	"flag"
	"mubi2letterboxd/cli"
	"mubi2letterboxd/gui"
)

func main() {
	noGUI := flag.Bool("disable-gui", false, "Disable GUI")
	flag.Parse()

	if *noGUI {
		cli.ProcessCli()
	} else {
		gui.ProcessGui()
	}
}
