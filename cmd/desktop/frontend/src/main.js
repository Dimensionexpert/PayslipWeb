import "./style.css";
import "./app.css";

import { GetEmployee } from "../wailsjs/go/main/App";

document.querySelector("#app").innerHTML = `
  <div class="employee-lookup">
    <h1>Employee Lookup</h1>

    <div class="search-row">
      <input
        id="shalarthId"
        type="text"
        placeholder="Enter Shalarth ID"
      />

      <button id="searchButton">Search</button>
    </div>

    <div id="result"></div>
  </div>
`;

const shalarthIdInput = document.getElementById("shalarthId");
const searchButton = document.getElementById("searchButton");
const resultElement = document.getElementById("result");

searchButton.addEventListener("click", async () => {
  const shalarthId = shalarthIdInput.value.trim();

  if (!shalarthId) {
    resultElement.innerText = "Enter a Shalarth ID.";
    return;
  }

  resultElement.innerText = "Searching...";

  try {
    const employee = await GetEmployee(shalarthId);

    resultElement.innerHTML = `
      <div class="employee-card">
        <h2>${employee.name}</h2>
        <p>${employee.designation}</p>
        <p>${employee.schoolName}</p>
        <p>UDISE: ${employee.udiseCode}</p>
        <p>Shalarth ID: ${employee.shalarthId}</p>
      </div>
    `;
  } catch (err) {
    console.error(err);
    resultElement.innerText = "Employee not found.";
  }
});
