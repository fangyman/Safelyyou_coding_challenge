package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"safelyyou_coding_challenge_go/database"
	"safelyyou_coding_challenge_go/models"
)

// CSVService handles CSV file operations
type CSVService struct {
	db *database.DatabaseOperations
}

// NewCSVService creates a new CSV service
func NewCSVService(db *database.DatabaseOperations) *CSVService {
	return &CSVService{db: db}
}

// LoadDevicesFromCSV loads devices from CSV and adds new ones to database
func (s *CSVService) LoadDevicesFromCSV(csvPath string) error {
	csvDeviceIDs, err := s.readDeviceIDsFromCSV(csvPath)
	if err != nil {
		return err
	}

	newDevices := s.findNewDevices(csvDeviceIDs)

	return s.createNewDevices(newDevices)
}

// readDeviceIDsFromCSV reads all device IDs from the CSV file
func (s *CSVService) readDeviceIDsFromCSV(csvPath string) (map[string]bool, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("devices.csv file not found at %s: %w", csvPath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // Skip header row

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv: %w", err)
	}

	csvDeviceIDs := make(map[string]bool)
	for _, row := range records {
		if len(row) > 0 {
			deviceID := strings.TrimSpace(row[0])
			if deviceID != "" {
				csvDeviceIDs[deviceID] = true
			}
		}
	}

	return csvDeviceIDs, nil
}

// findNewDevices compares CSV devices with existing database devices
func (s *CSVService) findNewDevices(csvDeviceIDs map[string]bool) []models.DeviceModel {
	existingIDs := s.db.GetAllDeviceIDs()
	existingIDMap := make(map[string]bool)
	for _, id := range existingIDs {
		existingIDMap[id] = true
	}

	var newDevices []models.DeviceModel
	for id := range csvDeviceIDs {
		if !existingIDMap[id] {
			newDevices = append(newDevices, models.DeviceModel{DeviceID: id})
		}
	}

	return newDevices
}

// createNewDevices bulk inserts new devices into the database
func (s *CSVService) createNewDevices(newDevices []models.DeviceModel) error {
	if len(newDevices) == 0 {
		log.Println("No new devices to load from CSV.")
		return nil
	}

	if err := s.db.CreateDevices(newDevices); err != nil {
		return fmt.Errorf("failed to bulk insert new devices: %w", err)
	}

	log.Printf("Loaded %d new devices from CSV.", len(newDevices))
	return nil
}
