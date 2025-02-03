# Receipt Processor

A simple Go API that processes receipts and calculates reward points using these rules:

1. One point for every alphanumeric character in the retailer name.
2. 50 points if the total is a round dollar amount with no cents.
3. 25 points if the total is a multiple of `0.25`.
4. 5 points for every two items on the receipt.
5. If the (trimmed) item description length is a multiple of 3, multiply the price by `0.2`, round up, and add that to points.
6. (LLM rule) 5 points if the total is greater than 10.00.
7. 6 points if the day of the purchase date is odd.
8. 10 points if the time of purchase is after 2:00pm and before 4:00pm.

Receipts are stored in memory; data will be lost if the application restarts.

## How to Start Docker
docker build -t receipt-processor .
docker run -p 8080:8080 receipt-processor

## API Endpoints

1. **POST** `/receipts/process`  
   - Accepts a JSON receipt payload.  
   - Returns a JSON object containing an `id` for the receipt.  
   - Example response:  
     ```json
     { "id": "fd7c4f42-0b2f-4778-b932-7c643d026be8" }
     ```

2. **GET** `/receipts/{id}/points`  
   - Returns the JSON object containing the total points awarded.  
   - Example response:  
     ```json
     { "points": 28 }
     ```

See [`api.yml`](./api.yml) for an OpenAPI 3.0 definition of this service.

---

## Running Locally (without Docker)

1. **Install Go** (version 1.18+ recommended).  
2. **Clone or copy** this folder to your machine.  
3. From the root folder `my-receipt-processor/`, run:
   ```bash
   go mod tidy
   go build -o receipt-processor ./cmd/receipt-processor
   ./receipt-processor
