package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"JSON-File-DB/tinydb"
)

const (
	ColProducts  = "products"
	ColSuppliers = "suppliers"
)

type Product struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

type Supplier struct {
	CompanyName string `json:"company_name"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
}

func main() {
	db, err := tinydb.New("database.json")
	if err != nil {
		log.Fatalf("Fatal error starting database: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== Warehouse Inventory System ===")
		fmt.Println("1. Add New Product")
		fmt.Println("2. Add New Supplier")
		fmt.Println("3. View a Record (by ID)")
		fmt.Println("4. Update a Record")
		fmt.Println("5. Delete a Record")
		fmt.Println("6. Exit")
		fmt.Print("Select an option (1-6): ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			addProduct(db, scanner)
		case "2":
			addSupplier(db, scanner)
		case "3":
			viewRecord(db, scanner)
		case "4":
			updateRecord(db, scanner)
		case "5":
			deleteRecord(db, scanner)
		case "6":
			fmt.Println("Shutting down system. Goodbye.")
			return
		default:
			fmt.Println("Invalid option. Please enter a number between 1 and 6.")
		}
	}
}

func addProduct(db *tinydb.DB, scanner *bufio.Scanner) {
	fmt.Println("\n--- Add Product ---")

	name := prompt(scanner, "Enter Product Name: ")
	category := prompt(scanner, "Enter Category: ")

	priceStr := prompt(scanner, "Enter Price (e.g., 19.99): ")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Println("Error: Price must be a valid number. Registration cancelled.")
		return
	}

	stockStr := prompt(scanner, "Enter Initial Stock Quantity: ")
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		fmt.Println("Error: Stock must be a whole number. Registration cancelled.")
		return
	}

	newProduct := Product{
		Name:     name,
		Category: category,
		Price:    price,
		Stock:    stock,
	}

	id, err := db.Insert(ColProducts, newProduct)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return
	}
	fmt.Printf("Success! Product registered with ID: %s\n", id)
}

func addSupplier(db *tinydb.DB, scanner *bufio.Scanner) {
	fmt.Println("\n--- Add Supplier ---")

	company := prompt(scanner, "Enter Company Name: ")
	contact := prompt(scanner, "Enter Contact Person: ")
	phone := prompt(scanner, "Enter Phone Number: ")

	newSupplier := Supplier{
		CompanyName: company,
		ContactName: contact,
		Phone:       phone,
	}

	id, err := db.Insert(ColSuppliers, newSupplier)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return
	}
	fmt.Printf("Success! Supplier registered with ID: %s\n", id)
}

func viewRecord(db *tinydb.DB, scanner *bufio.Scanner) {
	fmt.Println("\n--- View Record ---")

	collection, valid := getValidCollection(scanner)
	if !valid {
		return
	}

	id := prompt(scanner, "Enter Record ID: ")

	data, err := db.Read(collection, id)
	if err != nil {
		fmt.Printf("Error finding record: %v\n", err)
		return
	}

	fmt.Printf("\nRecord Found in '%s':\n%+v\n", collection, data)
}

func updateRecord(db *tinydb.DB, scanner *bufio.Scanner) {
	fmt.Println("\n--- Update Record ---")

	collection, valid := getValidCollection(scanner)
	if !valid {
		return
	}

	id := prompt(scanner, "Enter Record ID to update: ")

	// Verify the record actually exists before asking for new data
	_, err := db.Read(collection, id)
	if err != nil {
		fmt.Printf("Error: Record not found (%v)\n", err)
		return
	}

	if collection == ColProducts {
		name := prompt(scanner, "Enter New Product Name: ")
		category := prompt(scanner, "Enter New Category: ")

		priceStr := prompt(scanner, "Enter New Price (e.g., 19.99): ")
		price, _ := strconv.ParseFloat(priceStr, 64)

		stockStr := prompt(scanner, "Enter New Stock Quantity: ")
		stock, _ := strconv.Atoi(stockStr)

		updatedProduct := Product{
			Name:     name,
			Category: category,
			Price:    price,
			Stock:    stock,
		}

		err = db.Update(collection, id, updatedProduct)
	} else if collection == ColSuppliers {
		company := prompt(scanner, "Enter New Company Name: ")
		contact := prompt(scanner, "Enter New Contact Person: ")
		phone := prompt(scanner, "Enter New Phone Number: ")

		updatedSupplier := Supplier{
			CompanyName: company,
			ContactName: contact,
			Phone:       phone,
		}

		err = db.Update(collection, id, updatedSupplier)
	}

	if err != nil {
		fmt.Printf("Database error: %v\n", err)
	} else {
		fmt.Println("Success! Record updated.")
	}
}

func deleteRecord(db *tinydb.DB, scanner *bufio.Scanner) {
	fmt.Println("\n--- Delete Record ---")

	collection, valid := getValidCollection(scanner)
	if !valid {
		return
	}

	id := prompt(scanner, "Enter Record ID to delete: ")

	err := db.Delete(collection, id)
	if err != nil {
		fmt.Printf("Error: Failed to delete record (%v)\n", err)
	} else {
		fmt.Println("Success! Record permanently deleted.")
	}
}

// prompt prints a message and waits for the user to type a response
func prompt(scanner *bufio.Scanner, message string) string {
	fmt.Print(message)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func getValidCollection(scanner *bufio.Scanner) (string, bool) {
	col := prompt(scanner, fmt.Sprintf("Enter Collection (%s or %s): ", ColProducts, ColSuppliers))
	col = strings.ToLower(col)

	if col == ColProducts || col == ColSuppliers {
		return col, true
	}

	fmt.Printf("Error: Invalid collection. Must be '%s' or '%s'.\n", ColProducts, ColSuppliers)
	return "", false
}
