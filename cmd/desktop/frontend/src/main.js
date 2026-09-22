import "./style.css";
import "./app.css";

import {
  GetSchoolsByCluster,
  GetEmployeesBySchool,
  GetOutputDirectory,
  GenerateMonthlyPayslip,
  GenerateYearlyPayslip,
  GenerateMonthlyPayslips,
  GenerateYearlyPayslips,
  SetOutputDirectory,
  OpenMonthlyPayslip,
  OpenYearlyPayslip,
} from "../wailsjs/go/main/App";

const monthNames = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

const currentDate = new Date();
const currentMonth = currentDate.getMonth() + 1;
const currentYear = currentDate.getFullYear();

document.querySelector("#app").innerHTML = `
  <div class="bulk-generation">
    <h2>Bulk Generation</h2>

    <div class="payslip-options">
      <label for="bulkMonth">Month</label>
      <select id="bulkMonth">
        ${monthNames
          .map(
            (month, index) => `
              <option value="${index + 1}" ${index + 1 === currentMonth ? "selected" : ""}>
                ${month}
              </option>
            `,
          )
          .join("")}
      </select>

      <label for="bulkYear">Year</label>
      <select id="bulkYear">
        ${Array.from({ length: 5 }, (_, index) => currentYear - index)
          .map(
            (year) => `
              <option value="${year}" ${year === currentYear ? "selected" : ""}>
                ${year}
              </option>
            `,
          )
          .join("")}
      </select>
    </div>

    <button id="generateMonthlyPayslipsButton">
      Generate All Monthly Payslips
    </button>

    <div class="payslip-options">
      <label for="bulkFinancialYear">Financial Year</label>
      <select id="bulkFinancialYear">
        ${Array.from({ length: 5 }, (_, index) => currentYear - index)
          .map(
            (year) => `
              <option value="${year}" ${year === currentYear ? "selected" : ""}>
                ${year}-${year + 1}
              </option>
            `,
          )
          .join("")}
      </select>
    </div>

    <button id="generateYearlyPayslipsButton">
      Generate All Yearly Payslips
    </button>
  </div>

  <button id="setOutputDirectoryButton">
    Set Output Directory
  </button>

  <p id="configuredOutputDirectory">Loading...</p>

  <div class="employee-lookup">
    <h1>Payslip Lookup</h1>

    <div class="search-row">
      <input id="cluster" type="text" placeholder="Enter cluster name" />
      <button id="clusterSearchButton">Search</button>
    </div>

    <div id="schoolsResult"></div>
    <div id="employeesResult"></div>
    <div id="selectedEmployeeResult"></div>
    <div id="payslipResult"></div>
  </div>
`;

const setOutputDirectoryButton = document.getElementById("setOutputDirectoryButton");
const configuredOutputDirectory = document.getElementById("configuredOutputDirectory");
const generateMonthlyPayslipsButton = document.getElementById(
  "generateMonthlyPayslipsButton",
);
const generateYearlyPayslipsButton = document.getElementById(
  "generateYearlyPayslipsButton",
);
const clusterInput = document.getElementById("cluster");
const clusterSearchButton = document.getElementById("clusterSearchButton");
const schoolsResultElement = document.getElementById("schoolsResult");
const employeesResultElement = document.getElementById("employeesResult");
const selectedEmployeeResultElement = document.getElementById(
  "selectedEmployeeResult",
);
const payslipResultElement = document.getElementById("payslipResult");

async function loadOutputDirectory() {
  try {
    const path = await GetOutputDirectory();

    if (!path) {
      configuredOutputDirectory.innerText = "No output directory configured";
      return;
    }

    configuredOutputDirectory.innerText = path;
  } catch (err) {
    console.error("GetOutputDirectory failed:", err);
    configuredOutputDirectory.innerText = `Error: ${err}`;
  }
}

function bindBulkGenerationHandlers() {
  generateMonthlyPayslipsButton.addEventListener("click", async () => {
    const month = Number(document.getElementById("bulkMonth").value);
    const year = Number(document.getElementById("bulkYear").value);

    payslipResultElement.innerText = "Generating all monthly payslips...";

    try {
      await GenerateMonthlyPayslips(month, year);
      payslipResultElement.innerText = "All monthly payslips generated successfully.";
    } catch (err) {
      console.error("GenerateMonthlyPayslips failed:", err);
      payslipResultElement.innerText = `Bulk monthly generation error: ${err}`;
    }
  });

  generateYearlyPayslipsButton.addEventListener("click", async () => {
    const year = Number(document.getElementById("bulkFinancialYear").value);

    payslipResultElement.innerText = "Generating all yearly payslips...";

    try {
      await GenerateYearlyPayslips(year);
      payslipResultElement.innerText = "All yearly payslips generated successfully.";
    } catch (err) {
      console.error("GenerateYearlyPayslips failed:", err);
      payslipResultElement.innerText = `Bulk yearly generation error: ${err}`;
    }
  });
}

