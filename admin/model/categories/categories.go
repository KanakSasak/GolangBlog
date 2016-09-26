/*	-- Categories Model --
 * 	All functions in this file are called by the corresponding controller or by
 	functions from itself.
 */
package categories

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	cfg "github.com/jschalkwijk/GolangBlog/admin/config"
)

/* Category struct will hold data about a category and can be added to the Data struct */
type Category struct {
	Category_ID int
	Title string
	Description string
	Content string
	Keywords string
	Approved int
	Author string
	Cat_Type string
	Date string
	Parent_ID int
	Trashed int
}

/*
 * Declaring vars corresponding to the struct. When scanning data from the database, the
   data will be stored on the memory address of these vars.
*/
var category_id int
var title string
var description string
var content string
var keywords string
var approved int
var author string
var cat_type string
var date string
var parent_id int
var trashed int

/* Stores a single category, or multiple categories in a Slice which can be iterated over in the template */
type Data struct {
	Categories []Category
	Deleted bool
}

/* -- Get all categories --
 * 	Connects to the database and gets all category rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	For every row get the values and set the values to the memory address of the named variable.
 		- Instantiate a new Category Struct and insert values.
 		- Append the category struct to the Data.Categories Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func GetCategories(trashed int) *Data {
	db, err := sql.Open("mysql", cfg.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	rows, err := db.Query("SELECT * FROM categories WHERE trashed = ? ORDER BY categorie_id DESC", trashed)
	checkErr(err)

	data := new(Data)

	for rows.Next() {
		err = rows.Scan(&category_id, &title, &description, &content,&keywords,&approved,
		&author,&cat_type,&date,&parent_id,&trashed)
		checkErr(err)

		category := Category{category_id,title,description,content,keywords,approved,author,cat_type,date,parent_id,trashed}

		data.Categories = append(data.Categories , category)
	}

	if(trashed == 1) {
		data.Deleted = true
	} else {
		data.Deleted = false
	}

	return data
}

/* -- Get a single categories -- */
/* GetSingleCategory gets a category from the DB and returns a pointer to the Struct. It takes a id and category_title.
 * 	Connects to the database and gets all category rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	Get a single row from the DB and get the values and set the values to the memory address of the named variable.
 *	Instantiate a new Category Struct and insert values.
 *	Append the category struct to the Data.Categories Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func GetSingleCategory(id string,category_title string) *Data {
	db, err := sql.Open("mysql", cfg.DB)
	checkErr(err)
	fmt.Println("Connection established")
	defer db.Close()
	defer fmt.Println("Connection Closed")

	rows := db.QueryRow("SELECT * FROM categories WHERE categorie_id=? AND title=? LIMIT 1", id,category_title)

	data := new(Data)

	err = rows.Scan(&category_id, &title, &description, &content,&keywords,&approved,
	&author,&cat_type,&date,&parent_id,&trashed)
	checkErr(err)

	category := Category{category_id,title,description,content,keywords,approved,author,cat_type,date,parent_id,trashed}

	data.Categories = append(data.Categories , category)

	//fmt.Println(collection.categories)
	return data
}

/* -- Category Methods -- */

/* saveCategory updates the values of an existing category to the database and is a method to Category
 * Called by EditCategory
 * Connect to the DB and prepares query.
 * Execute query with the inserted struct values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (p *Category) saveCategory() error {
	db, err := sql.Open("mysql", cfg.DB)
	checkErr(err)

	defer db.Close()

	stmt, err := db.Prepare("UPDATE categories SET title=?, description=? WHERE categorie_id=?")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(p.Title,p.Description,p.Category_ID)
	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
	fmt.Println(res)
	return err
}

/* AddCategory saves the values of a new category to the database and is a method to Category.
 * Called by NewCategory
 * Connect to the DB and prepares query.
 * Execute query with the inserted values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (p *Category) AddCategory() error {
	db, err := sql.Open("mysql", cfg.DB)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO categories (title,description) VALUES(?,?) ")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(p.Title,p.Description)
	affect, err := res.RowsAffected()
	fmt.Println(affect)
	fmt.Println(res)
	checkErr(err)
	return err
}
// End category methods

/* EditCategory takes updated form values from the http.request to populate a Category and call the saveCategory method.
 * The request delivers the FormValues if asked.
 * Convert category_id to an INT. The category ID is pulled from the from as a string.
 * FormValues are appointed to to the memory address of the Category struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call saveCategory, a method of the Category Struct, to update the DB
*/
func EditCategory(w http.ResponseWriter, r *http.Request,id string) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	id_string := r.FormValue("category_id")
	category_id,error := strconv.Atoi(id_string)
	checkErr(error)
	p := &Category{Category_ID: category_id, Title: title,Description: description}
	fmt.Println(p)
	err := p.saveCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/categories/"+id+"/"+title, http.StatusFound)
}

/* NewCategory takes updated form values from the http.request to populate a Category and call the AddCategory method.
 * The request delivers the FormValues if asked.
 * FormValues are appointed to to the memory address of the Category struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call AddCategory, a method of the Category Struct, to insert new category in the DB.
*/
func NewCategory(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")

	p := &Category{Title: title ,Description: description}
	fmt.Println(p)
	err := p.AddCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/categories/", http.StatusFound)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

