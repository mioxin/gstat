# gstat
__gstat__ receive an individual business info by the Individual Identification Number (IIN) from "http://old.stat.gov.kz/api". If it registred.

__gstat__ can work in the multithread mode. It use connection pool. The number of concurency connections define -t option. By default the log writen to the file "gstat.log". If defined option -v log otput to the StdOut.

__gstat__ continue work from the last received position after cancel work by emergency.

If request rejected by timeout or statuscode 429 (Too many requests) before next request will be add a wait time 3 sec and each next wait will be increase on 3 sec will repeat 20 once.

Using:
---
__gstat <-h> <-v> <-t N> <-o output_file> <-i input_file>__

#### Flags:
    -h, -help:  Show help (Default: false)
    -i:         The input file (Default: in.txt)
    -o:         The output file (Default: out.json)
    -t:         The number of threads (Default: 1)
    -v, -verbouse:  Output log to StdOut (Default: false)

Any option have default value.\
The input file format is CSV with semicolon separated. The IIN should in 5th position.
#### The output file is json include data:

	Bin             string `json:"bin"`
	Name            string `json:"name"`         
	RegisterDate    string `json:"registerDate"` 
	OkedCode        string `json:"okedCode"`     
	OkedName        string `json:"okedName"`     
	KrpCode         string `json:"krpCode"`   
	KrpName         string `json:"krpName"`   
	KrpBfCode       string `json:"krpBfCode"` 
	KrpBfName       string `json:"krpBfName"` 
	KseCode         string `json:"kseCode"`   
	KseName         string `json:"kseName"`   
	KatoAddress     string `json:"katoAddress"` 
	Fio             string `json:"fio"`         
	Ip              bool   `json:"ip"`

