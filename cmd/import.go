/*
	package cmd

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
package cmd

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v8"
	"golang.org/x/sync/errgroup"
	"log"
	"path"
	_import "snowlastic-cli/cmd/import"
	"snowlastic-cli/pkg/es"
	types "snowlastic-cli/pkg/orm"
	"snowlastic-cli/pkg/snowflake"
	"sort"
	"text/template"
	"time"
)

var (
	by            string
	givenSegments []string
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   `import index_name --id ID [--from ./path/to/query-index_name.sql] [--by Database [--in Value1 --in "Value 2"  --in ...]]`,
	Short: "Import documents into an elasticsearch index",
	Long: `A document is a representation of any kind of record. This tool allows
for importing data from pre-defined sources (such as snowflake tables/views) or 
from a json file containing a list of documents.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		setLogLevel()
		if cmd.Flags().Lookup("in").Changed && !cmd.Flags().Lookup("by").Changed {
			return errors.New("you must specify a 'by' field in order to use the 'in' argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// arguments?
		//var batchSize = 10
		var indexName = args[0]

		// initialization
		var err error
		var pCtx, cancel = context.WithCancelCause(context.TODO())
		var g, ctx = errgroup.WithContext(pCtx)
		g.SetLimit(5)
		var progress = mpb.NewWithContext(ctx)

		var c = make(chan types.SnowlasticDocument, es.BulkInsertSize)
		var batches chan []types.SnowlasticDocument
		var totalSize int64
		var client *elasticsearch.Client
		var start = time.Now()

		log.Println("connecting to elasticsearch...")
		client, err = generateDefaultElasticClient()
		if err != nil {
			cancel(err)
		}

		switch viper.GetString("file") {
		default:
			{
				var docs []types.SnowlasticDocument
				log.Println("reading records from file", viper.GetString("file"))
				docs, err = _import.GetRecords(viper.GetString("file"))
				if err != nil {
					cancel(err)
				}
				totalSize = int64(len(docs))
				log.Printf("sending %d documents to pipeline...\n", len(docs))
				go func() {
					g.Go(func() error {
						for i := 0; i < len(docs); i++ {
							select {
							case <-ctx.Done():
							default:
								{

									var doc = docs[i]
									c <- doc
								}
							}
						}
						log.Println("completed sending documents...")
						return nil
					})
					err = g.Wait()
					if err != nil {
						cancel(err)
					}
					close(c)
				}()

			}
		case "":
			{
				var (
					segmenter        *string
					schemas          = viper.GetStringSlice("snowflakeSchemas")
					segments         = make(map[string]map[any]int64)
					segmentedQueries = make(map[string]string)
					tmpl             *template.Template

					db *sql.DB
				)
				switch {
				// as additional, non-sap sources get added,
				case indexName == "navex_cases":
					schemas = []string{"SQL_NAVEX"}
				default:
					sort.Strings(schemas)
				}
				if by != "" {
					segmenter = &by
				} else {
					segmenter = nil
				}
				tmpl, err := template.ParseFiles(path.Join(
					viper.GetString("settingsDirectory"),
					fmt.Sprintf("query-%s.sql", indexName),
				))

				log.Println("connecting to database...")
				db, err = snowflake.NewDB(snowflake.Config{
					Account:   viper.GetString("snowflakeAccount"),
					Warehouse: viper.GetString("snowflakeWarehouse"),
					Database:  viper.GetString("snowflakeDatabase"),
					User:      viper.GetString("snowflakeUser"),
					Password:  viper.GetString("snowflakePassword"),
					Role:      viper.GetString("snowflakeRole"),
				})
				if err != nil {
					cancel(fmt.Errorf("error encountered while connecting to the database %s: %s", indexName, err))
				}

				log.Printf("determining segments for each schema (%d)...\n", len(schemas))
				var bar = _import.AddBarMomentary(
					progress,
					int64(len(schemas)),
					viper.GetString("snowflakeDatabase"),
					by,
				)

				start := time.Now()

				var h, segCtx = errgroup.WithContext(ctx)
				h.SetLimit(5)
				// segmentation:
				// 1) construct a query to get distinct segments
				// 2) construct a paramterized query to get the actual records in each segment
				// 3) query the database with the query constructed in 1) to get segments and their counts
				// 4) return the segments, their counts, and a query constructed to get the records
				for i := range schemas {
					var (
						schema         = schemas[i]
						segmentCounts  map[any]int64
						segmentedQuery string
					)
					h.Go(func() error {
						segmentedQuery, err = getSegmentedQuery(segCtx, tmpl, viper.GetString("snowflakeDatabase"), schema, segmenter)
						segmentCounts, err = getSegmentCounts(segCtx, db, tmpl, viper.GetString("snowflakeDatabase"), schema, segmenter)
						if err != nil {
							return err
						}
						segments[schema] = segmentCounts
						segmentedQueries[schema] = segmentedQuery
						for _, count := range segmentCounts {
							totalSize += count
						}
						bar.EwmaIncrement(time.Since(start))
						return nil
					})
				}
				err = h.Wait()
				if err != nil {
					cancel(fmt.Errorf("error encountered during waiting for h %s: %s", indexName, err))
				}

				go func() {
					for i := 0; i < len(schemas); i++ {
						var (
							schema      = schemas[i]
							segmentKeys = SortMapKeys(segments[schema])
							query       = segmentedQueries[schema]
						)
						for j := 0; j < len(segmentKeys); j++ {
							var (
								segment = segmentKeys[j]
								count   = segments[schema][segment]
							)
							g.Go(func() error {
								var bar = _import.AddBarMomentary(progress, count, schema, segment)
								var rows *sql.Rows
								var start = time.Now()

								rows, err = db.QueryContext(ctx, query, segment)
								if err != nil {
									cancel(fmt.Errorf("error encountered getting rows of %s, %s at %s: %s", indexName, schema, segment, err))
								}
								for rows.Next() {
									var doc = types.NewDocument()
									err = doc.ScanFrom(rows)
									if err != nil {
										log.Println(err)
										cancel(fmt.Errorf("error encountered iterating rows of %s, %s at %s: %s", indexName, schema, segment, err))
									}
									c <- doc
									bar.EwmaIncrement(time.Since(start))
								}
								return nil
							})
						}
					}
					err = g.Wait()
					if err != nil {
						cancel(err)
					}
					close(c)
				}()
			}
		}

		select {
		case <-ctx.Done():
		default:
			{
				bar := _import.AddBarPersistent(progress, totalSize, indexName, "all")
				//log.Println("batching documents for import...")
				batches = es.BatchEntities(c, es.BulkInsertSize)
				numIndexed, numErrors, err := es.BulkImport(client, batches, indexName, bar)
				if err != nil {
					cancel(fmt.Errorf("problem with bulk importing for segment %s: %s", indexName, err))
				}
				if numErrors > 0 {
					cancel(errors.New(fmt.Sprintf(
						"%s:\tIndexed [%s] documents with [%s] errors in %s (%s docs/sec)",
						fmt.Sprintf("%s: %s=%v", indexName, viper.GetString("by"), "all"),
						humanize.Comma(numIndexed),
						humanize.Comma(numErrors),
						time.Since(start).Truncate(time.Millisecond),
						humanize.Comma(int64(1000.0/float64(time.Since(start)/time.Millisecond)*float64(numIndexed))),
					)))
				}
				progress.Wait()
			}
		}
		return context.Cause(pCtx)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().String("id", "id", "field which uniquely identifies each record (usually 'ID')")
	//_ = importCmd.MarkFlagRequired("id")
	_ = viper.BindPFlag("identifier", importCmd.Flags().Lookup("id"))

	importCmd.Flags().String("from", "", "location of query file (defaults to settingsDirectory/query-[index_name].sql)")
	importCmd.Flags().String("file", "", "location of a json file with records to import into an elasticsearch index")
	importCmd.MarkFlagsMutuallyExclusive("from", "file")
	_ = viper.BindPFlag("from", importCmd.Flags().Lookup("from"))
	_ = viper.BindPFlag("file", importCmd.Flags().Lookup("file"))

	importCmd.Flags().StringVar(&by, "by", "", "a field or SQL aggregating function used to split the import")
	importCmd.Flags().StringSliceVar(&givenSegments, "in", nil, "limit the import of the field defined in the 'by' argument to a comma seperated list of values")
}

func getSegmentCounts(ctx context.Context, db *sql.DB, baseQuery *template.Template, database, schema string, segmenter *string) (segmentCounts map[any]int64, err error) {
	var (
		_segmentCounts    = make(map[any]int64)
		segmentationQuery = bytes.Buffer{}
		templateData      = map[string]string{"database": database, "schema": schema}

		rows *sql.Rows
	)

	select {
	case <-ctx.Done():
		return segmentCounts, err
	default:
		// generate both the counting query for each segment
		// and prepare a parameterized query
		{
			switch {
			case segmenter == nil:
				// all rows without segmentation (nil segmenter)
				{
					{
						// segmentation query building: used to get counts
						segmentationQuery.WriteString(`SELECT NULL AS SEG, COUNT(*) AS CNT FROM (`)
						segmentationQuery.WriteString("\n")
						err = baseQuery.Execute(&segmentationQuery, templateData)
						if err != nil {
							return segmentCounts, err
						}
						segmentationQuery.WriteString("\n")
						segmentationQuery.WriteString(")")
					}
				}
			default: // partitioned import by segmenter field
				{
					// segmentation query building: used to get counts
					var quotedIdentifier = snowflake.QuoteIdentifier(*segmenter)
					segmentationQuery.WriteString("SELECT DISTINCT ")
					segmentationQuery.WriteString(quotedIdentifier)
					segmentationQuery.WriteString(", COUNT(*) AS CNT FROM (")
					segmentationQuery.WriteString("\n")
					err = baseQuery.Execute(&segmentationQuery, templateData)
					if err != nil {
						return segmentCounts, err
					}
					segmentationQuery.WriteString("\n")
					segmentationQuery.WriteString(")")
					segmentationQuery.WriteString("\n")
					segmentationQuery.WriteString("GROUP BY ")
					segmentationQuery.WriteString(quotedIdentifier)
					segmentationQuery.WriteString("\n")
					segmentationQuery.WriteString("ORDER BY ")
					segmentationQuery.WriteString(quotedIdentifier)
				}
			}
		}
	}

	rows, err = db.QueryContext(ctx, segmentationQuery.String())
	if err != nil {
		return segmentCounts, err
	}
	for rows.Next() {
		var segment struct {
			Value any
			Count int64
		}
		err = rows.Scan(&segment.Value, &segment.Count)
		if err != nil {
			return segmentCounts, err
		}
		if segment.Count != 0 { // no need to add schema segments which have no rows
			_segmentCounts[segment.Value] = segment.Count
		}
	}

	return _segmentCounts, nil
}
func getSegmentedQuery(ctx context.Context, baseQuery *template.Template, database, schema string, segmenter *string) (segmentedQuery string, err error) {
	var (
		_segmentedQuery = bytes.Buffer{}
		templateData    = map[string]string{"database": database, "schema": schema}
	)

	select {
	case <-ctx.Done():
		return _segmentedQuery.String(), err
	default:
		{
			switch {
			case segmenter == nil:
				// all rows without segmentation (nil segmenter)
				{
					// to get actual rows, just use the query
					err = baseQuery.Execute(&_segmentedQuery, templateData)
					if err != nil {
						return _segmentedQuery.String(), err
					}
				}
			default: // partitioned import by segmenter field
				{
					// prepare the base query to accept a segment parameter
					var quotedIdentifier = snowflake.QuoteIdentifier(*segmenter)

					_segmentedQuery.WriteString("SELECT * FROM (")
					_segmentedQuery.WriteString("\n")
					err = baseQuery.Execute(&_segmentedQuery, templateData)
					if err != nil {
						return _segmentedQuery.String(), err
					}
					_segmentedQuery.WriteString(")")
					_segmentedQuery.WriteString("\n")
					_segmentedQuery.WriteString("WHERE ")
					_segmentedQuery.WriteString(quotedIdentifier)
					_segmentedQuery.WriteString(" = ?")
				}
			}
		}
	}
	return _segmentedQuery.String(), err
}

func SortMapKeys(m map[any]int64) []any {
	var ret = make([]any, len(m))
	var i int
	for k := range m {
		ret[i] = k
		i++
	}
	sort.Slice(ret, func(i, j int) bool {
		var a = fmt.Sprintf("%v", ret[i])
		var b = fmt.Sprintf("%v", ret[j])
		if a < b {
			return true
		}
		return false
	})
	return ret
}
