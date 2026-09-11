import "./style.css";
import "./app.css";

import logo from "./assets/images/logo-universal.png";
import { GetEmployees, GetPayslip } from "../wailsjs/go/main/App";

document.querySelector("#app").innerHTML = `
  <img id="logo" class="logo">

  <h1>Payslip</h1>

  <div id="result">Loading employees...</div>

  <ul id="employees"></ul>
`;

document.getElementById("logo").src = logo;

const resultElement = document.getElementById("result");
const employeesElement = document.getElementById("employees");

GetEmployees()
  .then((employees) => {
    resultElement.innerText = `Employees loaded: ${employees.length}`;

    employeesElement.innerHTML = employees
      .map((employee) => `<li>${employee.Name}</li>`)
      .join("");
  })
  .catch((err) => {
    console.error("Failed to get employees:", err);
    resultElement.innerText = "Failed to load employees";
  });

GetPayslip("02DEDAAMF6901", 8, 2026)
  .then((payslip) => {
    console.log("Payslip:", payslip);
  })
  .catch((err) => {
    console.error("Failed to get payslip:", err);
  });
