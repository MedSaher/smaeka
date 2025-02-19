package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Define a struct for categories
type Category struct {
	ID    int    // json:"id"
	Name  string // json:"name"
	Image string // json:"image"
}

// Define a struct for products
type Product struct {
	ID    int     // json:"id"
	Name  string  // json:"name"
	Desc  string  // json:"desc"
	Price float64 // json:"price"
	Image string  // json:"image"
}

type Command struct {
	ID           int
	FullName     string
	City         string
	Email        string
	Number       int
	ProductID    int
	ProductName  string
	ProductPrice float64
	DateTime     *string
	Confirmed    string
}

// Global database variable
var db *sql.DB

// mod function to be used in templates
func mod(x, y int) int {
	return x % y
}

func add(x, y int) int {
	return x + y
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("session")
		log.Println("Session Value:", session) // Log the session value
		if err != nil || session != "admin" {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "You need to log in as an admin to access this page",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	var err error
	// Connect to the database
	db, err = sql.Open("mysql", "root:Saher2005@tcp(localhost:3306)/smaeka")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Ping the database to ensure the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	// Serve HTML files from the "templates" directory
	// Initialize the Gin router
	router := gin.Default()
	// Route to handle logout
	router.POST("/logout", func(c *gin.Context) {
		// Clear the session cookie
		c.SetCookie("session", "", -1, "/", "", false, true)
		// Redirect to login page or home page
		c.Redirect(http.StatusSeeOther, "/login")
	})
	// Add the custom function to the template engine
	router.SetFuncMap(template.FuncMap{
		"mod": mod,
		"add": add,
	})
	// Serve static files
	router.Static("/img", "./img")
	router.Static("/fonts", "./fonts")
	// Load the HTML templates from the templates directory
	router.LoadHTMLGlob("templates/*")

	// Define the route for the home page (index.html)
	router.GET("/", func(c *gin.Context) {
		// Fetch NewArrivals from the database
		newArrivals, err := getNewArrivals(db)
		if err != nil {
			log.Println("Error retrieving new arrivals:", err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		// Fetch Categories from the database
		categories, err := getCategories(db)
		if err != nil {
			log.Println("Error retrieving categories:", err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		// Render the template with NewArrivals and Categories data
		c.HTML(http.StatusOK, "index.html", gin.H{
			"NewArrivals": newArrivals,
			"Categories":  categories,
		})

	})
	router.GET("/about", func(c *gin.Context) {
		// Fetch Categories from the database
		categories, err := getCategories(db)
		if err != nil {
			log.Println("Error retrieving categories:", err)
			c.HTML(http.StatusInternalServerError, "about.html", gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		// Render the template with NewArrivals and Categories data
		c.HTML(http.StatusOK, "about.html", gin.H{
			"Categories": categories,
		})

	})
	router.GET("/contact", func(c *gin.Context) {
		// Fetch Categories from the database
		categories, err := getCategories(db)
		if err != nil {
			log.Println("Error retrieving categories:", err)
			c.HTML(http.StatusInternalServerError, "contact.html", gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		// Render the template with NewArrivals and Categories data
		c.HTML(http.StatusOK, "contact.html", gin.H{
			"Categories": categories,
		})

	})
	router.GET("/shop", func(c *gin.Context) {
		newArrivals, err := getNewArrivals(db)
		if err != nil {
			log.Println("Error retrieving new arrivals:", err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		// Fetch Categories from the database
		categories, err := getCategories(db)
		if err != nil {
			log.Println("Error retrieving categories:", err)
			c.HTML(http.StatusInternalServerError, "shop.html", gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		// Render the template with NewArrivals and Categories data
		c.HTML(http.StatusOK, "shop.html", gin.H{
			"Categories":  categories,
			"NewArrivals": newArrivals,
		})

	})

	// Apply the AdminAuthMiddleware to all routes starting with /admin
	adminRoutes := router.Group("/admin").Use(AdminAuthMiddleware())
	{
		adminRoutes.GET("", func(c *gin.Context) {
			categories, err := getCategories(db)
			if err != nil {
				log.Println("Error retrieving categories:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Internal Server Error",
				})
				return
			}

			products, err := getProducts(db)
			if err != nil {
				log.Println("Error retrieving products:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Internal Server Error",
				})
				return
			}

			commands, err := getCommands(db)
			if err != nil {
				log.Println("Error retrieving commands:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Internal Server Error",
				})
				return
			}

			c.HTML(http.StatusOK, "admin.html", gin.H{
				"Categories": categories,
				"Products":   products,
				"Commands":   commands,
			})
		})

		// Route to add a category
		adminRoutes.POST("/add-category", func(c *gin.Context) {
			categoryName := c.PostForm("categoryName")

			// Handle image upload
			file, err := c.FormFile("categoryImage")
			var fileName string
			if err == nil {
				fileName = filepath.Base(file.Filename)
				filePath := "./img/" + fileName
				if err := c.SaveUploadedFile(file, filePath); err != nil {
					c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
						"error": "Failed to save image: " + err.Error(),
					})
					return
				}
			} else {
				fileName = "" // No file uploaded
			}

			if categoryName == "" {
				c.HTML(http.StatusBadRequest, "admin.html", gin.H{
					"error": "Category name is required",
				})
				return
			}

			_, err = db.Exec("INSERT INTO categories (category_name, category_image) VALUES (?, ?)", categoryName, fileName)
			if err != nil {
				log.Println("Error inserting category:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Failed to add category",
				})
				return
			}

			// Redirect to refresh the category list
			c.Redirect(http.StatusSeeOther, "/admin")
		})

		// Route to delete a category
		adminRoutes.POST("/delete-category", func(c *gin.Context) {
			categoryID := c.PostForm("categoryID")

			if categoryID == "" {
				c.HTML(http.StatusBadRequest, "admin.html", gin.H{
					"error": "Category ID is required",
				})
				return
			}

			_, err := db.Exec("DELETE FROM categories WHERE id = ?", categoryID)
			if err != nil {
				log.Println("Error deleting category:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Failed to delete category",
				})
				return
			}

			// Redirect to refresh the category list
			c.Redirect(http.StatusSeeOther, "/admin")
		})

		// Route to add a product
		adminRoutes.POST("/add-product", func(c *gin.Context) {
			productName := c.PostForm("productName")
			productDesc := c.PostForm("productDesc")
			productPrice := c.PostForm("productPrice")
			productCategory := c.PostForm("productCategory")

			// Handle image upload
			file, err := c.FormFile("productImage")
			var fileName string
			if err == nil {
				fileName = filepath.Base(file.Filename)
				filePath := "./img/" + fileName
				if err := c.SaveUploadedFile(file, filePath); err != nil {
					c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
						"error": "Failed to save image",
					})
					return
				}
			}

			_, err = db.Exec("INSERT INTO products (product_name, product_desc, product_price, product_category, date_img) VALUES (?, ?, ?, ?, ?)",
				productName, productDesc, productPrice, productCategory, fileName)
			if err != nil {
				log.Println("Error inserting product:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Failed to add product",
				})
				return
			}

			c.Redirect(http.StatusSeeOther, "/admin")
		})

		// Route to handle command confirmation
		adminRoutes.POST("/confirm-command", func(c *gin.Context) {
			commandID := c.PostForm("commandID")

			// Update the command status to confirmed
			_, err := db.Exec("UPDATE commands SET confirmed = 'yes' WHERE command_id = ?", commandID)
			if err != nil {
				log.Println("Error confirming command:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Failed to confirm command",
				})
				return
			}

			// Redirect back to the commands page
			c.Redirect(http.StatusSeeOther, "/admin")
		})
		adminRoutes.POST("/delete-command", func(c *gin.Context) {
			commandID := c.PostForm("commandID")

			if commandID == "" {
				c.HTML(http.StatusBadRequest, "admin.html", gin.H{
					"error": "Command ID is required",
				})
				return
			}

			_, err := db.Exec("DELETE FROM commands WHERE command_id = ?", commandID)
			if err != nil {
				log.Println("Error deleting product:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Failed to delete product",
				})
				return
			}

			// Redirect to refresh the product list
			c.Redirect(http.StatusSeeOther, "/admin")
		})
		// Route to delete a product
		adminRoutes.POST("/delete-product", func(c *gin.Context) {
			productID := c.PostForm("productID")

			if productID == "" {
				c.HTML(http.StatusBadRequest, "admin.html", gin.H{
					"error": "Product ID is required",
				})
				return
			}

			_, err := db.Exec("DELETE FROM products WHERE product_id = ?", productID)
			if err != nil {
				log.Println("Error deleting product:", err)
				c.HTML(http.StatusInternalServerError, "admin.html", gin.H{
					"error": "Failed to delete product",
				})
				return
			}

			// Redirect to refresh the product list
			c.Redirect(http.StatusSeeOther, "/admin")
		})
	}

	// Route for login (not implemented in this example, but you should handle it)
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "",
		})
	})

	// Route to handle login POST request
	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// var dbPassword string
		var userCount int

		// Query to get the number of users with the provided username and password
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND password = ?", username, password).Scan(&userCount)
		if err != nil {
			log.Println("Error querying user:", err)
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Check if exactly one user matches
		if userCount == 1 {
			// Create a session cookie
			c.SetCookie("session", "admin", 3600, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/admin")
		} else {
			// Handle invalid login
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "Invalid username or password",
			})
		}
	})
	// / Route to handle form submission
	router.POST("/submit-command", func(c *gin.Context) {
		productID := c.PostForm("productID")
		productName := c.PostForm("productName")
		productPrice := c.PostForm("productPrice")
		fullname := c.PostForm("fullname")
		city := c.PostForm("city")
		email := c.PostForm("email")
		number := c.PostForm("number")

		// Insert the command into the database
		_, err := db.Exec("INSERT INTO commands (commander_fullname, city, email, number, product_id, datetime, confirmed) VALUES (?, ?, ?, ?, ?, ?,?)",
			fullname, city, email, number, productID, nil, "not")
		if err != nil {
			log.Println("Error inserting command:", err)
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "Failed to place order",
			})
			return
		}

		// Redirect to the confirmation page with order details
		c.Redirect(http.StatusSeeOther, "/confirmation?productID="+productID+"&productName="+productName+"&productPrice="+productPrice)
	})
	// Route to display products by category
	router.GET("/category/:id", CategoryPage)
	// Route to display product details
	router.GET("/product/:id", ProductPage)
	// Route to display the confirmation page
	router.GET("/confirmation", ConfirmationPage)
	// router.GET("/about", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "about.html", nil)
	// })

	router.GET("/faq", func(c *gin.Context) {
		c.HTML(http.StatusOK, "faq.html", nil)
	})
	router.GET("/privacy", func(c *gin.Context) {
		c.HTML(http.StatusOK, "privacy.html", nil)
	})
	router.GET("/terms", func(c *gin.Context) {
		c.HTML(http.StatusOK, "terms.html", nil)
	})

	// Start the server
	router.Run()
}

