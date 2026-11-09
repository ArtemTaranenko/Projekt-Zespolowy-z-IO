# Moduł układania zdań

Etap z karty projektu: „Implementacja modułu układania zdań”, 2.11.2026.
Użytkownik otrzymuje polskie zdanie, wybiera kolejność angielskich słów,
a serwer porównuje odpowiedź ze wzorcem zapisanym w SQLite.

## Uruchomienie na Windows w PowerShell

W katalogu z go.mod, po zainstalowaniu Go i GCC z MSYS2 UCRT64:

```powershell
$env:Path = "C:\msys64\ucrt64\bin;$env:Path"
$env:CGO_ENABLED = "1"
$env:CC = "gcc"
go test ./...
go run . -seed-sentences
```

Komendy wykonuj kolejno. Jeśli testy zgłoszą błąd, wyjaśnij go przed uruchomieniem.
Otwórz http://127.0.0.1:8080/sentences albo link na stronie głównej.
Zatrzymaj serwer przez Ctrl+C. Następnym razem wystarczy `go run .`
w terminalu z ustawionymi powyższymi zmiennymi.

Flaga -seed-sentences dodaje 10 przykładów wyłącznie do pustej tabeli sentences.
Nie modyfikuje słówek, tagów ani quizów. Powtórzenie flagi nie duplikuje przykładów.
Dotychczasowa flaga -seed nadal dodaje tylko przykładowe quizy.
Sama tabela sentences jest tworzona podczas inicjalizacji bazy.

## Obsługa

1. Przeczytaj polskie zdanie.
2. Klikaj angielskie słowa w wybranej kolejności.
3. Kliknięcie wybranego słowa odkłada je do puli; przycisk Wyczyść odkłada wszystkie.
4. Po użyciu wszystkich słów kliknij Sprawdź odpowiedź.
5. Zobacz wynik i wzorcowe zdanie. Kliknij Następne zdanie, aby kontynuować.

Powtórzone słowa są osobnymi kafelkami. Obsługiwane są kliknięcia, dotyk
oraz klawiatura (Tab i Enter/Spacja). Po sprawdzeniu nie można ponownie
naliczyć tego samego zadania. Losowanie może ponownie wybrać ten sam rekord;
nowe zadanie jest wtedy nową próbą.

Sprawdzenie dotyczy jednego wzorcowego szyku. Nie jest to analiza wszystkich
możliwych gramatycznie poprawnych tłumaczeń. Wielkość liter i białe znaki
na brzegach tokenów są ignorowane. Końcowe znaki . ! ? wzorca nie są osobnymi
kafelkami. Znaków interpunkcyjnych nie należy dodawać do words w żądaniu.

## API

GET /api/sentences/random zwraca id, polskie zdanie i przemieszane słowa:

```json
{"id":1,"polish":"Codziennie piję wodę.","words":["water","I","day","drink","every"]}
```

Odpowiedź GET nie zawiera wzorcowej kolejności ani pola correct_answer.
Przy pustej tabeli zwracany jest HTTP 404 z kodem SENTENCE_NOT_FOUND.

POST /api/sentences/check, Content-Type: application/json:

```json
{"sentence_id":1,"words":["I","drink","water","every","day"]}
```

Odpowiedź:

```json
{"correct":true,"correct_answer":"I drink water every day."}
```

Niepoprawny szyk jest normalnym wynikiem HTTP 200 z correct=false.
Pusta lista, nieprawidłowy identyfikator, puste tokeny, token zawierający kilka
słów lub niepoprawny JSON dają 400. Nieistniejący rekord daje 404.
Limit treści żądania: 8 KiB (413 po przekroczeniu); maksymalnie 64 tokeny
po 128 bajtów. Niepoprawny Content-Type: 415. Błąd bazy: 500.
Błędy mają format `{"error":{"code":"...","message":"..."}}`.

Przykłady identyfikatorów odnoszą się do nowo utworzonej tabeli. W istniejącej
bazie zawsze używaj id zwróconego przez endpoint losujący.

## Podział kodu

- db/sentences.go: tabela, model, SELECT i przykłady do inicjalizacji.
- sentence_handlers.go: losowanie kolejności, JSON, walidacja i sprawdzenie odpowiedzi.
- router.go: trasy strony, zasobów i API.
- frontend/sentences.html, sentences.css, sentences.js: interfejs ćwiczenia.
- sentence_test.go: testy nowego modułu.
- api_test.go: zaktualizowany test dawnej zaślepki, teraz realnego endpointu.

Postęp zapisuje się w localStorage pod kluczem
english-learning.sentences.progress.v1. To liczniki sprawdzonych zadań
i poprawnych odpowiedzi w konkretnej przeglądarce i dla konkretnego adresu
(origin). Zmiana 127.0.0.1 na localhost daje inny magazyn postępu.
Zablokowany zapis nie blokuje ćwiczenia, a strona pokazuje komunikat.

## Kryteria akceptacji

- Losowe polskie zdanie pochodzi z SQLite.
- Angielskie słowa są przemieszane, każde wystąpienie słowa zachowane.
- Użytkownik może złożyć i poprawić odpowiedź przed jej sprawdzeniem.
- Poprawność sprawdza backend na podstawie rekordu z bazy.
- Interfejs wyświetla wynik i wzorcową odpowiedź oraz udostępnia następne zadanie.
- Brak danych i błędy żądań mają komunikaty.
- Dotychczasowe endpointy słownika i quizu działają nadal.

## Baza w Git

Nie dołączaj zmienionego db/data.db do commita tego modułu.
Kod tworzy tabelę, a każdy członek zespołu może wykonać -seed-sentences lokalnie.
Przykładowe rekordy znajdują się w kodzie, dzięki czemu nie trzeba wymieniać
plików binarnych bazy. Istniejąca baza w repozytorium nie została zastąpiona.
