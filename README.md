# AdaGui with Go

Waehrend das Package [adatft](https://github.com/stefan-muehlebach/adatft)
ein Low-Level Interface zum Bildschirm und zum Touchscreen des TFT-Displays
von AdaFruit enthalten, versteht sich dieses Package als Sammlung von
High-Level Typen und Funktionen zur Erstellung von graphischen Oberflaechen.

Dieses Package enthaelt Datentypen, welche die Verbindung zur Hardware weiter
abstrahieren und viele komplizierte Verarbeitungen uebernehmen. Der Typ
[Screen] zum Beispiel, vereinigt in sich alles, was fuer die Kommunikation
mit der Hardware noch notwendig ist. Ihm untergeordnet ist der Typ [Window]
die Repraesentation einer Sammlung von GUI-Objekten, welche eine konkrete
Bildschirmseite konstituieren.



