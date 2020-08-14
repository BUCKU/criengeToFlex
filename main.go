package main

import (
	"bufio"
	"flag"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	input := flag.String("i", "./input/", "input path")
	output := flag.String("o", "./output/", "output path")
	relative := flag.Bool("r", false, "is relative")

	flag.Parse()

	workdir, _ := os.Getwd()
	in := *input
	out := *output
	if *relative {
		in = workdir + "/" +in
		out = workdir + "/" +out
	}



	log.Printf("\nInput: %s\nOutput: %s\n", in, out)

	err := copyFiles(in, out)
	if err != nil {
		log.Printf("ERR: At copying files %s", err.Error())
	}
}

func DecodeWindows1251(enc []byte) string {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(enc)
	return string(out)
}

func EncodeWindows1251(inp string) string {
	enc := charmap.Windows1251.NewEncoder()
	out, _ := enc.String(inp)
	return out
}

func copyFiles(in, out string) (err error){
	os.MkdirAll(out, 0664)
	dirContents, err := ioutil.ReadDir(in)
	if err != nil {
		log.Printf("Cant read dir: %s, error: %s", in, err.Error())
	}


	for _, fi := range dirContents {
		if !fi.IsDir() {
			log.Printf("Copying file: %s", in + "/" + fi.Name())
			file, err := os.Open(in+"/"+fi.Name())
			if err != nil {
				log.Printf("ERR: cant open file %s, error: %s", in +"/" +fi.Name(), err.Error())
			}
			defer file.Close()

			fileOut, err := os.OpenFile(out + "/" + fi.Name(), os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("ERR: cant open file %s, error: %s", out +"/" +fi.Name(), err.Error())
			}
			defer fileOut.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan(){
				decoded := DecodeWindows1251([]byte(scanner.Text()))
				fileOut.WriteString(decoded + "\n")
			}
		} else {
			copyFiles(in+"/"+fi.Name(), out+"/"+fi.Name())
		}
	}
	return nil
}
