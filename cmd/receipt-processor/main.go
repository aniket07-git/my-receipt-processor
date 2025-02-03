package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Receipt represents the JSON structure of a receipt.
type Receipt struct {
	Retailer     string        `json:"retailer"`
	PurchaseDate string        `json:"purchaseDate"`
	PurchaseTime string        `json:"purchaseTime"`
	Items        []ReceiptItem `json:"items"`
	Total        string        `json:"total"`
}

// ReceiptItem represents a single item in the receipt.
type ReceiptItem struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// We'll store the points in memory by a map of receiptID -> points
var receiptPoints = make(map[string]int)

func main() {
	r := mux.NewRouter()

	// POST /receipts/process
	r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")

	// GET /receipts/{id}/points
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// processReceiptHandler handles POST /receipts/process
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("processReceiptHandler invoked!")
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid receipt payload", http.StatusBadRequest)
		return
	}

	// Compute points
	points, err := computePoints(&receipt)
	if err != nil {
		http.Error(w, "Cannot compute points: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Generate an ID
	id := uuid.New().String()

	// Store points in memory
	receiptPoints[id] = points

	// Return ID
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// getPointsHandler handles GET /receipts/{id}/points
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	points, found := receiptPoints[id]
	if !found {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{"points": points})
}

// computePoints calculates the points for a given receipt.
func computePoints(r *Receipt) (int, error) {
	var totalPoints int

	// 1) One point for every alphanumeric character in the retailer name.
	re := regexp.MustCompile("[A-Za-z0-9]")
	matches := re.FindAllString(r.Retailer, -1)
	totalPoints += len(matches)

	// Parse total as float
	totalFloat, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid total value")
	}

	// 2) 50 points if the total is a round dollar amount with no cents
	if isRoundDollar(totalFloat) {
		totalPoints += 50
	}

	// 3) 25 points if the total is a multiple of 0.25
	if isMultipleOfQuarter(totalFloat) {
		totalPoints += 25
	}

	// 4) 5 points for every two items
	totalPoints += (len(r.Items) / 2) * 5

	// 5) If trimmed length of item description is a multiple of 3,
	//    multiply price by 0.2, round up, add that to points
	for _, item := range r.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		itemPrice, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			// skip invalid price
			continue
		}
		if len(desc)%3 == 0 && len(desc) != 0 {
			bonus := math.Ceil(itemPrice * 0.2)
			totalPoints += int(bonus)
		}
	}

	// 6) (LLM rule) +5 points if total > 10.00
	if totalFloat > 10.00 {
		totalPoints += 5
	}

	// 7) +6 points if purchase day is odd
	purchaseDay, err := parseDay(r.PurchaseDate)
	if err == nil {
		if purchaseDay%2 == 1 {
			totalPoints += 6
		}
	}

	// 8) +10 points if purchase time is after 2:00pm and before 4:00pm
	if isTimeBetween(r.PurchaseTime, "14:00", "16:00") {
		totalPoints += 10
	}

	return totalPoints, nil
}

// isRoundDollar checks if the float is an integer (e.g. 9.0, 10.0, etc.).
func isRoundDollar(value float64) bool {
	return value == float64(int(value))
}

// isMultipleOfQuarter checks if the float is a multiple of 0.25.
func isMultipleOfQuarter(value float64) bool {
	epsilon := 0.000001
	remainder := math.Mod(value, 0.25)
	return remainder < epsilon || math.Abs(remainder-0.25) < epsilon
}

// parseDay extracts the day integer from "YYYY-MM-DD".
func parseDay(dateStr string) (int, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, err
	}
	return t.Day(), nil
}

// isTimeBetween checks if timeStr is strictly between startStr and endStr (24h).
func isTimeBetween(timeStr, startStr, endStr string) bool {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false
	}
	start, _ := time.Parse("15:04", startStr)
	end, _ := time.Parse("15:04", endStr)
	return t.After(start) && t.Before(end)
}
