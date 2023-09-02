package gov

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"gl-source/database"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ProcessDownload() {

	// Define the URL for the CSV file download.
	csvURL := "http://download.companieshouse.gov.uk/BasicCompanyDataAsOneFile-2023-09-01.zip"

	// Define the local file path where the CSV will be saved permanently.
	localFilePath := "BasicCompanyDataAsOneFile-2023-09-01.csv"

	var file *os.File
	// Check if the file already exists before downloading.
	if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
		// File does not exist, proceed with downloading.
		log.Println("File does not exist, downloading...")
		file, err = downloadCSVFile(csvURL, localFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	} else {
		fmt.Println("File already exists. Skipping download.")
		// File already exists, open it for use.
		var err error
		file, err = os.Open(localFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	// Open the ZIP file for reading.
	log.Println("Opening zipfile..")
	zipFile, err := zip.OpenReader(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing zip..")
		if err := zipFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Iterate through the files in the ZIP archive and look for the CSV file.
	log.Println("Extracting csv..")
	var csvFile *zip.File
	for _, file := range zipFile.File {
		if filepath.Ext(file.Name) == ".csv" {
			csvFile = file
			break
		}
	}
	if csvFile == nil {
		log.Fatal("CSV file not found in the ZIP archive")
	}

	// Open the CSV file within the ZIP archive for reading.
	log.Println("Opening CSV...")
	reader, err := csvFile.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing CSV...")
		if err := reader.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Create a CSV reader for the file.
	csvReader := csv.NewReader(reader)

	// Read and discard the header row (assuming the headers are not needed).
	log.Println("Discarding header..")
	_, err = csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Rolling back due to error:", r)
			tx.Rollback()
		} else {
			log.Println("Committing... ")
			tx.Commit()
		}
	}()

	recordCh := make(chan []string, 1000) // Create a channel to send records to Goroutines.
	doneCh := make(chan bool)             // Create a channel to signal when all Goroutines are done.

	// Start Goroutines to process records concurrently.
	for i := 0; i < 4; i++ { // Adjust the number of Goroutines as needed.
		go func() {
			for record := range recordCh {
				var company database.Company
				if err := company.MapCSVData(record); err != nil {
					log.Println("Error mapping CSV data:", err)
					continue
				}

				result := tx.Create(&company)
				if result.Error != nil {
					log.Println("Error saving record to the database:", result.Error)
					continue
				}
			}
			doneCh <- true // Signal that this Goroutine is done.
		}()
	}

	i := 0
	for {
		i++
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV record %d: %s", i, err)
			continue
		}

		recordCh <- record // Send the record to a Goroutine for processing.

		if i%10000 == 0 {
			log.Printf("Processed: %d records", i)
			//break
		}
	}

	close(recordCh) // Close the record channel to signal that no more records will be sent.

	// Wait for all Goroutines to finish.
	for i := 0; i < 4; i++ { // Adjust the number based on the number of Goroutines started.
		<-doneCh
	}
}

func downloadCSVFile(url, destination string) (*os.File, error) {
	// Create the destination file.
	outputFile, err := os.Create(destination)
	if err != nil {
		return nil, err
	}

	// Send an HTTP GET request to download the file.
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	// Check if the response status code is valid.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	// Copy the contents of the response body to the destination file.
	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Downloaded CSV file: %s\n", destination)

	return outputFile, nil
}