setOutputDirectoryButton.addEventListener("click", async () => {
  try {
    const path = await SetOutputDirectory();

    if (!path) {
      return;
    }

    configuredOutputDirectory.innerText = path;
  } catch (err) {
    console.error("SetOutputDirectory failed:", err);
    configuredOutputDirectory.innerText = `Error: ${err}`;
  }
});

loadOutputDirectory();
bindBulkGenerationHandlers();

clusterSearchButton.addEventListener("click", async () => {
  const cluster = clusterInput.value.trim();

  if (!cluster) {
    schoolsResultElement.innerText = "Enter a cluster name.";
    employeesResultElement.innerHTML = "";
    selectedEmployeeResultElement.innerHTML = "";
    payslipResultElement.innerHTML = "";
    return;
  }

  schoolsResultElement.innerText = "Searching...";
  employeesResultElement.innerHTML = "";
  selectedEmployeeResultElement.innerHTML = "";
  payslipResultElement.innerHTML = "";

  try {
    const schools = await GetSchoolsByCluster(cluster);

    if (schools.length === 0) {
      schoolsResultElement.innerText = "No schools found.";
      return;
    }

    schoolsResultElement.innerHTML = `
      <h2>Schools</h2>
      <div class="employees-list">
        ${schools
          .map(
            (school) => `
              <button
                class="employee-card school-button"
                data-udise="${school.udiseCode}"
              >
                <strong>${school.name}</strong>
                <span>UDISE: ${school.udiseCode}</span>
              </button>
            `,
          )
          .join("")}
      </div>
    `;

    document.querySelectorAll(".school-button").forEach((button) => {
      button.addEventListener("click", async () => {
        const udise = button.dataset.udise;
        await loadEmployees(udise);
      });
    });
  } catch (err) {
    console.error("GetSchoolsByCluster failed:", err);
    schoolsResultElement.innerText = `Error: ${err}`;
  }
});

async function loadEmployees(udise) {
  employeesResultElement.innerText = "Loading employees...";
  selectedEmployeeResultElement.innerHTML = "";
  payslipResultElement.innerHTML = "";

  try {
    const employees = await GetEmployeesBySchool(udise);

    if (employees.length === 0) {
      employeesResultElement.innerText = "No employees found.";
      return;
    }

    employeesResultElement.innerHTML = `
      <h2>Employees</h2>
      <div class="employees-list">
        ${employees
          .map(
            (employee) => `
              <button
                class="employee-card employee-button"
                data-shalarth-id="${employee.shalarthId}"
              >
                <strong>${employee.name}</strong>
                <span>${employee.designation}</span>
                <small>${employee.shalarthId}</small>
              </button>
            `,
          )
          .join("")}
      </div>
    `;

    document.querySelectorAll(".employee-button").forEach((button) => {
      button.addEventListener("click", () => {
        const shalarthId = button.dataset.shalarthId;
        const employee = employees.find(
          (item) => item.shalarthId === shalarthId,
        );

        if (employee) {
          renderSelectedEmployee(employee);
        }
      });
    });
  } catch (err) {
    console.error("GetEmployeesBySchool failed:", err);
    employeesResultElement.innerText = `Error: ${err}`;
  }
}

