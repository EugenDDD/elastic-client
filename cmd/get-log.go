package cmd

import (
	"context"
	"elastic-search/global"
	"elastic-search/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var getLogCmd = &cobra.Command{
	Use:   "get-log",
	Short: "Get elastic log",
	Long:  `Get elastic config details`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		configSerial, err = ioutil.ReadFile(global.ConfigFile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(configSerial, &config)
		if err != nil {
			panic(err)
		}

		// check connection to global.Server
		if ok, err := isReacheable(config.Server); err == nil && !ok {
			panic(fmt.Sprintf("No ping response to kibana server: %s", config.Server))
		}

		key := "elastic the the best way to use logs"
		config.Pass, err = util.DecryptString(config.Pass, key)
		if err != nil {
			panic("Error when dec string")
		}

		client, err = elastic.NewClient(
			elastic.SetURL(fmt.Sprintf("%s:%s", config.Server, fmt.Sprintf("%v", config.Port))),
			elastic.SetHealthcheck(false),
			elastic.SetSniff(false),
			elastic.SetBasicAuth(config.User, config.Pass),
		)

		if err != nil {
			fmt.Println("Cannot create client for server:", config.Server)
			return
		}
	},

	Run: func(cmd *cobra.Command, args []string) {

		index := getFlags.Index

		noFlag := true
		query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery())

		if cmd.PersistentFlags().Lookup("route-id").Changed {
			updateQuery(query, noFlag, "RouteId", getFlags.RouteID)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("country").Changed {
			updateQuery(query, noFlag, "Country", getFlags.Country)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("store").Changed {
			updateQuery(query, noFlag, "Store", getFlags.Store)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("uuid").Changed {
			updateQuery(query, noFlag, "UUID", getFlags.UUID)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("service").Changed {
			updateQuery(query, noFlag, "Service", getFlags.Service)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("key1").Changed {
			updateQuery(query, noFlag, "Key1", getFlags.Key1)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("key2").Changed {
			updateQuery(query, noFlag, "Key2", getFlags.Key2)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("key3").Changed {
			updateQuery(query, noFlag, "Key3", getFlags.Key3)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("key4").Changed {
			updateQuery(query, noFlag, "Key4", getFlags.Key4)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("status").Changed {
			updateQuery(query, noFlag, "Status", getFlags.Status)
			noFlag = false
		}

		if cmd.PersistentFlags().Lookup("results").Changed {
			if getFlags.Results > 10000 {
				fmt.Println("Please enter a value lower than 10000 for \"results\" flag!")
				return
			}
			noFlag = false
		}

		now := time.Now()
		backTo := now.AddDate(0, 0, -getFlags.DaysBack)

		query.Must(elastic.NewRangeQuery("@timestamp").Gte(backTo).Lt(now))

		if noFlag {
			fmt.Printf("Please provide at least one flag!\n\n")
			cmd.Usage()
			return
		}

		// //  print flags
		// variable, err := query.Source()
		// if err != nil {
		// 	panic(err)
		// }

		// data, err := json.Marshal(variable)
		// if err != nil {
		// 	panic(err)
		// }
		// got := string(data)

		// fmt.Println(got)
		searchResults, err := client.Search().Index(index).Query(query).From(0).Size(getFlags.Results).Do(context.Background())
		if err != nil {
			panic(err)
		}

		fmt.Printf("Query took %d milliseconds\n", searchResults.TookInMillis)
		if searchResults.Hits.TotalHits > 0 {
			fmt.Printf("Found %d results\n\n", searchResults.Hits.TotalHits)
			for _, hit := range searchResults.Hits.Hits {
				var result global.Result
				json.Unmarshal(*hit.Source, &result)

				if !cmd.PersistentFlags().Lookup("to-file").Changed {
					fmt.Println(result.String())
				}

				results = append(results, result)
			}
		} else {
			fmt.Println("No results found!")
		}

		if cmd.PersistentFlags().Lookup("to-file").Changed {
			util.WriteResults(results, "./results.log")
		}
	},
}

func init() {

	getLogCmd.PersistentFlags().StringVarP(&getFlags.Salesline, "salesline", "", "mcc", "Set salesline for index search")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Index, "index", "", "*mw_mcc_prod*", "Set index")

	getLogCmd.PersistentFlags().StringVarP(&getFlags.RouteID, "route-id", "", "", "Set RouteID")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Country, "country", "", "", "Set country")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Store, "store", "", "", "Set store")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.UUID, "uuid", "", "", "Set UUID")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Service, "service", "", "", "Set service")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Key1, "key1", "", "", "Set Key1")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Key2, "key2", "", "", "Set Key2")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Key3, "key3", "", "", "Set Key3")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Key4, "key4", "", "", "Set Key4")
	getLogCmd.PersistentFlags().StringVarP(&getFlags.Status, "status", "", "", "Set status")

	getLogCmd.PersistentFlags().IntVarP(&getFlags.DaysBack, "days-back", "", 1, "Set days back")
	getLogCmd.PersistentFlags().IntVarP(&getFlags.Results, "results", "", 100, "Set results - max 10000 records")
	getLogCmd.PersistentFlags().BoolVarP(&getFlags.ToFile, "to-file", "F", false, "Send results to \"results.log\" file")
}

func getPhraseQueryList(field, params string) []*elastic.MatchPhraseQuery {

	values := strings.Split(params, ",")
	query := []*elastic.MatchPhraseQuery{}

	for _, value := range values {
		aQuery := elastic.NewMatchPhraseQuery(field, value)
		query = append(query, aQuery)
	}

	return query
}

func updateQuery(query *elastic.BoolQuery, isFirstFlag bool, name, values string) {

	aQuery := elastic.NewBoolQuery()
	for _, value := range getPhraseQueryList(name, values) {
		aQuery.Should(value)
	}
	aQuery.MinimumNumberShouldMatch(1)

	if isFirstFlag {
		*query = *aQuery
	} else {
		query.Must(aQuery)
	}
}

func isReacheable(ipAddress string) (bool, error) {

	var exp = regexp.MustCompile(`^http://`)
	ipAddress = exp.ReplaceAllString(ipAddress, "")

	if runtime.GOOS == "windows" {
		fmt.Println(ipAddress)
		if err := exec.Command("ping", ipAddress, "-n 2").Run(); err != nil {
			return false, err
		}
	} else {
		if err := exec.Command("ping", ipAddress, "-c 1", "-i 1").Run(); err != nil {
			return false, err
		}
	}

	return true, nil
}
