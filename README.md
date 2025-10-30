# Go Shopping Cart

This repository contains a Go application for managing products and their prices, including functionalities for listing products, viewing product details, and seeding the database with initial data.

---

## 📁 Project Structure

- **`cmd/`**: Entry points
    - `server/main.go`: Starts the REST API
    - `seed/main.go`: Seeds the database

- **`app/`**: Application logic and handlers
- **`models/`**: Data models and repository interfaces
- **`sql/`**: Database schema and seed scripts
- **`.env`**: Configuration file for environment variables

---


## Environment Configuration

1. Copy the example `.env` file:
   ```bash
   cp .env.example .env

---

## 🚀 Application Setup

### Prerequisites

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/products/docker-desktop)

### Running the Project

```bash
make tidy         # Installs Go dependencies
make docker-up    # Starts the database container
make seed         # ⚠️ Resets and seeds the database with test data
make run          # Starts the Go server on http://localhost:8084
make test         # Runs unit tests
make docker-down  # Stops and removes docker containers
```

---

## 📦 API Endpoints

### `GET /catalog`

Returns a list of products with support for pagination and filtering.

#### Query Parameters:

| Param       | Type    | Description                                  |
|-------------|---------|----------------------------------------------|
| `offset`    | int     | Optional. Default is `0`.                    |
| `limit`     | int     | Optional. Default is `10`. Max is `100`.     |
| `category`  | string  | Optional. Filters products by category name. |
| `price_lt`  | float   | Optional. Filters products with price less than this value. |

#### Example:

```bash
curl "http://localhost:8084/catalog?offset=0&limit=5&category=Shoes&price_lt=20"
```

#### Sample Response:

```json
{
  "total": 1,
  "products": [
    {
      "code": "PROD007",
      "price": 18.2,
      "category": "Shoes"
    }
  ]
}
```

---

### `GET /catalog/:id`

Returns full details of a product including its variants and category.

#### Path Parameters:

| Param  | Type | Description                |
|--------|------|----------------------------|
| `id`   | int  | Required. Product ID.      |

#### Example:

```bash
curl "http://localhost:8084/catalog/1"
```

#### Sample Response:

```json
{
  "code": "PROD001",
  "price": 10.99,
  "category": "Clothing",
  "variants": [
    {
      "name": "Red",
      "sku": "SKU123",
      "price": 10.99
    },
    {
      "name": "Blue",
      "sku": "SKU124",
      "price": 12.99
    }
  ]
}
```

> **Note:** If a variant does not have a specific price, it will inherit the product's base price.
