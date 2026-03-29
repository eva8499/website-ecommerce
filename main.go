package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "modernc.org/sqlite"
)

type Product struct {
	Name     string
	Price    int
	Image    string
	Category string
}

var products = []Product{
	{"Jeruk bali", 45000, "/static/img/jeruk Bali.jpg", "buah"},
	{"Buah naga merah 1 KG", 80000, "/static/img/Buah Naga.jpg", "buah"},
	{"Jeruk Sunkist Navel 500 gr", 40000, "/static/img/jeruk.jpg", "buah"},
	{"Pepaya California", 25000, "/static/img/Pepaya.jpg", "buah"},
	{"Pear Century 1 KG", 44000, "/static/img/pir.jpg", "buah"},
	{"Anggur Red Globe ", 60000, "/static/img/anggur.jpg", "buah"},
	{"Durian", 50000, "/static/img/durian.jpg", "buah"},
	{"Kelengkeng Thailand 1 KG", 45000, "/static/img/kelengkeng.jpg", "buah"},
	{"Nanas madu per buah", 15000, "/static/img/nanas.jpg", "buah"},
	{"Semangka merah per buah", 42000, "/static/img/semangka.jpg", "buah"},
	{"Strawberry Fruit Fresh 250gr", 15000, "/static/img/stroberi.jpg", "buah"},
	{"Delima Merah Jumbo ", 18000, "/static/img/Delima.jpg", "buah"},

	{"Cider blus putih bordir bunga", 100000, "/static/img/baju 1.jpg", "baju"},
	{"Maison Special Cotton Fabric Lace Colored Cardigan", 150000, "/static/img/baju 2.jpg", "baju"},
	{"Cardingan rajut crop wanita", 80000, "/static/img/baju 3.jpg", "baju"},
	{"Colorbox Contrast Details Denim Jacket off-white", 200000, "/static/img/baju 4.jpg", "baju"},
	{"Barrie Contrast-Stitching Denim-Effect Jacket", 95000, "/static/img/baju 5.jpg", "baju"},
	{"Collar flower embroidery cardigan", 85000, "/static/img/baju 6.jpg", "baju"},
	{"Kaos Polo Distro ", 50000, "/static/img/baju 7.jpg", "baju"},
	{"Sweater rugby lengan panjang", 160000, "/static/img/baju 8.jpg", "baju"},
	{"Abercrombie & Fitch Men's Tipped Johnny Collar Sweater Polo. ", 100000, "/static/img/baju 9.jpg", "baju"},
	{"Abercrombie & Fitch Men's Essential Rugby Polo Sweatshirt. ", 110000, "/static/img/baju 10.jpg", "baju"},
	{"Striped Color Blok Sweater for men", 100000, "/static/img/baju 11.jpg", "baju"},
	{"Zara Striped Knit Sweater", 140000, "/static/img/baju 12.jpg", "baju"},

	{"Rinnai kompor gas model RI-202S 2 tungku", 450000, "/static/img/KOMPOR GAS.jpg", "elektronik"},
	{"Kipas Angin Portable Mini Fan Karakter Sanrio", 35000, "/static/img/KIPAS ANGIN MINI.jpg", "elektronik"},
	{"Anker Soundcore Space One Pro", 250000, "/static/img/aerphone.jpg", "elektronik"},
	{"Miyako Blender BL-151 GF blender kaca 2-in-1", 300000, "/static/img/BLENDER.jpg", "elektronik"},
	{"Lunalife Sterika baju uap", 190000, "/static/img/strika.jpg", "elektronik"},
	{"Printer HP Smart Tank 580 All-in-One", 2300000, "/static/img/printer.jpg", "elektronik"},
	{" Mecoo Belle Air Fryer 4L", 950000, "/static/img/Air fryer.jpg", "elektronik"},
	{"Philco Planetária PHP500 Turbo", 500000, "/static/img/mixer.jpg", "elektronik"},
	{"oven listrik mini HAN RIVER kapasitas 12 liter. ", 1500000, "/static/img/oven.jpg", "elektronik"},
	{"Maspion Rice Cooker 1 Liter MRJ-1003 BS", 370000, "/static/img/rice cooker.jpg", "elektronik"},
	{"Tramontina 3-Piece Black Aluminum Frying Pan Set", 500000, "/static/img/teflon.jpg", "elektronik"},
	{"Sharp LED Digital TV 32 Inch 2T-C32GH3000i", 2600000, "/static/img/tv.jpg", "elektronik"},
}

