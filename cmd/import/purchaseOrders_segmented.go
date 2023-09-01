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
	"github.com/spf13/cobra"
	orm "snowlastic-cli/pkg/orm"
)

// purchaseOrdersCmd represents the import/purchaseOrders command
var purchaseOrdersSegmentedCmd = &cobra.Command{
	Use:   "purchaseOrders-segmented",
	Short: "Index all Purchase Orders contained in the snowflake database",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			dbSchema  = "CMP"
			dbTable   = "PURCHASEORDERS_FLAGGED" // there are multiple, which are handled in the GetQuery return
			indexName = "purchaseorders"
			docType   = (&orm.PurchaseOrder{}).New()
		)
		return runSegmentedImport(dbSchema, dbTable, indexName, docType)
	},
}

func init() {
	purchaseOrdersSegmentedCmd.Flags().StringVar(&segmenter, "by", "", "a field or SQL aggregating function used to split the import")
}
