
const tagSelect = document.getElementById("tagSelect");
const randomButton = document.getElementById("randomButton");
const message = document.getElementById("message");
const wordResult = document.getElementById("wordResult");

const polishWord = document.getElementById("polishWord");
const englishWord = document.getElementById("englishWord");
const wordTags = document.getElementById("wordTags");

// Pobieranie tagów z backendu
async function loadTags() {
    try {
        const response = await fetch("/api/tags");

        if (!response.ok) {
            throw new Error("Nie udało się pobrać tagów.");
        }

        const tags = await response.json();

        tagSelect.replaceChildren();

        const placeholder = document.createElement("option");
        placeholder.value = "";
        placeholder.textContent = "-- Wybierz tag --";
        tagSelect.appendChild(placeholder);

        tags.forEach((tag) => {
            const option = document.createElement("option");
            option.value = tag.id;
            option.textContent = tag.name;
            tagSelect.appendChild(option);
        });

        message.textContent = tags.length
            ? ""
            : "Brak tagów w bazie danych.";
    } catch (error) {
        message.textContent = error.message;
    }
}

// Losowanie słówka
async function randomWord() {
    const tagID = tagSelect.value;

    if (!tagID) {
        message.textContent = "Najpierw wybierz tag.";
        wordResult.hidden = true;
        return;
    }

    randomButton.disabled = true;
    message.textContent = "Losowanie...";
    wordResult.hidden = true;

    try {
        const response = await fetch(
            `/api/words/random?tag_id=${encodeURIComponent(tagID)}`
        );

        const word = await response.json();

        if (!response.ok) {
            throw new Error(word.error || "Wystąpił błąd.");
        }

        polishWord.textContent = word.polish;
        englishWord.textContent = word.english;
        wordTags.textContent = word.tags.join(", ");

        message.textContent = "";
        wordResult.hidden = false;
    } catch (error) {
        message.textContent = error.message;
    } finally {
        randomButton.disabled = false;
    }
}

randomButton.addEventListener("click", randomWord);

loadTags();
