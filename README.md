# Aplikacja do nauki języka angielskiego

Projekt zespołowy realizowany na Politechnice Częstochowskiej.
Aplikacja ma umożliwiać naukę przez układanie zdań, fiszki i quizy.
Backend: Go i Gin. Baza danych: SQLite. Frontend: HTML, CSS i JavaScript.
Aplikacja nie wymaga konta użytkownika.

## Aktualny zakres

| Moduł | Stan |
| --- | --- |
| Słownik i tagi | Baza słów, lista tagów, losowanie słowa po tagu i testowa strona słownika |
| Układanie zdań | Strona ćwiczenia, przemieszane słowa, sprawdzanie odpowiedzi przez API i lokalny postęp |
| Quiz | Odczyt losowego pytania z SQLite; pełny interfejs quizu pozostaje do wykonania |
| Fiszki | Podstawa danych dostępna przez API słownika; pełna obsługa odwracania fiszek pozostaje do wykonania |
| Własne słówka | Dodawanie pozostaje do wykonania |

## Uruchomienie

Wymagane są Go zgodny z go.mod oraz GCC i włączone CGO dla sterownika SQLite.
W Windows PowerShell, w katalogu z go.mod:

```powershell
$env:Path = "C:\msys64\ucrt64\bin;$env:Path"
$env:CGO_ENABLED = "1"
$env:CC = "gcc"
go test ./...
go run . -seed-sentences
```

Ścieżka GCC zakłada instalację MSYS2 UCRT64 w C:\msys64.
Po pomyślnych testach uruchom serwer. Adresy:

- http://127.0.0.1:8080 — testowy słownik;
- http://127.0.0.1:8080/sentences — układanie zdań;
- http://127.0.0.1:8080/api/health — stan backendu i połączenia z bazą.

Po utworzeniu przykładów kolejne uruchomienia: `go run .`.
Zatrzymanie serwera: Ctrl+C. Nie ma automatycznego restartu po zmianie Go.
Dane przechowywane są w db/data.db. Ścieżki są względne do katalogu uruchomienia.

## Dokumentacja

- [Instrukcja modułu układania zdań](docs/Sentences.md) — API, obsługa, kryteria akceptacji.
- [Instrukcja serwera](Server.md) — konfiguracja środowiska i podstawy uruchamiania;
  opis zakresu modułów w tym starszym pliku pochodzi z etapu 4.

Testy: `go test ./...`. Analiza kodu: `go vet ./...`.
Zmiany rozwijamy w osobnych gałęziach i łączymy przez Pull Request.
