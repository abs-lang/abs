package install

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func valid(module string) bool {
	return strings.HasPrefix(module, "github.com/")
}

func Install(module string) {
	if !valid(module) {
		fmt.Printf(`Error reading URL. Please use "github.com/USER/REPO" format to install`)
		return
	}

	err := getZip(module)
	if err != nil {
		return
	}

	err = unzip(module)
	if err != nil {
		fmt.Printf(`Error unpacking: %v`, err)
		return
	}

	alias, err := createAlias(module)
	if err != nil {
		return
	}

	fmt.Printf("\nInstall Success. You can use the module with `require(\"%s\")`\n", alias)
	return
}

func printLoader(done chan int64, message string) {
	var stop bool = false
	symbols := []string{"🌑 ", "🌒 ", "🌓 ", "🌔 ", "🌕 ", "🌖 ", "🌗 ", "🌘 "}
	i := 0

	for {
		select {
		case <-done:
			stop = true
		default:
			fmt.Printf("\r" + symbols[i] + " - " + message)
			time.Sleep(100 * time.Millisecond)
			i++
			if i > len(symbols)-1 {
				i = 0
			}
		}

		if stop {
			break
		}
	}
}

func getZip(module string) error {
	path := fmt.Sprintf("./vendor/%s-master.zip", module)
	// Create all the parent directories if needed
	err := os.MkdirAll(filepath.Dir(path), 0755)

	if err != nil {
		fmt.Printf("Error making directory %s\n", err)
		return err
	}

	out, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		fmt.Printf("Error opening file %s\n", err)
		return err
	}
	defer out.Close()

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	url := fmt.Sprintf("https://%s/archive/master.zip", module)

	done := make(chan int64)
	go printLoader(done, "Downloading archive")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating new request %s", err)
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Could not get module: %s\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Errorf("Bad response code: %d", resp.StatusCode)
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error copying file %s", err)
		return err
	}
	done <- 1
	return err
}

// Unzip will decompress a zip archive
func unzip(module string) error {
	fmt.Printf("\nUnpacking...")
	src := fmt.Sprintf("./vendor/%s-master.zip", module)
	dest := filepath.Dir(src)

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		filename := f.Name
		parts := strings.Split(f.Name, string(os.PathSeparator))
		if len(parts) > 1 {
			if strings.HasSuffix(parts[0], "-master") {
				// Trim "master" suffix due to github's naming convention for archives
				parts[0] = strings.TrimSuffix(parts[0], "-master")
				filename = strings.Join(parts, string(os.PathSeparator))
			}
		}
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, filename)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, 0755)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func createAlias(module string) (string, error) {
	fmt.Printf("\nCreating alias...")
	f, err := os.OpenFile("./packages.abs.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Could not open alias file %s\n", err)
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)

	data := make(map[string]string)
	moduleName := filepath.Base(module)
	modulePath := fmt.Sprintf("./vendor/%s", module)

	// If package.abs.json file is empty
	if len(b) == 0 {
		// Add alias key-value pair to file
		data[moduleName] = modulePath
	} else {
		err = json.Unmarshal(b, &data)
		if err != nil {
			fmt.Printf("Could not unmarshal alias json %s\n", err)
			return "", err
		}
		// module already installed and aliased
		if data[moduleName] == modulePath {
			return moduleName, nil
		}

		if data[moduleName] != "" {
			fmt.Printf("This module could not be aliased because module of same name exists\n")
			return modulePath, nil
		}

		data[moduleName] = modulePath
	}

	newData, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		fmt.Printf("Could not marshal alias json when installing module %s\n", err)
		return "", err
	}

	_, err = f.WriteAt(newData, 0)
	if err != nil {
		fmt.Printf("Could not write to alias file %s\n", err)
		return "", err
	}
	return moduleName, err

}
