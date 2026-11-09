"use strict";

const elements = Object.fromEntries([
  "status", "exercise-body", "polish", "answer", "bank", "word-count",
  "check", "reset", "next", "feedback", "feedback-title", "correct-answer",
  "checked-count", "correct-count", "storage-note"
].map(id => [id, document.getElementById(id)]));
const progressKey = "english-learning.sentences.progress.v1";
let exercise = null;
let selected = []; // Indeksy słów, nie ich tekst: powtórzone słowa pozostają oddzielnymi kafelkami.
let busy = false;
let checked = false;
let progress = { checked: 0, correct: 0 };

try {
  const stored = JSON.parse(localStorage.getItem(progressKey));
  if (stored && Number.isSafeInteger(stored.checked) && Number.isSafeInteger(stored.correct)
      && stored.correct >= 0 && stored.checked >= stored.correct) progress = stored;
} catch (_) {
  elements["storage-note"].hidden = false;
}

function renderProgress() {
  elements["checked-count"].textContent = progress.checked;
  elements["correct-count"].textContent = progress.correct;
}

function render() {
  elements["answer"].replaceChildren();
  elements["bank"].replaceChildren();
  const words = exercise ? exercise.words : [];
  const used = new Set(selected);
  for (const index of selected) addWord(index, elements["answer"], true);
  words.forEach((_, index) => { if (!used.has(index)) addWord(index, elements["bank"], false); });
  if (selected.length === 0) {
    const hint = document.createElement("p");
    hint.className = "placeholder";
    hint.textContent = "Tutaj pojawi się Twoje zdanie…";
    elements["answer"].appendChild(hint);
  }
  elements["word-count"].textContent = `${selected.length} / ${words.length} słów`;
  elements["check"].disabled = busy || checked || !exercise || selected.length !== words.length;
  elements["reset"].disabled = busy || checked || selected.length === 0;
  elements["next"].disabled = busy;
}

function addWord(index, parent, inAnswer) {
  const button = document.createElement("button");
  button.type = "button";
  button.className = "word";
  button.textContent = exercise.words[index];
  button.dataset.index = index;
  button.setAttribute("aria-label", `${inAnswer ? "Odłóż" : "Dodaj"} słowo: ${exercise.words[index]}`);
  button.disabled = busy || checked;
  button.addEventListener("click", () => {
    if (busy || checked) return;
    if (inAnswer) selected = selected.filter(item => item !== index);
    else selected.push(index);
    render();
    const target = inAnswer ? elements["bank"] : elements["answer"];
    target.querySelector(`[data-index="${index}"]`).focus();
  });
  parent.appendChild(button);
}

async function request(url, options = {}) {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), 10000);
  try {
    const response = await fetch(url, { ...options, signal: controller.signal, cache: "no-store" });
    const data = await response.json();
    if (!response.ok) throw new Error(data.error?.message || "Nie udało się wykonać żądania.");
    return data;
  } catch (error) {
    if (error.name === "AbortError") throw new Error("Serwer nie odpowiada. Spróbuj ponownie.");
    throw error;
  } finally {
    clearTimeout(timer);
  }
}

async function loadSentence() {
  if (busy) return;
  busy = true;
  checked = false;
  exercise = null;
  selected = [];
  elements["exercise-body"].hidden = true;
  elements["feedback"].hidden = true;
  elements["status"].textContent = "Ładowanie zdania…";
  render();
  try {
    exercise = await request("/api/sentences/random");
    elements["polish"].textContent = exercise.polish;
    elements["exercise-body"].hidden = false;
    elements["status"].textContent = "";
  } catch (error) {
    elements["status"].textContent = error.message || "Nie udało się połączyć z serwerem.";
  } finally {
    busy = false;
    render();
  }
}

async function checkSentence() {
  if (busy || checked || !exercise || selected.length !== exercise.words.length) return;
  busy = true;
  elements["status"].textContent = "Sprawdzanie…";
  render();
  try {
    const result = await request("/api/sentences/check", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ sentence_id: exercise.id, words: selected.map(index => exercise.words[index]) })
    });
    checked = true;
    elements["feedback"].className = result.correct ? "success" : "incorrect";
    elements["feedback-title"].textContent = result.correct ? "Dobrze! Zdanie jest poprawne." : "Jeszcze nie tym razem. Spójrz na wzór:";
    elements["correct-answer"].textContent = result.correct_answer;
    elements["feedback"].hidden = false;
    elements["status"].textContent = "";
    progress.checked += 1;
    if (result.correct) progress.correct += 1;
    try { localStorage.setItem(progressKey, JSON.stringify(progress)); }
    catch (_) { elements["storage-note"].hidden = false; }
    renderProgress();
  } catch (error) {
    elements["status"].textContent = error.message || "Nie udało się połączyć z serwerem.";
  } finally {
    busy = false;
    render();
    if (checked) elements["next"].focus();
  }
}

elements["next"].addEventListener("click", loadSentence);
elements["check"].addEventListener("click", checkSentence);
elements["reset"].addEventListener("click", () => { if (!busy && !checked) { selected = []; render(); } });
renderProgress();
loadSentence();
