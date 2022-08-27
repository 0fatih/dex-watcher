package utils

import "github.com/fatih/color"

type Colors struct {
	RED    string
	GREEN  string
	YELLOW string
}

var PrintColors = Colors{RED: "red", GREEN: "green", YELLOW: "yellow"}

func ColoredPrint(msg, clr string) {
	var pr *color.Color
	switch clr {
	case "red":
		pr = color.New(color.FgRed)
	case "green":
		pr = color.New(color.FgGreen)
	case "yellow":
		pr = color.New(color.FgYellow)
	default:
		pr = color.New(color.FgWhite)
	}

	pr.Println(msg)
}
