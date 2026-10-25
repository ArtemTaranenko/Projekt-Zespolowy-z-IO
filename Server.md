# Serwer lokalny — instrukcja uruchamiania

Instrukcja dla projektu aplikacji do nauki języka angielskiego.
Środowisko: Windows, Git Bash, Go i GCC z MSYS2 UCRT64.

## Co zostało uruchomione

Backend Go jest programem działającym na Twoim komputerze. Po uruchomieniu
nasłuchuje pod adresem http://127.0.0.1:8080 i obsługuje żądania HTTP.
Gin kieruje żądania do odpowiednich funkcji, a SQLite przechowuje dane
w pliku db/data.db. SQLite nie wymaga osobnego serwera ani usługi.

Serwer działa, dopóki działa uruchomiony proces. Po Ctrl+C, zamknięciu
terminala lub ponownym uruchomieniu komputera trzeba uruchomić go ponownie.
Zamknięcie samej karty przeglądarki nie zatrzymuje serwera.

Nie jest to wdrożenie w internecie. Adres 127.0.0.1 wskazuje komputer,
na którym otwierasz przeglądarkę. Inni członkowie zespołu uruchamiają własną
kopię projektu. Sam push na GitHub nie uruchamia serwera.

## Wymagania

- Go zgodny z go.mod projektu (obecnie plik deklaruje 1.27.2).
- GCC dla Windows; w tej instrukcji MSYS2 UCRT64 w C:\msys64.
- CGO_ENABLED=1, ponieważ sterownik go-sqlite3 korzysta z kodu C.

Jeśli środowisko jest już skonfigurowane i testy przechodzą, nie instaluj
ponownie tych narzędzi. Przy pierwszym uruchomieniu Go pobiera zależności;
kolejne uruchomienia korzystają z lokalnej pamięci podręcznej.

Jeśli GCC jeszcze nie jest zainstalowany, zainstaluj MSYS2 z
https://www.msys2.org/ i w terminalu **MSYS2 UCRT64** wykonaj:

```bash
pacman -S mingw-w64-ucrt-x86_64-gcc
```

Dalsze polecenia wykonuj w **Git Bash**.

## Zwykłe uruchomienie

1. Otwórz Git Bash w folderze projektu, w którym znajduje się go.mod.
2. W nowym oknie terminala ustaw środowisko:

```bash
export PATH="/c/msys64/ucrt64/bin:$PATH"
export CGO_ENABLED=1
export CC=gcc
```

Jeżeli MSYS2 jest zainstalowany gdzie indziej, dostosuj ścieżkę.
Polecenia export dotyczą bieżącego okna terminala.

3. Uruchom backend:

```bash
go run .
```

Nie używaj `go run main.go`: pakiet main zawiera kilka plików.
Po komunikacie `API: http://127.0.0.1:8080/api/health` pozostaw terminal otwarty.
Brak nowego znaku `$` jest normalny: terminal jest zajęty działającym serwerem.

4. Otwórz w przeglądarce:

- http://127.0.0.1:8080/api/health
- http://127.0.0.1:8080/api/quiz/random

Aby wykonywać inne polecenia podczas działania serwera, otwórz drugi terminal.

## Pierwsze uruchomienie i przykładowe dane

Jeżeli tabela quizzes jest pusta, można dodać trzy przykładowe pytania:

```bash
go run . -seed
```

Ta komenda dodaje dane i uruchamia serwer. Nie uruchamiaj jej równolegle
z drugim serwerem na tym samym porcie.

Flaga -seed dodaje przykłady wyłącznie do pustej tabeli. Jeśli rekordy już
istnieją, nie dodaje kolejnych i nie nadpisuje danych. Na co dzień wystarczy
`go run .`. Dane pozostają w db/data.db także po wyłączeniu serwera.
Nie usuwaj ani nie zastępuj tego pliku, jeśli chcesz zachować lokalne dane.

## Zatrzymanie i zmiany w kodzie

Aby zatrzymać serwer, w jego terminalu naciśnij **Ctrl+C**.
Aby uruchomić ponownie, wpisz `go run .`.

Po zmianie plików Go zatrzymaj i uruchom backend ponownie. Projekt nie ma
mechanizmu automatycznego restartu. Samo odświeżenie strony nie przeładuje kodu Go.
Po zmianie plików HTML, CSS lub JS zwykle wystarczy odświeżenie przeglądarki.

## Testy

