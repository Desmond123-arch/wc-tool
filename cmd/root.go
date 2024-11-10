package cmd

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	// homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cFlag bool
var lFlag bool
var wFlag bool
var mFlag bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "default is $Home/ccwc.yml")
	rootCmd.PersistentFlags().StringP("author", "a", "Desmond Opoku Mends", "Author name for copyright attribution")
	rootCmd.PersistentFlags().BoolVarP(&cFlag, "bytes", "c", false, "Prints the number of bytes in a file")
	rootCmd.PersistentFlags().BoolVarP(&lFlag, "lines", "l", false, "Prints the number of lines in a file")
	rootCmd.PersistentFlags().BoolVarP(&wFlag, "words", "w", false, "Prints the number of words in a file")
	rootCmd.PersistentFlags().BoolVarP(&mFlag, "chars", "m", false, "Prints the number of characters in a file(locale)")
}

func initConfig() {
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
	Use:   "ccwc",
	Short: "A word count (wc tool)",
	Long: `A mordern day word count tool built with
			cobra in go`,
	Run: func(cmd *cobra.Command, args []string) {
		var file []byte
		file_name := ""
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			val, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Printf("Error reading piped value")
			}
			file = val
		} else {
			if args[0] == "" {
				fmt.Print("A file path is expected")
				return
			}
			file_name = args[0]
			file, _ = os.ReadFile(file_name)
		}
			
		if cFlag {
			bytes := countBytes(file)
			fmt.Printf("\t%d %s", bytes, file_name)
		} else if lFlag {
			lines := lineCount(file)
			fmt.Printf("\t%d %s", lines, file_name)
		} else if wFlag {
			words := wordCount(file)
			fmt.Printf("\t%d %s", words, file_name)
		} else if mFlag {
			chars := charCount(file)
			fmt.Printf("\t%d %s", chars, file_name)
		} else {
			bytes := countBytes(file)
			lines := lineCount(file)
			words := wordCount(file)
			fmt.Printf("\t%d \t%d \t%d %s", lines, words, bytes, file_name)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Count the number of bytes
func countBytes(file []byte) int {
	return len(file)
}

// Count the number of lines
func lineCount(file []byte) int {
	lines := 0
	for _, b := range file {
		if string(b) == "\n" || string(b) == "/r" {
			lines++
		}
	}
	return lines
}

func wordCount(file []byte) int {
	if len(file) == 0 {
		return 0
	}
	words := 0
	var inWord bool
	for _, b := range file {
		letter := string(b)
		if letter == " " || letter == "\t" || letter == "\n" || letter == "\r" {
			if inWord {
				//we are outside a word now
				words++
				inWord = false
			}
		} else {
			//in a word
			inWord = true
		}
	}
	if inWord {
		words++
	}
	return words
}

func charCount(file []byte) int {
	return utf8.RuneCount(file)
}
