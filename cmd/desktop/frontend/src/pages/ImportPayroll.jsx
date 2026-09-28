import { useState } from "react";
import "./ImportPayroll.css";

import {
  SelectPayrollFile,
  ImportPayroll as ImportPayrollAPI,
} from "../../wailsjs/go/main/App";

function ImportPayroll({ onHome }) {
  const date = new Date();

  const [filePath, setFilePath] = useState("");
  const [month, setMonth] = useState(date.getMonth() + 1);
  const [year, setYear] = useState(date.getFullYear());
  const [status, setStatus] = useState("idle");
  const [feedback, setFeedback] = useState("");
  const [report, setReport] = useState(null);

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

  async function handleChooseFile() {
    try {
      const path = await SelectPayrollFile();

      if (path) {
        setFilePath(path);
      }
    } catch (error) {
      console.error("Failed to select payroll file:", error);
    }
  }

  async function handleImport() {
    if (!filePath) {
      setStatus("error");
      setFeedback("Please select a payroll file");
      return;
    }

    setStatus("loading");
    setFeedback("");
    setReport(null);

    try {
      const result = await ImportPayrollAPI(filePath, month, year);

      console.log("Import report:", result);

      setReport(result);
      setStatus("success");
    } catch (error) {
      console.error("Import failed:", error);

      setStatus("error");
      setFeedback(error?.message || String(error));
    }
  }
  return (
    <main className="import-payroll">
      <button className="payslip-back" onClick={onHome}>
        ← Back
      </button>

      <header className="import-header">
        <h1>Import Payroll</h1>
        <p>Select the payroll Excel file to import.</p>
      </header>

      <section className="import-file">
        <button onClick={handleChooseFile}>Choose XLSX File</button>

        {filePath && <p>{filePath}</p>}
      </section>

      <section className="import-fields">
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
      </section>

      <button
        className="import-action"
        onClick={handleImport}
        disabled={status === "loading"}
      >
        {status === "loading" ? "Importing..." : "Import"}
      </button>

      {status === "error" && (
        <div className="import-feedback error">{feedback}</div>
      )}

      {report && status === "success" && (
        <div className="import-report">
          <div className="import-report-title">
            ✓ Payroll imported successfully
          </div>

          <div className="import-report-summary">
            <div>
              <strong>{report.EmployeeCount}</strong>
              <span>Employees</span>
            </div>

            <div>
              <strong>{report.SchoolCount}</strong>
              <span>Schools</span>
            </div>

            <div>
              <strong>{report.ClusterMappingCount}</strong>
              <span>Cluster mappings</span>
            </div>
          </div>

          <div className="import-report-details">
            <div>
              <span>Missing clusters</span>
              <strong>{report.MissingClusterCount}</strong>
            </div>

            <div>
              <span>Unknown UDISE</span>
              <strong>{report.UnknownUDISECount}</strong>
            </div>

            <div>
              <span>Duplicate UDISE</span>
              <strong>{report.DuplicateUDISECount}</strong>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}

export default ImportPayroll;