Po zmianach w kodzie wykonaj:

```bash
go test ./...
```

Możesz połączyć testy z uruchomieniem:

```bash
go test ./... && go run .
```

Operator && uruchomi serwer tylko wtedy, gdy testy zakończą się powodzeniem.
Testów nie trzeba wykonywać przy każdym zwykłym uruchomieniu.
Korzystają z tymczasowych baz i nie zmieniają db/data.db.

- `ok english-learning-app` oznacza zaliczone testy.
- `db [no test files]` oznacza brak osobnych plików testowych w pakiecie db;
  baza jest sprawdzana przez testy głównego pakietu.

## Dostępne endpointy

| Metoda | Ścieżka | Oczekiwany rezultat |
| --- | --- | --- |
| GET | /api/health | 200 i JSON ze statusem połączenia z bazą |
| GET | /api/quiz/random | 200 i losowy rekord; 404 przy pustej tabeli |
| GET | /api/sentences/random | 501 — moduł jeszcze niezaimplementowany |
| GET | /api/flashcards/random | 501 — moduł jeszcze niezaimplementowany |
| POST | /api/flashcards | 501 — moduł jeszcze niezaimplementowany |

Odpowiedź health:

```json
{"database":"ok","status":"ok"}
```

Quiz zwraca id, question, correct_answer, wrong1, wrong2 i wrong3.
Powtórzenie pytania po odświeżeniu strony jest poprawnym wynikiem losowania.
Mieszanie odpowiedzi i ocena odpowiedzi użytkownika należą do późniejszego etapu.

Widok JSON jest prawidłowy: endpoint API zwraca dane. Frontend będzie je później
wyświetlał jako ćwiczenia. W obecnej wersji pliki frontendu są puste,
więc pod samym adresem http://127.0.0.1:8080 nie ma jeszcze gotowej strony.

## Typowe komunikaty

| Komunikat | Znaczenie i działanie |
| --- | --- |
| `CGO_ENABLED=0` | Ustaw `export CGO_ENABLED=1` w terminalu uruchamiającym Go. |
| `gcc: command not found` | Sprawdź instalację GCC i ustaw PATH zgodnie z instrukcją. |
| `[GIN-debug] [WARNING]` | Komunikaty trybu deweloperskiego; same w sobie nie oznaczają awarii. |
| `404 /favicon.ico` | Brak ikonki karty przeglądarki; API nadal działa. |
| `QUIZ_NOT_FOUND` | Tabela jest pusta. Zatrzymaj serwer i uruchom `go run . -seed`. |
| Port 8080 zajęty / `Only one usage...` | Sprawdź, czy w innym terminalu nie działa już serwer; zatrzymaj własny poprzedni proces. |
| Przeglądarka nie może się połączyć | Sprawdź, czy serwer działa, czy nie zgłosił błędu i czy używasz adresu z portem 8080. |

## Opcjonalnie — ustawienia dla kolejnych okien Git Bash

Aby nie wpisywać trzech poleceń export w każdym nowym terminalu, otwórz
plik ~/.bashrc:

```bash
touch ~/.bashrc
notepad.exe ~/.bashrc
```

Dopisz poniższe linie jeden raz, zachowując dotychczasową zawartość:

```bash
export PATH="/c/msys64/ucrt64/bin:$PATH"
export CGO_ENABLED=1
export CC=gcc
```

Zapisz plik, a w Git Bash wykonaj:

```bash
source ~/.bashrc
```

W nowo otwartym Git Bash sprawdź `gcc --version` oraz `go env CGO_ENABLED`.
Jeżeli ustawienia nie są wczytywane w sesji logowania, sprawdź swój
~/.bash_profile i dodaj do niego poniższą linię, jeśli jeszcze jej nie ma:

```bash
[ -f ~/.bashrc ] && source ~/.bashrc
```

Te ustawienia dotyczą Twojego konta i sesji Bash, także przy innych projektach.
Nie włączają autostartu serwera i nie konfigurują PowerShell.

## Dodanie instrukcji do repozytorium

Umieść Server.md obok go.mod. W swojej gałęzi zad4 możesz zapisać samą instrukcję:

```bash
git add Server.md
git commit -m "Add local server instructions"
```

Ten commit obejmie również inne pliki, jeśli były wcześniej dodane do stagingu;
przed zatwierdzeniem sprawdź `git status`.
Zmiany w lokalnej bazie po -seed nie muszą trafiać do repozytorium.
