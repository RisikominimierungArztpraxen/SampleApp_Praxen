# Beispiel ClientApp für Arztpraxen

Die im Hackathon erstellte App dient als Beispiellösungprototyp, der entweder Termine durch die eigene Weboberfläche (noch nicht implementiert) oder via REST-API, durch die vorhandene Terminierungssoftware in den Praxen mit der App kommunizieren kann, entgegen nimmt. Der Prototyp sendet dann API-Calls zur QueuingApp, die nur funktional notwendige Daten enthalten. Da die Beispielösungsapp in der Praxis selbst läuft, verlassen alle anderen Daten die Praxis nicht.

Zu diesem Zeitpunkt ist die Beispiellösung bereits funktionstüchtig. Um "in production" verwendet werden zu können, fehlen aber noch z.B. Error Handling, die Anbindung des konzipierten UI für das eigen Webinterface, eine Anbindung zu einer internen DB, und eine bessere Dokumentation. Da der Prototyp in Go geschrieben ist, kann er auf allen gängigen Betriebssystemen (Windows, MacOS, Linux) installiert werden.

### Beispieldatensatz der an die praxisinterne App gesendet wird
```
{
        "time": "13:30",
        "patientId": "E73",
        "patientName": "Max Mustermann",
        "patientDoB": "03.07.1963",
        "notifications": [
            {
                "type": "app",
                "identifier": "P51kyoH8zx"
            },
            {
                "type": "sms",
                "identifier": "0049159123456"
            }
        ],
        "estimatedInMinutes": 15,
        "urgent": false,
        "potentialCOVID-19": true,
        "queuingApp": true
    }
```

### Beispieldatensatz der an die externe QueuingApp gesendet wird
```
{
        "time": "13:30",
        "patientId": "E73",
        "notifications": [
            {
                "type": "app",
                "identifier": "P51kyoH8zx"
            },
            {
                "type": "sms",
                "identifier": "0049159123456"
            }
        ],
        "estimatedInMinutes": 15
    }
```