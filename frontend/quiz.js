"use strict";

const elements = Object.fromEntries([
    "status",
    "exercise-body",
    "question",
    "answers",
    "check",
    "next",
    "feedback",
    "feedback-title",
    "correct-answer",
    "checked-count",
    "correct-count",
    "storage-note"
].map(id => [id, document.getElementById(id)]));

const progressKey = "english-learning.quiz.progress.v1";

let quiz = null;
let selectedAnswer = null;
let busy = false;
let checked = false;
let progress = { checked: 0, correct: 0 };

try {
    const stored = JSON.parse(localStorage.getItem(progressKey));
    if (
        stored &&
        Number.isSafeInteger(stored.checked) &&
        Number.isSafeInteger(stored.correct) &&
        stored.checked >= stored.correct &&
        stored.correct >= 0
    ) {
        progress = stored;
    }
} catch (_) {
    elements["storage-note"].hidden = false;
}

function renderProgress() {
    elements["checked-count"].textContent = progress.checked;
    elements["correct-count"].textContent = progress.correct;
}

function renderAnswers() {
    elements["answers"].replaceChildren();

    if (!quiz) return;

    quiz.answers.forEach(answer => {
        const button = document.createElement("button");
        button.type = "button";
        button.className = "answer";
        button.textContent = answer;
        button.disabled = busy || checked;

        if (answer === selectedAnswer) {
            button.classList.add("selected");
        }

        button.addEventListener("click", () => {
            if (busy || checked) return;

            selectedAnswer = answer;
            renderAnswers();
            elements["check"].disabled = false;
        });

        elements["answers"].appendChild(button);
    });
}

function render() {
    renderAnswers();
    elements["check"].disabled = busy || checked || !quiz || !selectedAnswer;
    elements["next"].disabled = busy;
}

async function request(url, options = {}) {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), 10000);

    try {
        const response = await fetch(url, {
            ...options,
            signal: controller.signal,
            cache: "no-store"
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error?.message || "Nie udało się wykonać żądania.");
        }

        return data;
    } catch (error) {
        if (error.name === "AbortError") {
            throw new Error("Serwer nie odpowiada. Spróbuj ponownie.");
        }

        throw error;
    } finally {
        clearTimeout(timer);
    }
}

async function loadQuiz() {
    if (busy) return;

    busy = true;
    checked = false;
    quiz = null;
    selectedAnswer = null;

    elements["exercise-body"].hidden = true;
    elements["feedback"].hidden = true;
    elements["status"].textContent = "Ładowanie pytania…";
    render();

    try {
        quiz = await request("/api/quiz/random");
        elements["question"].textContent = quiz.question;
        elements["exercise-body"].hidden = false;
        elements["status"].textContent = "";
    } catch (error) {
        elements["status"].textContent =
            error.message || "Nie udało się połączyć z serwerem.";
    } finally {
        busy = false;
        render();
    }
}

async function checkQuiz() {
    if (busy || checked || !quiz || !selectedAnswer) return;

    busy = true;
    elements["status"].textContent = "Sprawdzanie…";
    render();

    try {
        const result = await request("/api/quiz/check", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                quiz_id: quiz.id,
                answer: selectedAnswer
            })
        });

        checked = true;
        progress.checked += 1;

        if (result.correct) {
            progress.correct += 1;
        }

        elements["feedback"].className = result.correct
            ? "success"
            : "incorrect-feedback";

        elements["feedback-title"].textContent = result.correct
            ? "Dobrze! Poprawna odpowiedź."
            : "Nie tym razem. Poprawna odpowiedź:";

        elements["correct-answer"].textContent = result.correct_answer;
        elements["feedback"].hidden = false;
        elements["status"].textContent = "";

        try {
            localStorage.setItem(progressKey, JSON.stringify(progress));
        } catch (_) {
            elements["storage-note"].hidden = false;
        }

        renderProgress();
    } catch (error) {
        elements["status"].textContent =
            error.message || "Nie udało się połączyć z serwerem.";
    } finally {
        busy = false;
        render();

        if (checked) {
            elements["next"].focus();
        }
    }
}

elements["check"].addEventListener("click", checkQuiz);
elements["next"].addEventListener("click", loadQuiz);

renderProgress();
loadQuiz();