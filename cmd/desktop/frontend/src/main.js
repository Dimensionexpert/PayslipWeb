import "./style.css";
import "./app.css";

import {
  GetSchool,
  GetEmployeesBySchool,
  GetPayslip,
} from "../wailsjs/go/main/App";

document.querySelector("#app").innerHTML = `
  <div class="employee-lookup">
    <h1>School Lookup</h1>

    <div class="search-row">
      <input
        id="udise"
        type="text"
        placeholder="Enter UDISE code"
      />

      <button id="schoolSearchButton">Search</button>
    </div>

    <div id="schoolResult"></div>
    <div id="employeesResult"></div>
    <div id="selectedEmployeeResult"></div>
    <div id="payslipResult"></div>
  </div>
`;

const udiseInput = document.getElementById("udise");
const schoolSearchButton = document.getElementById("schoolSearchButton");
const schoolResultElement = document.getElementById("schoolResult");
const employeesResultElement = document.getElementById("employeesResult");
const selectedEmployeeResultElement = document.getElementById(
  "selectedEmployeeResult",
);
const payslipResultElement = document.getElementById("payslipResult");

schoolSearchButton.addEventListener("click", async () => {
  const udise = udiseInput.value.trim();

  if (!udise) {
    schoolResultElement.innerText = "Enter a UDISE code.";
    employeesResultElement.innerHTML = "";
    selectedEmployeeResultElement.innerHTML = "";
    payslipResultElement.innerHTML = "";
    return;
  }

  schoolResultElement.innerText = "Searching...";
  employeesResultElement.innerHTML = "";
  selectedEmployeeResultElement.innerHTML = "";
  payslipResultElement.innerHTML = "";

  try {
    const school = await GetSchool(udise);

    schoolResultElement.innerHTML = `
      <div class="school-card">
        <h2>${school.name}</h2>
        <p>UDISE: ${school.udiseCode}</p>
      </div>
    `;

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

        renderSelectedEmployee(employee);
      });
    });
  } catch (err) {
    console.error("School lookup failed:", err);

    schoolResultElement.innerText = `Error: ${err}`;
    employeesResultElement.innerHTML = "";
    selectedEmployeeResultElement.innerHTML = "";
    payslipResultElement.innerHTML = "";
  }
});

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

  const currentYear = new Date().getFullYear();

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
                <option value="${index + 1}" ${
                  index + 1 === new Date().getMonth() + 1 ? "selected" : ""
                }>
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
                <option value="${year}">
                  ${year}
                </option>
              `,
            )
            .join("")}
        </select>
      </div>

      <button id="viewPayslipButton" class="payslip-button">
        View Payslip
      </button>
    </div>
  `;

  const viewPayslipButton = document.getElementById("viewPayslipButton");

  viewPayslipButton.addEventListener("click", async () => {
    const month = Number(document.getElementById("month").value);
    const year = Number(document.getElementById("year").value);

    payslipResultElement.innerText = "Loading payslip...";

    try {
      const payslip = await GetPayslip(employee.shalarthId, month, year);

      console.log("Payslip:", payslip);

      payslipResultElement.innerHTML = `
        <div class="employee-card">
          <h2>Salary</h2>
          <p>Employee: ${employee.name}</p>
          <p>Month: ${months[month - 1]} ${year}</p>

          <h1>₹${payslip.Payslip.EmployeeNetSalary}</h1>
        </div>
      `;
    } catch (err) {
      console.error("GetPayslip failed:", err);

      payslipResultElement.innerText = `Payslip error: ${err}`;
    }
  });
}
