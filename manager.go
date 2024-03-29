package main

import (
	"os"
	"time"
	"encoding/json"
	"bufio"
	"fmt"
	//"log"
	"os/exec"
	"strconv"
	"bytes"
	"strings"
)
// Url Define el tipo de dato del objeto que sera json (URLS)
type Url struct {
	Urls []string
}
// Geo Define el tipo de dato del objeto que sera json (Lat,Long,Elev)
type Geo struct {
	Geolocs []string
}

// containerUp crea, levanta un worker (contenedor), y le da trabajo
// j: Json tipo string, es el trabajo que sera dado al worker (urls)
// co: Json tipo string, son los datos (Lat, Lon, Elev) de cada estacion, que sera gregado a cada url descargado
func containerUp(j string, co string) {
	//vol := "C:\\Users\\JOSE\\Desktop\\Project\\pruebas\\data:/usr/src/app/data"
	//c, err := exec.Command("docker", "run", "--name", "extractor", "--rm", "-v", vol, "dataextractor:v1", "-j", j, "-c", co).Output()
	/*_, err := exec.Command("python", "worker.py", "-j", j, "-c", co).Output()
	if err != nil {
		log.Fatal(err)
	}*/
	c := exec.Command("python", "worker.py", "-j", j, "-c", co)
	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    	return
	}
	//fmt.Printf("%s", c)
}

func main() {
	mx := 267 // MX = 267
	i, _ := strconv.Atoi(os.Args[1]) // Año inicio
	f, _ := strconv.Atoi(os.Args[2]) // Año fin
	replicas := ((f - i) + 1) * 6 // N replicas a realizar
	nUrls := ((mx * (f - i + 1)) / replicas) + 1 // N de urls por replica
	var urls []string
	var geos []string

	file, err := os.Open("./idAntenas.txt") // Abrir archivo para poder escanear
	if err != nil {
		fmt.Println("Error")
	}

	scanner := bufio.NewScanner(file) // Instanciar objeto Scanner
	y := 0 // Contador, n de lineas de MX, MX = 267
	z := 1 // Contador, controlar y validar la asignacion de numero de urls por replica, n de replicas
	ti := time.Now()
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, " MX ") {
			slice := strings.Split(line, " ") // Convert str to []str
			stnCode := strings.Join(slice[:2], "-") // Obtain station code
			coord := strings.Join(slice[len(slice) - 5: len(slice) - 2], "_") // Lat, Long, Elev of station
			for x := 0; x <= (f - i); x++ {
				y++
				urls = append(urls, "pub/data/gsod/"+strconv.Itoa(i+x)+"/"+stnCode+"-"+strconv.Itoa(i+x)+".op.gz") // create a url, add to []urls
				geos = append(geos, coord) // add coord to []geos
				if y == (nUrls*z) || y == (mx*(f-i+1)) { // Si llego al numero maximo de url por replicas, o llego al final
					z++
					//fmt.Println("Replica", z-1)
					url := Url{ // Create a instancia of Url
						Urls: urls,
					}

					geo := Geo{
						Geolocs: geos,
					}
					
					jUrls, _ := json.Marshal(url) // Create a json URLS
					jGeos, _ := json.Marshal(geo) // Create a json Geos
					
					//fmt.Println(string(jUrls))
					//fmt.Println(string(jGeos))

					go containerUp(string(jUrls), string(jGeos)) // Levantamos workers (contenedores)

					urls = nil // Clear []urls to fill again with new urls
					geos = nil // Clear []geos to fill again with new geos
				}
			}
		}
	}

	tf := time.Now()
	time.Sleep(300 * time.Second);
	fmt.Println(tf.Sub(ti))
	file.Close() // Close the file

}
