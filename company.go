package main

import (
	"fmt"
)

const URL string = "https://old.stat.gov.kz/api/juridical/counter/api/" //?bin=840629300619&lang=ru"


type Company struct {
	Bin          string `json:"bin"`          //"840629300619",
	Name         string `json:"name"`         // "ИП КОВАЛЕВ ИВАН АЛЕКСАНДРОВИЧ",
	RegisterDate string `json:"registerDate"` // null,
	OkedCode     string `json:"okedCode"`     // "62099",
	OkedName     string `json:"okedName"`     // "Другие виды деятельности в области информационных технологий и информационных систем, не включенные в другие группировки",
	//SecondOkeds  string `json:"secondOkeds"`  // null,
	KrpCode   string `json:"krpCode"`   // "105",
	KrpName   string `json:"krpName"`   // "Малые предприятия (<= 5)",
	KrpBfCode string `json:"krpBfCode"` // "105",
	KrpBfName string `json:"krpBfName"` // "Малые предприятия (<= 5)",
	KseCode   string `json:"kseCode"`   // "1122",
	KseName   string `json:"kseName"`   // "Национальные частные нефинансовые корпорации – ОПП",
	//KfsCode      string `json:"kfsCode"`      // null,
	//KfsName      string `json:"kfsName"`      // null,
	//KatoCode     string `json:"katoCode"`     // "631010000",
	//KatoId       string `json:"katoId"`       // 264992,
	KatoAddress string `json:"katoAddress"` // "ВОСТОЧНО-КАЗАХСТАНСКАЯ ОБЛАСТЬ, УСТЬ-КАМЕНОГОРСК Г.А., Г.УСТЬ-КАМЕНОГОРСК",
	Fio         string `json:"fio"`         // "КОВАЛЕВ ИВАН АЛЕКСАНДРОВИЧ",
	Ip          bool   `json:"ip"`          // "true
}


func (c *Company) String() string {
	return fmt.Sprintf(`
БИН:			%v
Наименование:	%v
ОКЭД:			%v
Описание ОКЭД:	%v
Адрес:			%v
ФИО:			%v
Размерность:	%v
ИП:				%v
`, c.Bin, c.Name, c.OkedCode, c.OkedName, c.KatoAddress, c.KrpName, c.Fio, c.Ip)
}
