const supportedLanguage = hljs.listLanguages().filter(hljs.getLanguage);

function guessProgrammingLanguage(value) {
  const match = supportedLanguage.reduce(
    (previous, next) => {
      const result = hljs.highlight(value, {
        language: next,
        ignoreIllegals: false,
      });

      if (result.relevance > previous.relevance) {
        return { ...result, language: next };
      }

      return previous;
    },
    { relevance: 0, value }
  );

  return match.language || null;
}

function renderHighlightedCode(content) {
  const highlighted = hljs.highlightAuto(content);
  const lines = highlighted.value.split(/\r\n|\r|\n/);

  const codeElement = document.getElementById("code");
  codeElement.innerHTML = "";

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    const lineElement = document.createElement("div");
    lineElement.classList.add("view-line");

    const lineNumberElement = document.createElement("div");
    lineNumberElement.classList.add("view-line-number");
    lineNumberElement.innerText = i + 1;

    const lineCodeElement = document.createElement("code");
    lineCodeElement.classList.add("view-line-code");
    lineCodeElement.innerHTML = line;

    lineElement.appendChild(lineNumberElement);
    lineElement.appendChild(lineCodeElement);

    codeElement.appendChild(lineElement);
  }
}

document.addEventListener("htmx:configRequest", (e) => {
  const content = e.detail.parameters["content"];

  if (!content) return;

  const language = guessProgrammingLanguage(content);
  e.detail.parameters["language"] = language || "";
});