function renderSelectedEmployee(employee) {
  const months = [
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
  ];

  const currentDate = new Date();
  const currentMonth = currentDate.getMonth() + 1;
  const currentYear = currentDate.getFullYear();

  selectedEmployeeResultElement.innerHTML = `
    <div class="selected-employee">
      <h2>Selected Employee</h2>

      <div class="selected-employee-info">
        <h3>${employee.name}</h3>
        <p>${employee.designation}</p>
        <p>Shalarth ID: ${employee.shalarthId}</p>
      </div>

      <div class="payslip-type">
        <button id="monthlyTab" class="payslip-type-button active">Monthly</button>
        <button id="yearlyTab" class="payslip-type-button">Yearly</button>
      </div>

      <div id="monthlyOptions">
        <div class="payslip-options">
          <label for="month">Month</label>
          <select id="month">
            ${months
              .map(
                (month, index) => `
                  <option value="${index + 1}" ${index + 1 === currentMonth ? "selected" : ""}>
                    ${month}
                  </option>
                `,
              )
              .join("")}
          </select>

          <label for="monthlyYear">Year</label>
          <select id="monthlyYear">
            ${Array.from({ length: 5 }, (_, index) => currentYear - index)
              .map(
                (year) => `
                  <option value="${year}" ${year === currentYear ? "selected" : ""}>
                    ${year}
                  </option>
                `,
              )
              .join("")}
          </select>
        </div>

        <button id="viewPayslipButton" class="payslip-button">View Payslip</button>
        <button id="generateMonthlyPayslipButton" class="payslip-button">
          Generate Monthly Payslip
        </button>
      </div>

      <div id="yearlyOptions" style="display: none;">
        <div class="payslip-options">
          <label for="financialYear">Financial Year</label>
          <select id="financialYear">
            ${Array.from({ length: 5 }, (_, index) => currentYear - index)
              .map(
                (year) => `
                  <option value="${year}" ${year === currentYear ? "selected" : ""}>
                    ${year}-${year + 1}
                  </option>
                `,
              )
              .join("")}
          </select>
        </div>

        <button id="viewYearlyPayslipButton" class="payslip-button">View Payslip</button>
        <button id="generateYearlyPayslipButton" class="payslip-button">
          Generate Yearly Payslip
        </button>
      </div>
    </div>
  `;

  const monthlyTab = document.getElementById("monthlyTab");
  const yearlyTab = document.getElementById("yearlyTab");
  const monthlyOptions = document.getElementById("monthlyOptions");
  const yearlyOptions = document.getElementById("yearlyOptions");

  monthlyTab.addEventListener("click", () => {
    monthlyOptions.style.display = "block";
    yearlyOptions.style.display = "none";
    monthlyTab.classList.add("active");
    yearlyTab.classList.remove("active");
  });

  yearlyTab.addEventListener("click", () => {
    monthlyOptions.style.display = "none";
    yearlyOptions.style.display = "block";
    monthlyTab.classList.remove("active");
    yearlyTab.classList.add("active");
  });

  const viewPayslipButton = document.getElementById("viewPayslipButton");
  const generateMonthlyPayslipButton = document.getElementById(
    "generateMonthlyPayslipButton",
  );

  viewPayslipButton.addEventListener("click", async () => {
    const month = Number(document.getElementById("month").value);
    const year = Number(document.getElementById("monthlyYear").value);

    payslipResultElement.innerText = "Opening payslip...";

    try {
      await OpenMonthlyPayslip(employee.shalarthId, month, year);
      payslipResultElement.innerText = "Payslip opened successfully.";
    } catch (err) {
      console.error("OpenMonthlyPayslip failed:", err);
      payslipResultElement.innerText = `Payslip error: ${err}`;
    }
  });

  generateMonthlyPayslipButton.addEventListener("click", async () => {
    const month = Number(document.getElementById("month").value);
    const year = Number(document.getElementById("monthlyYear").value);

    payslipResultElement.innerText = "Generating payslip...";

    try {
      const pdfPath = await GenerateMonthlyPayslip(employee.shalarthId, month, year);

      payslipResultElement.innerHTML = `
        <div class="employee-card">
          <h2>Payslip Generated</h2>
          <p>${pdfPath}</p>
        </div>
      `;

      console.log("Generated PDF:", pdfPath);
    } catch (err) {
      console.error("GenerateMonthlyPayslip failed:", err);
      payslipResultElement.innerText = `Generation error: ${err}`;
    }
  });

  const viewYearlyPayslipButton = document.getElementById(
    "viewYearlyPayslipButton",
  );
  const generateYearlyPayslipButton = document.getElementById(
    "generateYearlyPayslipButton",
  );

  viewYearlyPayslipButton.addEventListener("click", async () => {
    const financialYear = Number(document.getElementById("financialYear").value);

    payslipResultElement.innerText = "Opening yearly payslip...";

    try {
      await OpenYearlyPayslip(employee.shalarthId, financialYear);
      payslipResultElement.innerText = "Yearly payslip opened successfully.";
    } catch (err) {
      console.error("OpenYearlyPayslip failed:", err);
      payslipResultElement.innerText = `Yearly payslip error: ${err}`;
    }
  });

  generateYearlyPayslipButton.addEventListener("click", async () => {
    const financialYear = Number(document.getElementById("financialYear").value);

    payslipResultElement.innerText = "Generating yearly payslip...";

    try {
      const pdfPath = await GenerateYearlyPayslip(employee.shalarthId, financialYear);

      payslipResultElement.innerHTML = `
        <div class="employee-card">
          <h2>Yearly Payslip Generated</h2>
          <p>${pdfPath}</p>
        </div>
      `;

      console.log("Generated yearly PDF:", pdfPath);
    } catch (err) {
      console.error("GenerateYearlyPayslip failed:", err);
      payslipResultElement.innerText = `Yearly generation error: ${err}`;
    }
  });
}
