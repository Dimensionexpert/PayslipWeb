const tabs = document.querySelectorAll(".tab");
const tabPanels = document.querySelectorAll(".tab-panel");

const dropZone = document.querySelector("#drop-zone");
const fileInput = document.querySelector("#file-input");
const fileList = document.querySelector("#file-list");
const fileCount = document.querySelector("#file-count");

const generateButton = document.querySelector("#generate-button");
const progressSection = document.querySelector("#progress-section");
const progressBar = document.querySelector("#progress-bar");
const progressLabel = document.querySelector("#progress-label");
const progressCount = document.querySelector("#progress-count");
const progressMessage = document.querySelector("#progress-message");

const searchInput = document.querySelector("#search-input");
const pdfList = document.querySelector("#pdf-list");

let selectedFiles = [];

// Tab switching
tabs.forEach((tab) => {
  tab.addEventListener("click", () => {
    const selectedTab = tab.dataset.tab;

    tabs.forEach((item) => {
      item.classList.toggle("active", item === tab);
    });

    tabPanels.forEach((panel) => {
      panel.classList.toggle("active", panel.id === selectedTab);
    });
  });
});

// File selection
fileInput.addEventListener("change", () => {
  selectedFiles = [...fileInput.files];
  renderSelectedFiles();
});

// Drag and drop
dropZone.addEventListener("dragover", (event) => {
  event.preventDefault();
  dropZone.classList.add("dragging");
});

dropZone.addEventListener("dragleave", () => {
  dropZone.classList.remove("dragging");
});

dropZone.addEventListener("drop", (event) => {
  event.preventDefault();
  dropZone.classList.remove("dragging");

  selectedFiles = [...event.dataTransfer.files];
  renderSelectedFiles();
});

function renderSelectedFiles() {
  fileCount.textContent = `${selectedFiles.length} file${
    selectedFiles.length === 1 ? "" : "s"
  }`;

  if (selectedFiles.length === 0) {
    fileList.innerHTML = `
            <p class="empty-state">
                No files selected yet.
            </p>
        `;
    return;
  }

  fileList.innerHTML = "";

  selectedFiles.forEach((file) => {
    const item = document.createElement("div");
    item.className = "file-item";

    item.innerHTML = `
            <span class="file-name">${escapeHtml(file.name)}</span>
            <span class="file-type">${formatFileSize(file.size)}</span>
        `;

    fileList.appendChild(item);
  });
}

function formatFileSize(bytes) {
  if (bytes < 1024) {
    return `${bytes} B`;
  }

  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }

  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function escapeHtml(value) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

// Temporary local generation simulation
generateButton.addEventListener("click", async () => {
  if (selectedFiles.length === 0) {
    alert("Select at least one payroll file first.");
    return;
  }

  progressSection.classList.remove("hidden");
  generateButton.disabled = true;

  for (let progress = 0; progress <= 100; progress += 10) {
    await wait(150);

    progressBar.style.width = `${progress}%`;
    progressCount.textContent = `${progress}%`;

    if (progress < 100) {
      progressLabel.textContent = "Preparing generation...";
      progressMessage.textContent = "This is only a frontend test for now.";
    } else {
      progressLabel.textContent = "Ready";
      progressMessage.textContent =
        "Frontend test completed. Go will be connected next.";
    }
  }

  generateButton.disabled = false;
});

function wait(milliseconds) {
  return new Promise((resolve) => {
    setTimeout(resolve, milliseconds);
  });
}

// Temporary example PDF data
const examplePdfs = [
  {
    employee: "Gahin... Gaikwad",
    school: "ZPPS AADHALE",
    name: "September 2026.pdf",
    url: "#",
  },
  {
    employee: "Another Employee",
    school: "ZPPS Example",
    name: "Financial Year 2026–2027.pdf",
    url: "#",
  },
];

function renderPdfs(searchTerm = "") {
  const normalizedSearch = searchTerm.toLowerCase();

  const filteredPdfs = examplePdfs.filter((pdf) => {
    const searchableText = [pdf.employee, pdf.school, pdf.name]
      .join(" ")
      .toLowerCase();

    return searchableText.includes(normalizedSearch);
  });

  if (filteredPdfs.length === 0) {
    pdfList.innerHTML = `
            <p class="empty-state">
                No matching PDFs found.
            </p>
        `;
    return;
  }

  pdfList.innerHTML = "";

  filteredPdfs.forEach((pdf) => {
    const item = document.createElement("div");
    item.className = "pdf-item";

    item.innerHTML = `
            <div class="pdf-info">
                <strong class="pdf-name">
                    ${escapeHtml(pdf.name)}
                </strong>
                <div class="pdf-meta">
                    ${escapeHtml(pdf.employee)}
                    ·
                    ${escapeHtml(pdf.school)}
                </div>
            </div>

            <a
                class="pdf-link"
                href="${pdf.url}"
                target="_blank"
                rel="noopener"
            >
                Open PDF
            </a>
        `;

    pdfList.appendChild(item);
  });
}

searchInput.addEventListener("input", () => {
  renderPdfs(searchInput.value);
});

renderPdfs();

// async function testGoServer() {
//   const response = await fetch("http://localhost:8080/api/health");
//   const data = await response.json();

//   console.log(data);
// }

// testGoServer();