var tmpl = template.Must(template.ParseGlob("templates/*.html"))
var db *sql.DB

func main() {

	var err error
	db, err = sql.Open("sqlite", "ecommerce.db")
	if err != nil {
		panic(err)
	}
	createTables()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/produk", produkHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/hapus", hapusHandler)

	http.HandleFunc("/tambah", tambahHandler)
	http.HandleFunc("/keranjang", keranjangHandler)

	fmt.Println("Server jalan di http://localhost:4000")

	http.ListenAndServe(":4000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var username string

	cookie, err := r.Cookie("user")
	if err == nil {
		username = cookie.Value
	}

	data := map[string]string{
		"Username": username,
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbUser, dbPass string

		err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).
			Scan(&dbUser, &dbPass)

		if err != nil || dbPass != password {
			tmpl.ExecuteTemplate(w, "login.html", "Login gagal!")
			return
		}
		cookie := &http.Cookie{
			Name:  "user",
			Value: username,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func getUser(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		_, err := db.Exec("INSERT INTO users(username, password) VALUES (?, ?)", username, password)
		if err != nil {
			fmt.Println("ERROR DB:", err)
			http.Error(w, "Gagal register", 500)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "register.html", nil)
}

func produkHandler(w http.ResponseWriter, r *http.Request) {
	kategori := r.URL.Query().Get("kategori")

	var username string
	cookie, err := r.Cookie("user")
	if err == nil {
		username = cookie.Value
	}

	var hasil []map[string]string

	for _, p := range products {
		if p.Category == kategori {
			item := map[string]string{
				"Name":  p.Name,
				"Image": p.Image,
				"Price": formatRupiah(p.Price),
			}
			hasil = append(hasil, item)
		}
	}

	data := map[string]interface{}{
		"Products": hasil,
		"Username": username,
	}

	tmpl.ExecuteTemplate(w, "produk.html", data)
}

var cart []Product

func tambahHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	nama := r.URL.Query().Get("nama")

	var harga int
	for _, p := range products {
		if p.Name == nama {
			harga = p.Price
		}
	}
	db.Exec("INSERT INTO cart(username, product_name, price) VALUES (?, ?, ?)",
		user, nama, harga)

	http.Redirect(w, r, "/keranjang", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func keranjangHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, _ := db.Query("SELECT product_name, price FROM cart WHERE username = ?", user)
	defer rows.Close()

	var items []map[string]string

	for rows.Next() {
		var name string
		var price int
		rows.Scan(&name, &price)

		items = append(items, map[string]string{
			"Name":  name,
			"Price": formatRupiah(price),
		})
	}
	tmpl.ExecuteTemplate(w, "keranjang.html", items)
}

func formatRupiah(price int) string {
	str := fmt.Sprintf("%d", price)
	n := len(str)

	var result string
	count := 0

	for i := n - 1; i >= 0; i-- {
		result = string(str[i]) + result
		count++

		if count%3 == 0 && i != 0 {
			result = "." + result
		}
	}

	return result
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		password TEXT
	);`

	cartTable := `
	CREATE TABLE IF NOT EXISTS cart (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		product_name TEXT,
		price INTEGER
	);`

	db.Exec(userTable)
	db.Exec(cartTable)
}

func hapusHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	nama := r.URL.Query().Get("nama")

	_, err := db.Exec("DELETE FROM cart WHERE username = ? AND product_name = ?", user, nama)
	if err != nil {
		http.Error(w, "Gagal hapus", 500)
		return
	}

	http.Redirect(w, r, "/keranjang", http.StatusSeeOther)
}
