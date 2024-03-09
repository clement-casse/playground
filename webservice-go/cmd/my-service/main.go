package main

import (
	"github.com/spf13/cobra"

	"github.com/clement-casse/playground/webservice-go/cmd/my-service/internal"
)

func main() {
	cmd := internal.Command()
	cobra.CheckErr(cmd.Execute())
}
