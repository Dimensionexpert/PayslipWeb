import { useState } from "react";
import "./PayslipSelector.css";

import {
  OpenMonthlyPayslip,
  OpenYearlyPayslip,
  GenerateMonthlyPayslip,
  GenerateYearlyPayslip,
  GenerateMonthlyPayslipsForScope,
  GenerateYearlyPayslipsForScope,
} from "../../wailsjs/go/main/App";

function PayslipSelector({
  employee = null,
  mode = "view",
  scope = "all",
  cluster = "",
  udise = "",
  schoolName = "",
}) {
  const isViewMode = mode === "view";
  const isBulkMode = mode === "bulk";

  const date = new Date();

  const [type, setType] = useState("monthly");
  const [month, setMonth] = useState(date.getMonth() + 1);
  const [year, setYear] = useState(date.getFullYear());
  const [financialYear, setFinancialYear] = useState("2026-27");

  const [status, setStatus] = useState("idle");
  const [feedback, setFeedback] = useState("");

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

  async function handleAction() {
    setStatus("loading");
    setFeedback("");

    try {
      // --------------------------------
      // View individual payslip
      // --------------------------------
      if (isViewMode) {
        if (type === "monthly") {
          await OpenMonthlyPayslip(employee.shalarthId, month, year);

          setFeedback("Payslip opened successfully");
        } else {
          const financialYearStart = Number(financialYear.split("-")[0]);

          await OpenYearlyPayslip(employee.shalarthId, financialYearStart);

          setFeedback("Yearly payslip opened successfully");
        }

        setStatus("success");
        return;
      }

      // --------------------------------
      // Bulk generation
      // --------------------------------
      if (isBulkMode) {
        if (type === "monthly") {
          const count = await GenerateMonthlyPayslipsForScope(
            scope,
            cluster,
            udise,
            month,
            year,
          );

          setFeedback(
            count === 0
              ? "All payslips are already generated"
              : `${count} payslip${count === 1 ? "" : "s"} generated successfully`,
          );
        } else {
          const financialYearStart = Number(financialYear.split("-")[0]);

          const count = await GenerateYearlyPayslipsForScope(
            scope,
            cluster,
            udise,
            financialYearStart,
          );

          setFeedback(
            count === 0
              ? "All yearly payslips are already generated"
              : `${count} payslip${count === 1 ? "" : "s"} generated successfully`,
          );
        }

        setStatus("success");
        return;
      }

      // --------------------------------
      // Generate individual payslip
      // --------------------------------
      if (type === "monthly") {
        const result = await GenerateMonthlyPayslip(
          employee.shalarthId,
          month,
          year,
        );

        console.log("Monthly payslip generated:", result);
        setFeedback("Monthly payslip generated successfully");
      } else {
        const financialYearStart = Number(financialYear.split("-")[0]);

        const result = await GenerateYearlyPayslip(
          employee.shalarthId,
          financialYearStart,
        );

        console.log("Yearly payslip generated:", result);
        setFeedback("Yearly payslip generated successfully");
      }

      setStatus("success");
    } catch (error) {
      console.error("Payslip action failed:", error);

      setStatus("error");
      setFeedback(error?.message || String(error));
    }
  }

  return (
    <div className="payslip-selector">
      <div className="payslip-employee">
        <strong>{employee ? employee.name : "Bulk Generation"}</strong>

        <span>
          {employee
            ? employee.shalarthId
            : scope === "all"
              ? "All clusters and schools"
              : scope === "cluster"
                ? cluster
                : schoolName || "Selected school"}
        </span>
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

      <button
        className="payslip-action"
        onClick={handleAction}
        disabled={status === "loading"}
      >
        {status === "loading"
          ? isViewMode
            ? "Opening..."
            : "Generating..."
          : isViewMode
            ? "View Payslip"
            : isBulkMode
              ? "Generate Payslips"
              : "Generate Payslip"}
      </button>

      {status !== "idle" && (
        <div className={`payslip-feedback ${status}`}>{feedback}</div>
      )}
    </div>
  );
}

export default PayslipSelector;
