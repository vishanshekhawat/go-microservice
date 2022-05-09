package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/vishn001/build-microservice/data"
)

type Product struct {
	logger *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{
		logger: l,
	}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {

		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			p.logger.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.logger.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.logger.Println("Invalid URI unable to convert to numer", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, rw, r)
		return

	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Post Method")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal Json", http.StatusInternalServerError)
	}

	p.logger.Printf("Prod : %#v", prod)
	data.AddProduct(prod)
}
func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal Json", http.StatusInternalServerError)
	}
}
func (p *Product) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
