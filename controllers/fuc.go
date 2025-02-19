package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type APIFUC struct {
	token          string
	consumerKey    string
	consumerSecret string
}

type UserFUC struct {
	ID                         string      `json:"id"`
	IdentidadNumero            string      `json:"identidad_numero"`
	PrimerNombre               string      `json:"primer_nombre"`
	SegundoNombre              string      `json:"segundo_nombre"`
	PrimerApellido             string      `json:"primer_apellido"`
	SegundoApellido            string      `json:"segundo_apellido"`
	Sexo                       string      `json:"sexo"`
	Edad                       int64       `json:"edad"`
	NombrePadre                string      `json:"nombre_padre"`
	NombreMadre                string      `json:"nombre_madre"`
	Direccion                  string      `json:"direccion"`
	MunicipioResidenciaSid     int64       `json:"municipio_residencia_sid"`
	MunicipioResidencia        string      `json:"municipio_residencia"`
	MunicipioResidenciaCodDpa  string      `json:"municipio_residencia_cod_dpa"`
	MunicipioResidenciaCodSuin string      `json:"municipio_residencia_cod_suin"`
	ProvinciaResidenciaSid     int64       `json:"provincia_residencia_sid"`
	ProvinciaResidencia        string      `json:"provincia_residencia"`
	ProvinciaResidenciaCodDpa  string      `json:"provincia_residencia_cod_dpa"`
	ProvinciaResidenciaCodSuin string      `json:"provincia_residencia_cod_suin"`
	FechaActRid                time.Time   `json:"fecha_act_rid"`
	FechaActFuc                time.Time   `json:"fecha_act_fuc"`
	FechaFoto                  time.Time   `json:"fecha_foto"`
	FechaFirma                 time.Time   `json:"fecha_firma"`
	Fallecido                  bool        `json:"fallecido"`
	FechaActRec                interface{} `json:"fecha_act_rec"`
	NacimientoAnnoReg          string      `json:"nacimiento_anno_reg"`
	NacimientoFecha            string      `json:"nacimiento_fecha"`
	NacimientoTomo             string      `json:"nacimiento_tomo"`
	NacimientoFolio            string      `json:"nacimiento_folio"`
	NacimientoMunicipioSid     int64       `json:"nacimiento_municipio_sid"`
	NacimientoMunicipio        string      `json:"nacimiento_municipio"`
	NacimientoMunicipioCodDpa  string      `json:"nacimiento_municipio_cod_dpa"`
	NacimientoMunicipioCodSuin string      `json:"nacimiento_municipio_cod_suin"`
	NacimientoProvinciaSid     int64       `json:"nacimiento_provincia_sid"`
	NacimientoProvincia        string      `json:"nacimiento_provincia"`
	NacimientoProvinciaCodDpa  string      `json:"nacimiento_provincia_cod_dpa"`
	NacimientoProvinciaCodSuin string      `json:"nacimiento_provincia_cod_suin"`
	NacimientoRegistroCivilSid int64       `json:"nacimiento_registro_civil_sid"`
	NacimientoRegistroCivil    string      `json:"nacimiento_registro_civil"`
	DefuncionAnnoReg           interface{} `json:"defuncion_anno_reg"`
	DefuncionFecha             interface{} `json:"defuncion_fecha"`
	DefuncionTomo              interface{} `json:"defuncion_tomo"`
	DefuncionFolio             interface{} `json:"defuncion_folio"`
	DefuncionRegistroCivilSid  interface{} `json:"defuncion_registro_civil_sid"`
	DefuncionRegistroCivil     interface{} `json:"defuncion_registro_civil"`
	CiudadaniaSid              string      `json:"ciudadania_sid"`
	Ciudadania                 string      `json:"ciudadania"`
	CiudadaniaGentilicio       string      `json:"ciudadania_gentilicio"`
	CondicionMigratoriaSid     int64       `json:"condicion_migratoria_sid"`
	CondicionMigratoria        string      `json:"condicion_migratoria"`
}

func NewAPIFUC(consumerKey, consumerSecret string) (*APIFUC, error) {
	token, err := getToken(consumerKey, consumerSecret)
	if err != nil {
		return nil, err
	}
	return &APIFUC{
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
		token:          token,
	}, nil
}

// :return: Returns an API access token It depends on the input of the public and private keys.
func getToken(consumerKey, consumerSecret string) (string, error) {
	auth := consumerKey + ":" + consumerSecret
	bs4 := base64.StdEncoding.EncodeToString([]byte(auth))
	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("scope", "nivel10")
	data := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", "https://apis-fuc.minjus.gob.cu/token", data)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+bs4)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("ERROR: Get Token failed with HTTP status code: " + strconv.Itoa(resp.StatusCode))
	}
	var jsonResp map[string]interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Fatal(err)
	}
	return jsonResp["access_token"].(string), nil
}

func (f *APIFUC) getDataUser(IdCard string) ([]UserFUC, error) {
	// Define the parameters
	params := url.Values{}
	params.Add("identidad_numero", IdCard)
	params.Add("similar", "false")

	// Create the URL with the parameters
	apiUrl := "https://apis-fuc.minjus.gob.cu/pn-api-consulta/2.0.210131/api/v1/nivel10?" + params.Encode()

	// Realizar la solicitud HTTP GET
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con el API FUC: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+f.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Leer la respuesta del API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ERROR: No se pudo leer la respuesta del API FUC")
	}

	var response []UserFUC

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("ERROR: Formato de respuesta inválido desde el API FUC")
	}

	return response, nil
}
