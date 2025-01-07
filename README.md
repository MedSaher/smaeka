# E-Commerce Platform Prototype

## Overview
This is a prototype for an e-commerce platform developed using Go for backend logic and HTML with Bootstrap for the frontend. The project aims to demonstrate the integration of a robust backend with a responsive and user-friendly interface. It is designed to be a modular, scalable, and maintainable foundation for full-fledged e-commerce solutions.

---

## Features
- **User Authentication**: Secure user login and registration.
- **Product Management**: CRUD (Create, Read, Update, Delete) functionality for products.
- **Shopping Cart**: Add, update, and remove items from the cart.
- **Checkout System**: Simulated checkout process.
- **Responsive Design**: Frontend built with Bootstrap for cross-device compatibility.

---

## Tech Stack

| Component            | Technology           |
|----------------------|----------------------|
| Backend              | Go                  |
| Frontend             | HTML, Bootstrap     |
| Database (Optional)  | SQLite (can be extended) |

---

## Installation

### Prerequisites
1. Install [Go](https://golang.org/doc/install).
2. Install a code editor (e.g., VS Code).
3. Optionally, set up SQLite for database storage.

### Steps
1. Clone the repository:
   ```bash
   git clone <repository_url>
   ```
2. Navigate to the project directory:
   ```bash
   cd ecommerce-platform-prototype
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run main.go
   ```
5. Open your browser and visit `http://localhost:8080` to view the application.

---

## Project Structure

```plaintext
.
├── main.go           # Entry point for the application
├── controllers/      # Handles HTTP requests and responses
├── models/           # Contains data structures and business logic
├── templates/        # HTML templates for rendering pages
├── static/           # Static files (CSS, JS, images)
├── go.mod            # Go module file
└── README.md         # Project documentation
```

---

## How to Contribute
1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add your message here"
   ```
4. Push to the branch:
   ```bash
   git push origin feature/your-feature-name
   ```
5. Open a Pull Request.

---

## Future Enhancements
- **Payment Integration**: Add real payment gateways like Stripe or PayPal.
- **Admin Panel**: A dashboard for managing users, products, and orders.
- **API Support**: Expose REST APIs for third-party integrations.
- **Enhanced Security**: Implement HTTPS, OAuth, and data validation.

---

## License
This project is open-source and available under the MIT License.

---

## Contact
For inquiries or support, please reach out at (mohamed.saher.23@ump.ac.ma).

