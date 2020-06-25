// File: main.go
package main

import (
    "encoding/json"
    "fmt"
    "log"
	"net/http"
	"io/ioutil"
)

// ConsumoDia struct to daily intake come to POST
type ConsumoDia struct {
	Data string
	Km float64
	Litros float64
}

// Consumo struct to daily intake send by JSON
type Consumo struct {
	Dia             string
	Km              float64
	Litros          float64
	ConsumoPorLitro float64
}

// ListaConsumo struct to array daily intake send by JSON
type ListaConsumo struct {
	Consumos []Consumo
}

func newConsumption(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Declare a new Person struct.
		var consumo ConsumoDia
		var erroPost string

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&consumo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fQuilometragem := consumo.Km
		fVolumeAbastecimento := consumo.Litros

		/** LÊ O JSON COM OS DADOS ATUAIS */
		listaInicial := lerJSON("consumo.json")

		kmAnterior := 0.0
		fConsumoLitros := 0.0

		/** PEGA KM ANTERIOR PARA SUBTRAÇÃO **/
		if len(listaInicial.Consumos) > 0 {
			kmAnterior = listaInicial.Consumos[len(listaInicial.Consumos)-1].Km

			kmDiff := fQuilometragem - kmAnterior
			if kmDiff <= 0 {
				erroPost = "Você informou uma quilometragem errada"
			}

			fConsumoLitros = calculaKmPorLitro(kmDiff, fVolumeAbastecimento)
		} else {
			/** PRIMEIRO INPUT NÃO CALCULA CONSUMO */
			fConsumoLitros = 0.0
		}

		fmt.Println((erroPost))
		if len(erroPost) == 0 {
			/** GERA UM NOVO ITEM PARA O JSON */
			newItem := Consumo{
				Dia:             consumo.Data,
				Km:              fQuilometragem,
				Litros:          fVolumeAbastecimento,
				ConsumoPorLitro: fConsumoLitros,
			}

			/** ADICIONA O NOVO ITEM A LISTA PARA GERAR O NOVO JSON */
			listaInicial.AddItem(newItem)

			/** GERA O NOVO JSON */
			gerarJSON(listaInicial)
		}

		b, err := json.Marshal(listaInicial)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		fmt.Fprintf(w, string(b))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
    	w.Write([]byte("405 - Method Not Allowed!"))
	}
}

func listConsumption (w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		/** LÊ O JSON COM OS DADOS ATUAIS */
		listaInicial := lerJSON("consumo.json")

		b, err := json.Marshal(listaInicial)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		fmt.Fprintf(w, string(b))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
    	w.Write([]byte("405 - Method Not Allowed!"))
	}
}

func main() {
    mux := http.NewServeMux()
	mux.HandleFunc("/consumo/add", newConsumption)
	mux.HandleFunc("/consumo", listConsumption)

    err := http.ListenAndServe(":8080", mux)
    log.Fatal(err)
}

func calculaKmPorLitro(km float64, lt float64) (kmProLitro float64) {
	kmProLitro = km / lt
	return
}

func gerarJSON(lista ListaConsumo) bool {

	file, _ := json.MarshalIndent(lista, "", " ")
	_ = ioutil.WriteFile("consumo.json", file, 0644)

	return true
}

func lerJSON(filename string) ListaConsumo {

	plan, _ := ioutil.ReadFile(filename)
	var data ListaConsumo
	_ = json.Unmarshal(plan, &data)

	return data
}

// AddItem adding new item on ListaConsumo
func (lista *ListaConsumo) AddItem(item Consumo) []Consumo {
	lista.Consumos = append(lista.Consumos, item)
	return lista.Consumos
}
