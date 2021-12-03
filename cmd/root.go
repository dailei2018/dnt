package cmd

import (
	"fmt"
	"os"

	"github.com/dailei2018/dnt/lib"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	NWorker int
	Allow   []map[string]interface{}
	Output  []map[string]interface{}

	BufferPoolSize int
	BufferSize     int

	TLS struct {
		Tag   []uint8
		Topic []string
	}

	FILE struct {
		Tag []uint8
	}

	Kafka struct {
		Broker string
	}
}

var (
	cfgFile    string
	MainConfig lib.Config
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("aaaaaaa", args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "dnt.yaml", "")
}

func initConfig() {
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		logrus.Fatalf("load %s failed, %v", viper.ConfigFileUsed(), err)
	}

	if err := viper.Unmarshal(&MainConfig); err != nil {
		logrus.Fatalf("unable to decode into struct, %v", err)
	}
}