// Get all categories from the database
func getCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT id, category_name, category_image FROM categories")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Image); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}

	return categories, nil
}

// Get all products from the database
func getProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT product_id, product_name, product_desc, product_price, date_img FROM products")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Desc, &p.Price, &p.Image); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}

	return products, nil
}

// Get new arrivals from the database
func getNewArrivals(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT product_id, product_name, product_desc, product_price, date_img FROM products ORDER BY date_add DESC LIMIT 10")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Desc, &p.Price, &p.Image); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}

	return products, nil
}

// Handler to render the category page
func CategoryPage(c *gin.Context) {
	categoryID := c.Param("id") // Get the category ID from the URL parameter
	log.Println("Category ID:", categoryID)
	// Query the database to get the category details and products
	var category Category
	var products []Product

	// Query the database for the category details
	err := db.QueryRow("SELECT id, category_name, category_image FROM categories WHERE id = ?", categoryID).Scan(&category.ID, &category.Name, &category.Image)
	if err != nil {
		log.Println("Error retrieving category:", err)
		c.HTML(http.StatusInternalServerError, "category.html", gin.H{
			"error": "Failed to retrieve category",
		})
		return
	}

	// Query the database for the products in this category
	rows, err := db.Query("SELECT product_id, product_name, product_desc, product_price, date_img FROM products WHERE product_category = ?", categoryID)
	if err != nil {
		log.Println("Error retrieving products:", err)
		c.HTML(http.StatusInternalServerError, "category.html", gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Desc, &p.Price, &p.Image); err != nil {
			log.Println("Error scanning product row:", err)
			c.HTML(http.StatusInternalServerError, "category.html", gin.H{
				"error": "Failed to retrieve products",
			})
			return
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with product rows:", err)
		c.HTML(http.StatusInternalServerError, "category.html", gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	// Render the category page template with the category and products data
	c.HTML(http.StatusOK, "category.html", gin.H{
		"Category": category,
		"Products": products,
	})
}

// Handler to render the product page
func ProductPage(c *gin.Context) {
	productID := c.Param("id") // Get the product ID from the URL

	// Query the database to get the product details
	var product Product

	err := db.QueryRow("SELECT product_id, product_name, product_desc, product_price, date_img FROM products WHERE product_id = ?", productID).Scan(&product.ID, &product.Name, &product.Desc, &product.Price, &product.Image)
	if err != nil {
		log.Println("Error retrieving product:", err)
		c.HTML(http.StatusInternalServerError, "product.html", gin.H{
			"error": "Failed to retrieve product",
		})
		return
	}

	// Render the product page template with the product data
	c.HTML(http.StatusOK, "product.html", gin.H{
		"Product": product,
	})
}

func ConfirmationPage(c *gin.Context) {
	productName := c.Query("productName")
	productPrice := c.Query("productPrice")
	c.HTML(http.StatusOK, "confirmation.html", gin.H{
		"ProductName":  productName,
		"ProductPrice": productPrice,
	})
}

func getCommands(db *sql.DB) ([]Command, error) {
	rows, err := db.Query(`
		SELECT c.command_id, c.commander_fullname, c.city, c.email, c.number, 
			   c.product_id, p.product_name, p.product_price, c.datetime, c.confirmed 
		FROM commands c
		JOIN products p ON c.product_id = p.product_id
	`)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		if err := rows.Scan(&cmd.ID, &cmd.FullName, &cmd.City, &cmd.Email, &cmd.Number,
			&cmd.ProductID, &cmd.ProductName, &cmd.ProductPrice,
			&cmd.DateTime, &cmd.Confirmed); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		// If DateTime is nil, set it to "NULL"
		if cmd.DateTime == nil {
			nullStr := "NULL"
			cmd.DateTime = &nullStr
		}

		commands = append(commands, cmd)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}

	return commands, nil
}
