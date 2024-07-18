package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/unidoc/unidoc/pdf/core"
)

func debugInfo(message string) {
	if verbose {
		log.Println(message)
	}
}

func fatalIfError(err error, message string) {
	if err != nil {
		fmt.Printf("ERROR: %s \n", message)
		os.Exit(1)
	}
}

func parseAccessPermissions(permStr string) (core.AccessPermissions, error) {
	var permissions core.AccessPermissions

	// Delete the part before '{' and after '}'.
	start := strings.Index(permStr, "{")
	end := strings.LastIndex(permStr, "}")
	if start == -1 || end == -1 {
		return permissions, fmt.Errorf("incorrect string format")
	}

	permStr = permStr[start+1 : end]

	// Divide string into key/value pairs
	pairs := strings.Split(permStr, ", ")
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return permissions, fmt.Errorf("incorrect key/value pair format: %s", pair)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1]) == "true"

		switch key {
		case "Printing":
			permissions.Printing = value
		case "Modify":
			permissions.Modify = value
		case "ExtractGraphics":
			permissions.ExtractGraphics = value
		case "Annotate":
			permissions.Annotate = value
		case "FillForms":
			permissions.FillForms = value
		case "DisabilityExtract":
			permissions.DisabilityExtract = value
		case "RotateInsert":
			permissions.RotateInsert = value
		case "FullPrintQuality":
			permissions.FullPrintQuality = value
		default:
			return permissions, fmt.Errorf("unknown key: %s", key)
		}
	}

	return permissions, nil
}
