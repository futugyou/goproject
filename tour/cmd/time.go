package cmd

import (
	"github/go-project/tour/internal/timer"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:       "time",
	Short:     "time format",
	Long:      "time format",
	ValidArgs: []string{"now", "calc"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// nowTime := timer.GetNowTime()
		// log.Printf("output  %s,   %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}
var nowTimeCmd = &cobra.Command{
	Use:   "now",
	Short: "time now",
	Long:  "time now",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		log.Printf("output  %s,   %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}
var calcTimeCmd = &cobra.Command{
	Use:   "calc",
	Short: "calc time now",
	Long:  "calc time now",
	Run: func(cmd *cobra.Command, args []string) {
		location, _ := time.LoadLocation("Asia/Shanghai")
		var currentTime time.Time
		var layout = "2006-01-02 15:04:05"
		if calculateTime == "" {
			currentTime = timer.GetNowTime()
		} else {
			var err error
			if !strings.Contains(calculateTime, " ") {
				layout = "2016-01-02"
			}
			currentTime, err = time.ParseInLocation(layout, calculateTime, location)
			if err != nil {
				t, _ := strconv.Atoi(calculateTime)
				currentTime = time.Unix(int64(t), 0).In(location)
			}
		}
		calculateTime, err := timer.GetCalculateTime(currentTime, duration)
		if err != nil {
			log.Fatalf("timer getclaculation err : %v", err)
		}

		log.Printf("output : %s ,%d ", calculateTime.Format(layout), calculateTime.Unix())
	},
}

var calculateTime string
var duration string

func init() {
	timeCmd.AddCommand(nowTimeCmd)
	timeCmd.AddCommand(calcTimeCmd)
	calcTimeCmd.Flags().StringVarP(&calculateTime, "calculate", "c", "", "please input calculate !")
	calcTimeCmd.Flags().StringVarP(&duration, "duration", "d", "", "please intout duration !")
}
