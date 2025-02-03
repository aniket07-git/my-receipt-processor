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

### Endpoints

- `POST /receipts/process`  
  Accepts a JSON payload representing a receipt. Returns a UUID to reference that receipt.

- `GET /receipts/{id}/points`  
  Returns the number of points that were awarded to the referenced receipt.

### Running Locally (with Go)

1. Ensure [Go 1.20+](https://go.dev/dl/) is installed.
2. Clone this repo.
3. `cd my-receipt-processor`
4. `go mod tidy` (to ensure all dependencies are installed).
5. `go build -o receipt-processor ./cmd/receipt-processor`
6. `./receipt-processor` (or `receipt-processor.exe` on Windows).

The server starts on port `8080`.

### Running with Docker

1. Ensure Docker is installed.
2. `docker build -t my-receipt-processor .`
3. `docker run -p 8080:8080 my-receipt-processor`

Now the API is accessible on `http://localhost:8080`.

### Example Usage

Using `curl` to test:

```bash
# POST a receipt
curl -X POST http://localhost:8080/receipts/process \
     -H 'Content-Type: application/json' \
     -d '{
           "retailer": "Target",
           "purchaseDate": "2022-01-01",
           "purchaseTime": "13:01",
           "items": [
             { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
             { "shortDescription": "Emils Cheese Pizza", "price": "12.25" },
             { "shortDescription": "Knorr Creamy Chicken", "price": "1.26" },
             { "shortDescription": "Doritos Nacho Cheese", "price": "3.35" },
             { "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "12.00" }
           ],
           "total": "35.35"
         }'

# Example response:
# {"id":"7fb1377b-b223-49d9-a31a-5a02701dd310"}

# GET the points
curl http://localhost:8080/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points
# {"points":28}




Receipts are stored in memory; data will be lost if the application restarts.

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
