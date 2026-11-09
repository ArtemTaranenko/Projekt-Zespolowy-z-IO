# Weryfikacja modułu zdań

Baza implementacji: gałąź main, commit ee17948ce08eaa65f06fac55bed65ddd4e0b4579.

Wykonane sprawdzenia:

- Go 1.27.3 na Linux: go test ./... — wszystkie testy zakończone powodzeniem.
- go vet ./... — bez zgłoszonych problemów.
- go build — poprawna kompilacja programu.
- node --check frontend/sentences.js — poprawna składnia JavaScript.
- Testy zdarzeń DOM w jsdom: wybór, usuwanie i reset słów, powtórzone słowa,
  poprawna i błędna odpowiedź, blokada podwójnego wysłania, odtworzenie postępu,
  ponowienie po błędzie serwera, pusta baza i niedostępny localStorage — zaliczone.
  W tych testach odpowiedzi API były symulowane; rzeczywiste API testowano w Go.
- Migracja na kopii bazy zespołu: utworzono 10 zdań, zachowano 785 słów,
  80 tagów i 1662 powiązania. Tabela quizzes pozostała pusta, tak jak w oryginale.
- git diff --check — bez błędów białych znaków.

Nie zakończono wizualnego testu w prawdziwej przeglądarce: środowisko blokowało
uruchomienie lokalnego Chromium, a przeglądarka zewnętrzna nie otworzyła serwera.
Układ zawiera style responsywne, ale jego wygląd na komputerze i telefonie
należy sprawdzić po uruchomieniu aplikacji. Nie wykonywano kompilacji na Windows.

Przed prezentacją uruchom go test ./... i go run . -seed-sentences, otwórz
/sentences, ułóż jedną poprawną i jedną błędną odpowiedź, sprawdź Wyczyść oraz
odśwież stronę, aby sprawdzić zachowanie postępu. Zwęź okno do szerokości telefonu.
