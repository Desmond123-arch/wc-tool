package cmd

import (
	"fmt"
	"os"

	// homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cFlag bool
var lFlag bool
var wFlag bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "default is $Home/ccwc.yml")
	rootCmd.PersistentFlags().StringP("author", "a", "Desmond Opoku Mends", "Author name for copyright attribution")
	rootCmd.PersistentFlags().BoolVarP(&cFlag, "count", "c", false, "Prints the number of bytes in a file")
	rootCmd.PersistentFlags().BoolVarP(&lFlag, "lines", "l", false, "Prints the number of lines in a file")
	rootCmd.PersistentFlags().BoolVarP(&wFlag, "words", "w", false, "Prints the number of words in a file")
}

func initConfig(){
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// home, err := homedir.Dir()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		viper.AddConfigPath(".")
		viper.SetConfigFile("ccwc.yml")
		viper.SetConfigType("yml")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config: ", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: "ccwc",
	Short: "A word count (wc tool)",
	Long: `A mordern day word count tool built with
			cobra in go`,
	Run: func (cmd *cobra.Command, args []string)  {
		if len(args) == 0 {
			fmt.Println("A file path is expected")
			return
		}
		file_name := args[0]
		file, err := os.ReadFile(file_name)
		if err != nil {
			fmt.Println("Error while reading file")
			return 
		}
		count := countBytes(file)
		if cFlag {
			fmt.Println(count, file_name)
		}
		lines := lineCount(file)
		if lFlag {
			fmt.Println(lines, file_name)
		}
		words := wordCount(file)
		if wFlag {
			fmt.Println(words, file_name)
		}
	},
} 

func Execute(){
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Count the number of bytes
func countBytes(file []byte) int {
	return len(file)
}

//Count the number of lines
func lineCount(file []byte) int {
	lines := 0
	for _, b := range file{
		if string(b) == "\n" || string(b) == "/r"{
			lines++
		}
	}
	return lines
}

func wordCount(file []byte) int {
	if len(file) == 0{
		return 0
	}
	words := 0
	var inWord bool;
	for _, b := range file {
		letter := string(b)
		if letter == " " || letter == "\t" || letter == "\n" || letter == "\r" {
			if inWord{
				//we are outside a word now
				words ++
				inWord = false
			}
		} else {
			//in a word
			inWord = true
		}
	}
	if (inWord){
		words++
	}
	return words
}