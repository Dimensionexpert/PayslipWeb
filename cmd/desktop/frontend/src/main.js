import "./style.css";
import "./app.css";

import {
  GetSchoolsByCluster,
  GetEmployeesBySchool,
  GetPayslip,
  GenerateMonthlyPayslip,
  ChooseOutputDirectory,
} from "../wailsjs/go/main/App";

document.querySelector("#app").innerHTML = `
  <div class="employee-lookup">
    <h1>Payslip Lookup</h1>

    <div class="search-row">
      <input
        id="cluster"
        type="text"
        placeholder="Enter cluster name"
      />

      <button id="clusterSearchButton">Search</button>
    </div>

    <div id="schoolsResult"></div>
    <div id="employeesResult"></div>
    <div id="selectedEmployeeResult"></div>
    <div id="payslipResult"></div>
  </div>
`;

const clusterInput = document.getElementById("cluster");
const clusterSearchButton = document.getElementById("clusterSearchButton");

const schoolsResultElement = document.getElementById("schoolsResult");
const employeesResultElement = document.getElementById("employeesResult");
const selectedEmployeeResultElement = document.getElementById(
  "selectedEmployeeResult",
);
const payslipResultElement = document.getElementById("payslipResult");

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
          (employee) => employee.shalarthId === shalarthId,
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

      <div class="payslip-options">
        <label for="month">Month</label>

        <select id="month">
          ${months
            .map(
              (month, index) => `
                <option
                  value="${index + 1}"
                  ${index + 1 === currentMonth ? "selected" : ""}
                >
                  ${month}
                </option>
              `,
            )
            .join("")}
        </select>

        <label for="year">Year</label>

        <select id="year">
          ${Array.from({ length: 5 }, (_, index) => currentYear - index)
            .map(
              (year) => `
                <option
                  value="${year}"
                  ${year === currentYear ? "selected" : ""}
                >
                  ${year}
                </option>
              `,
            )
            .join("")}
        </select>
      </div>

      <div class="output-options">
        <button id="chooseOutputButton">
          Choose Output Folder
        </button>

        <p id="outputDirectory">No folder selected</p>
      </div>

      <button id="viewPayslipButton" class="payslip-button">
        View Payslip
      </button>

      <button id="generateMonthlyPayslipButton" class="payslip-button">
        Generate Monthly Payslip
      </button>
    </div>
  `;

  const viewPayslipButton = document.getElementById("viewPayslipButton");
  const chooseOutputButton = document.getElementById("chooseOutputButton");
  const outputDirectoryElement = document.getElementById("outputDirectory");
  const generateMonthlyPayslipButton = document.getElementById(
    "generateMonthlyPayslipButton",
  );

  let outputDir = "";

  chooseOutputButton.addEventListener("click", async () => {
    try {
      const selectedPath = await ChooseOutputDirectory();

      if (!selectedPath) {
        return;
      }

      outputDir = selectedPath;
      outputDirectoryElement.innerText = selectedPath;
    } catch (err) {
      console.error("ChooseOutputDirectory failed:", err);

      outputDirectoryElement.innerText = `Error: ${err}`;
    }
  });

  viewPayslipButton.addEventListener("click", async () => {
    const month = Number(document.getElementById("month").value);
    const year = Number(document.getElementById("year").value);

    payslipResultElement.innerText = "Loading payslip...";

    try {
      const payslip = await GetPayslip(employee.shalarthId, month, year);

      payslipResultElement.innerHTML = `
        <div class="employee-card">
          <h2>Salary</h2>

          <p>Employee: ${employee.name}</p>
          <p>Month: ${months[month - 1]} ${year}</p>

          <h1>₹${payslip.Payslip.EmployeeNetSalary}</h1>
        </div>
      `;

      console.log("Payslip:", payslip);
    } catch (err) {
      console.error("GetPayslip failed:", err);

      payslipResultElement.innerText = `Payslip error: ${err}`;
    }
  });

  generateMonthlyPayslipButton.addEventListener("click", async () => {
    if (!outputDir) {
      payslipResultElement.innerText = "Choose an output folder first.";
      return;
    }

    const month = Number(document.getElementById("month").value);
    const year = Number(document.getElementById("year").value);

    payslipResultElement.innerText = "Generating payslip...";

    try {
      const pdfPath = await GenerateMonthlyPayslip(
        employee.shalarthId,
        month,
        year,
        outputDir,
      );

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
}
