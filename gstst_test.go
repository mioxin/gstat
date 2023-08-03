package main

import "testing"

const OUTPUT string = `{"bin":"860161","name":"ИП \"АБДЕВА\"","registerDate":"","okedCode":"47261","okedName":"Розничная торговля табачными изделиями в специализированных магазинах, являющихся торговыми объектами, с торговой площадью менее 2000 кв.м","krpCode":"105","krpName":"Малые предприятия (\u003c= 5)","krpBfCode":"105","krpBfName":"Малые предприятия (\u003c= 5)","kseCode":"1122","kseName":"Национальные частные нефинансовые корпорации – ОПП","katoAddress":"Г.АЛМАТЫ, АЛМАЛИНСКИЙ РАЙОН","fio":"АБДЕВА АРК АХНА","ip":true}
{"bin":"8905548","name":"ИП \"Аутсорсинговая компания Аксултан\"","registerDate":"","okedCode":"69202","okedName":"Деятельность в области составления счетов и бухгалтерского учета","krpCode":"105","krpName":"Малые предприятия (\u003c= 5)","krpBfCode":"","krpBfName":"","kseCode":"1122","kseName":"Национальные частные нефинансовые корпорации – ОПП","katoAddress":"ЖАМБЫЛСКАЯ ОБЛАСТЬ, ТАРАЗ Г.А., Г.ТАРАЗ","fio":"АБВА БАН РЕВНА","ip":true}
{"bin":"90032396","name":"ИП \"ГУЛБАНУ\"","registerDate":"","okedCode":"96090","okedName":"Предоставление прочих индивидуальных услуг, не включенных в другие группировки","krpCode":"105","krpName":"Малые предприятия (\u003c= 5)","krpBfCode":"","krpBfName":"","kseCode":"1122","kseName":"Национальные частные нефинансовые корпорации – ОПП","katoAddress":"АЛМАТИНСКАЯ ОБЛАСТЬ, КАРАСАЙСКИЙ РАЙОН, КАСКЕЛЕНСКАЯ Г.А., Г.КАСКЕЛЕН","fio":"АБЕВА ГУНУ ОВНА","ip":true}
`
const INPUT string = `000045570;000045570;Абаев Рль Риич;25.09.1996;960117
000030026;000030026;Абев Тайс Нурвич;07.07.1994;94079028
000034944;000034944;Абева Айа Мур;03.02.1998;98020348
000003345;000003345;Абева Ай Абдуы;  .  .    ;
000009245;000009245;Абева Ай Абдуы;  .  .    ;
000025606;000025606;Абева Айа Базна;27.07.1989;89072635
000019695;000019695;Абева Аке Дауна;08.01.1990;90010903
000014591;000014591;Абева Бооз Сеа;27.06.1990;90062728
000023812;000023812;Абева Га Балувна;06.07.1989;89071908
000050319;000050319;Абева Га Сайлвна;06.09.1988;88090340
000020224;000020224;Абева Гуну Мана;23.03.1990;90032396
000023916;000023916;Абева Гуза Осовна ;05.01.1991;91400572
000037681;000037681;Абева Гуаз Аха;26.01.1975;75012606
000040506;000040506;Абева Гуаз Гана;13.11.1995;95111407
`

func TestSkipLines(t *testing.T) {

}