package utils

import "github.com/fatih/color"

func ColoredPrint(msg, clr string) {
	var pr *color.Color
	switch clr {
	case "red":
		pr = color.New(color.FgRed)
	case "yellow":
		pr = color.New(color.FgYellow)
	case "green":
		pr = color.New(color.FgGreen)
	default:
		pr = color.New(color.FgWhite)
	}

	pr.Println(msg)
}
