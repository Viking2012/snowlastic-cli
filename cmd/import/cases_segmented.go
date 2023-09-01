/*
Copyright © 2023 Alexander Orban <alexander.orban@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package _import

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	orm "snowlastic-cli/pkg/orm"
)

// import/casesCmd represents the import/cases command
var casesSegmentedCmd = &cobra.Command{
	Use:   "cases-segmented",
	Short: "Index all Navex Cases contained in the snowflake database",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("esClient") {
			return errors.New("elasticsearch was somehow not created by the `import` command prior to running `demos`")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			dbSchema  = "SQL_NAVEX"
			dbTable   = "" // there are multiple, which are handled in the GetQuery return
			indexName = "cases"
			docType   = (&orm.Case{}).New()
		)
		return runSegmentedImport(dbSchema, dbTable, indexName, docType)
	},
}

func init() {}
