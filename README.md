# Restaurant Order Management System

This is a backend API for managing restaurant operations, built with Go, Gin and PostgreSQL. It handles user authentication, role-based access, product and category management, and order tracking. 

I structured it as a clean monolith so it is easy to reason about, but the code is organized well enough that splitting it into services later would not be a big deal.

## What It Does

The system uses JWT for authentication. Users log in and get a token, and that token is used for every request after that. 

There are two main roles:
* **Admin:** Admins have full access to everything in the system.
* **Waiter:** Waiters can browse products and categories, and create and manage their own orders. They cannot see other waiters' orders or touch anything they are not supposed to.

Here is a breakdown of the core features:
* **Categories and Products:** Admins have full control — they can create, edit, and delete. Waiters can only read. The product list also supports search by name, filtering by category, and pagination, so it is easy to work with even when there are a lot of items.
* **Orders:** Orders are the core of the system. When an order is created, it automatically reserves stock for the items, calculates the subtotal, applies a 10% tax, and gives you the grand total. Order creation is wrapped in a database transaction so nothing gets saved halfway. To keep things simple, orders just use a two-step status flow: `pending` and `completed`. If you ever need to cancel an order, simply delete it, and the system will automatically return the item quantities back into the product stock when order status is pending! Admins can see and manage all orders, while waiters only see their own.
* **User Management:** This is admin-only. Admins can create new users, assign roles, update details, or remove accounts.

Every endpoint returns responses in the same consistent JSON format, whether it is a success, an error, or a paginated list.

## Tech Stack

Here is what is running under the hood:
* **Go 1.21** as the core language
* **Gin** as the web framework
* **golang-jwt** for authentication
* **bcrypt** for password hashing
* **PostgreSQL 16** for the database
* **Docker and Docker Compose** so you can run the whole thing without setting up anything locally

## Prerequisites

Before you begin, make sure you have the following installed:
* **Go** (version 1.21 or higher)
* **PostgreSQL** (version 16 or higher)
* **Docker and Docker Compose** (optional, if you want to run it via containers)

## Environment Variables

You need to set up a few environment variables to connect to the database and secure the JWT tokens. Create a `.env` file in the root of the project with the following:

DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=restaurant_db
DB_PORT=5432
JWT_SECRET=your_super_secret_key

## Running the Project

The easiest way to get it running is with Docker. As long as you have Docker and Docker Compose installed, just follow these steps:

1. Clone the repo: `git clone https://github.com/Deepak-S-M/restaurant-order-management.git`
2. Move into the directory: `cd restaurant-order-management`
3. Spin it up: `docker-compose up --build`

That will start both the database and the app together. The API will be available at http://localhost:8080 once everything is up. 

If you prefer to run it without Docker, you will need Go 1.21 or later and a PostgreSQL instance running locally:
1. Start by cloning the repo.
2. Create a database called `restaurant_db`. 
3. Copy `.env.example` to `.env` and fill in your database credentials and a JWT secret. 
4. Run `go mod tidy` to pull in the dependencies.
5. Run `go run .` to start the server.

On the first startup, the app will auto-migrate the database schema and seed some default data, so you do not need to set anything up manually.

## Default Seeded Data

When the server starts for the first time, it sets up some data to get you going:
* **Roles:** It creates admin and waiter roles.
* **Users:** It creates two users: 
  * `admin@restaurant.com` (password: admin123)
  * `waiter@restaurant.com` (password: waiter123)
* **Categories:** It creates Appetizers, Main Course, Desserts, and Beverages.
* **Products:** It adds a couple of sample products under Main Course, priced realistically in Indian Rupees (INR).

## API Reference

All endpoints are under the `/api` prefix.

**Authentication**
* `POST /api/login` — accepts an email and password and returns a JWT token. 
* `POST /api/register` — creates a new account and defaults the user to the waiter role.

**Categories**
* `GET /api/categories` — list all categories (admins and waiters)
* `GET /api/categories/:id` — get a specific category (admins and waiters)
* `POST /api/categories` — create a category (admin only)
* `PUT /api/categories/:id` — update a category (admin only)
* `DELETE /api/categories/:id` — delete a category (admin only)

**Products**
The list endpoint supports optional query parameters: page, limit, search (matches product name), and category_id.
* `GET /api/products` — list products (admins and waiters)
* `GET /api/products/:id` — get a single product (admins and waiters)
* `POST /api/products` — create a product (admin only)
* `PUT /api/products/:id` — update a product (admin only)
* `DELETE /api/products/:id` — delete a product (admin only)

**Orders**
You can filter the list endpoint by status or page through results.
* `GET /api/orders` — list orders (admins see all, waiters see their own)
* `POST /api/orders` — create a new order and automatically deduct stock (admins and waiters)
* `GET /api/orders/:id` — get full order details (admins and waiters)
* `PUT /api/orders/:id/status` — update the status to either `pending` or `completed` (admins and waiters)
* `DELETE /api/orders/:id` — remove an order and restore its stock back to the inventory (admin only)

**Users**
Everything here is admin-only.
* `GET /api/users` — list all users
* `POST /api/users` — create a user with a specific role
* `GET /api/users/:id` — get a single user
* `PUT /api/users/:id` — update a user
* `DELETE /api/users/:id` — delete an account

**Roles**
Everything here is admin-only.
* `GET /api/roles` — list all available roles
* `GET /api/roles/:id` — get details for a specific role

## Sample Requests & Responses

Here is an example of what it looks like to interact with the API when creating an order:

**Creating an Order (Request)**
```json
POST /api/orders
{
  "items": [
    {
      "product_id": "bc07baeb-d8dd-43db-97ec-8991d3df812b",
      "quantity": 3
    }
  ]
}

**Creating an Order (Response)**
json
{
  "status": "success",
  "message": "Order created successfully",
  "data": {
    "id": "849592e2-3c5b-4939-81a0-4d78c038c578",
    "user_id": "2cdaa24a-b37e-4225-873b-5f479fe0a8dd",
    "status": "pending",
    "subtotal": 897.00,
    "tax": 89.70,
    "grand_total": 986.70,
    "items": [
      {
        "product_id": "bc07baeb-d8dd-43db-97ec-8991d3df812b",
        "quantity": 3,
        "unit_price": 299.00,
        "subtotal": 897.00
      }
    ]
  }
}

## Testing with Postman

There is a Postman collection in the root of the project called `postman_collection.json`. 

1. Import it into Postman.
2. Go to the collection variables and set `base_url` to `http://localhost:8080/api`. 
3. Call the Login endpoint with one of the seeded credentials.
4. Copy the token from the response, and paste it into the `token` variable. 

After that, all the other requests in the collection will automatically include it in the authorization header.

## Project Structure

The codebase is organized pretty straightforwardly:
* **config** — handles the database connection.
* **controllers** — where all the request logic lives (one file per resource).
* **middlewares** — contains the JWT auth check and role guard that protect the routes. 
* **models** — database models.
* **utils** — shared utility functions like response formatting.

At the root level:
* **router.go** is where all the routes are registered.
* **seeder.go** handles the default data on first boot.
* **main.go** is the entry point that ties everything together. 
* **Dockerfile** and **docker-compose.yml** are there if you want to run the whole thing with Docker.

## Future Development

These are features planned for future versions of the system.

### Table Management
- Add `table_number` field to orders so each order is linked to a specific table
- Waiter can add extra items to an existing order for the same table without creating a new order
- System automatically recalculates subtotal, tax, and grand total when new items are added
- Admin can fetch full order details by table number for billing at checkout