import { useState } from "react";
import "./PayslipSelector.css";

function PayslipSelector({ employee, mode = "view" }) {
  const date = new Date();

  const [type, setType] = useState("monthly");
  const [month, setMonth] = useState(date.getMonth() + 1);
  const [year, setYear] = useState(date.getFullYear());
  const [financialYear, setFinancialYear] = useState("2026-27");

  const months = [
    { value: 1, name: "January" },
    { value: 2, name: "February" },
    { value: 3, name: "March" },
    { value: 4, name: "April" },
    { value: 5, name: "May" },
    { value: 6, name: "June" },
    { value: 7, name: "July" },
    { value: 8, name: "August" },
    { value: 9, name: "September" },
    { value: 10, name: "October" },
    { value: 11, name: "November" },
    { value: 12, name: "December" },
  ];

  return (
    <div className="payslip-selector">
      <div className="payslip-employee">
        <strong>{employee.name}</strong>
        <span>{employee.id}</span>
      </div>

      <div className="payslip-type">
        <button
          className={type === "monthly" ? "active" : ""}
          onClick={() => setType("monthly")}
        >
          Monthly
        </button>

        <button
          className={type === "yearly" ? "active" : ""}
          onClick={() => setType("yearly")}
        >
          Yearly
        </button>
      </div>

      {type === "monthly" ? (
        <div className="payslip-fields">
          <label>
            Month
            <select
              value={month}
              onChange={(event) => setMonth(Number(event.target.value))}
            >
              {months.map((month) => (
                <option key={month.value} value={month.value}>
                  {month.name}
                </option>
              ))}
            </select>
          </label>

          <label>
            Year
            <select
              value={year}
              onChange={(event) => setYear(Number(event.target.value))}
            >
              <option value={2026}>2026</option>
              <option value={2025}>2025</option>
              <option value={2024}>2024</option>
            </select>
          </label>
        </div>
      ) : (
        <div className="payslip-fields">
          <label>
            Financial Year
            <select
              value={financialYear}
              onChange={(event) => setFinancialYear(event.target.value)}
            >
              <option value="2026-27">2026-27</option>
              <option value="2025-26">2025-26</option>
              <option value="2024-25">2024-25</option>
            </select>
          </label>
        </div>
      )}

      <button className="payslip-action">
        {mode === "view" ? "View Payslip" : "Generate Payslip"}
      </button>
    </div>
  );
}

export default PayslipSelector;
