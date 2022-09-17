/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"context"

	"github.com/adshao/go-binance/v2"
	"github.com/dunstorm/crypto-data/database"
	"github.com/spf13/cobra"
)

var ticker string
var interval string
var start string
var end string
var output string

func parseDate(date string) int64 {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return t.Unix() * 1000
}

func tpToDateTime(tp int64) time.Time {
	return time.Unix(tp/1000, 0)
}

func klinesToRecords(klines []*binance.Kline) [][]string {
	records := make([][]string, len(klines))
	for i, kline := range klines {
		records[i] = []string{
			tpToDateTime(kline.OpenTime).Format("2006-01-02 15:04:05"),
			kline.Open,
			kline.High,
			kline.Low,
			kline.Close,
			kline.Volume,
			tpToDateTime(kline.CloseTime).Format("2006-01-02 15:04:05"),
			kline.QuoteAssetVolume,
			kline.TakerBuyBaseAssetVolume,
			kline.TakerBuyQuoteAssetVolume,
		}
	}
	return records
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download historical crypto data",
	Long: `
	This command downloads historical crypto data from Binance.
	Usage: crypto-data download --ticker=<ticker> --interval=<interval> --start=<start> --end=<end> --output=<output>`,
	Run: func(cmd *cobra.Command, args []string) {
		config := database.FindOrCreateConfig()
		// check for api key and secret if not present ask for it
		if config.BinanceAPIKey == "" || config.BinanceAPISecret == "" {
			fmt.Println("Please enter your Binance API Key:")
			fmt.Scanln(&config.BinanceAPIKey)
			fmt.Println("Please enter your Binance API Secret:")
			fmt.Scanln(&config.BinanceAPISecret)
			database.UpdateConfig(config)
		}

		// if output is not specified use the default
		if output == "" {
			output = "data.csv"
		}

		// parse start and end and convert into timestamp
		startTimestamp := parseDate(start)
		endTimestamp := parseDate(end)

		// download data
		fmt.Println("Downloading data...")

		// create binance client
		client := binance.NewClient(config.BinanceAPIKey, config.BinanceAPISecret)
		// get historical data
		klines, err := client.NewKlinesService().Symbol(ticker).Interval(interval).StartTime(startTimestamp).EndTime(endTimestamp).Do(context.Background())
		if err != nil {
			fmt.Println(err)
		}

		// display different of start and end
		fmt.Printf("Start: %s, End: %s, Data: %d days", start, end, (endTimestamp-startTimestamp)/(1000*60*60*24))

		// save data to csv
		fmt.Println("Saving data to csv...")
		// create csv
		f, err := os.Create(output)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		// write data to csv
		csvWriter := csv.NewWriter(f)
		// write header in title case
		csvWriter.Write([]string{"Date", "Open", "High", "Low", "Close", "Volume", "CloseTime", "QuoteAssetVolume", "TakerBuyBaseAssetVolume", "TakerBuyQuoteAssetVolume"})
		data := klinesToRecords(klines)
		csvWriter.WriteAll(data)

		fmt.Println("Done! Data saved to", output)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// turn off sorting
	downloadCmd.Flags().SortFlags = false

	downloadCmd.Flags().StringVarP(&ticker, "ticker", "t", "", "The ticker of the crypto asset (e.g. BTCUSDT)")
	downloadCmd.Flags().StringVarP(&interval, "interval", "i", "", "The interval of the data (1m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M)")
	downloadCmd.Flags().StringVarP(&start, "start", "s", "", "The start date of the data (format: yyyy-mm-dd)")
	downloadCmd.Flags().StringVarP(&end, "end", "e", "", "The end date of the data (format: yyyy-mm-dd)")
	downloadCmd.Flags().StringVarP(&output, "output", "o", "", "The output file (e.g. data.csv)")

	downloadCmd.MarkFlagRequired("ticker")
	downloadCmd.MarkFlagRequired("interval")
	downloadCmd.MarkFlagRequired("start")
	downloadCmd.MarkFlagRequired("end")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
