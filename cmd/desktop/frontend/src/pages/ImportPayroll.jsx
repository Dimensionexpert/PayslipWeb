import { useState } from "react";
import "./ImportPayroll.css";

import {
  SelectPayrollFile,
  ImportPayroll as ImportPayrollAPI,
  DeleteMonthlyPayslips,
} from "../../wailsjs/go/main/App";

function ImportPayroll({ onHome }) {
  const date = new Date();

  const [filePath, setFilePath] = useState("");
  const [month, setMonth] = useState(date.getMonth() + 1);
  const [year, setYear] = useState(date.getFullYear());
  const [status, setStatus] = useState("idle");
  const [feedback, setFeedback] = useState("");
  const [report, setReport] = useState(null);

  const [showDeletePanel, setShowDeletePanel] = useState(false);
  const [deleteMonth, setDeleteMonth] = useState(date.getMonth() + 1);
  const [deleteYear, setDeleteYear] = useState(date.getFullYear());
  const [deleteStatus, setDeleteStatus] = useState("idle");
  const [deleteFeedback, setDeleteFeedback] = useState("");
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

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

  const selectedDeleteMonth = months.find((m) => m.value === deleteMonth)?.name;

  const operationInProgress =
    status === "loading" || deleteStatus === "loading";

  async function handleChooseFile() {
    try {
      const path = await SelectPayrollFile();

      if (path) {
        setFilePath(path);
      }
    } catch (error) {
      console.error("Failed to select payroll file:", error);
      setStatus("error");
      setFeedback(error?.message || String(error));
    }
  }

  async function handleImport() {
    if (!filePath) {
      setStatus("error");
      setFeedback("Please select a payroll file.");
      return;
    }

    setStatus("loading");
    setFeedback("");
    setReport(null);

    try {
      const result = await ImportPayrollAPI(filePath, month, year);

      setReport(result);
      setStatus("success");
    } catch (error) {
      console.error("Import failed:", error);
      setStatus("error");
      setFeedback(error?.message || String(error));
    }
  }

  async function handleDelete() {
    setDeleteStatus("loading");
    setDeleteFeedback("");

    try {
      const rowsAffected = await DeleteMonthlyPayslips(deleteYear, deleteMonth);

      setDeleteStatus("success");
      setDeleteFeedback(
        `Deleted ${rowsAffected} payslip record${
          rowsAffected === 1 ? "" : "s"
        } for ${selectedDeleteMonth} ${deleteYear}.`,
      );

      setShowDeleteConfirm(false);
    } catch (error) {
      console.error("Failed to delete monthly payslips:", error);
      setDeleteStatus("error");
      setDeleteFeedback(error?.message || String(error));
    }
  }

  function handleOpenDeletePanel() {
    setShowDeletePanel(true);
    setDeleteStatus("idle");
    setDeleteFeedback("");
  }

  function handleCloseDeletePanel() {
    setShowDeletePanel(false);
    setDeleteStatus("idle");
    setDeleteFeedback("");
    setShowDeleteConfirm(false);
  }

  function handleDeleteMonthChange(value) {
    setDeleteMonth(Number(value));
    setDeleteStatus("idle");
    setDeleteFeedback("");
  }

  function handleDeleteYearChange(value) {
    setDeleteYear(Number(value));
    setDeleteStatus("idle");
    setDeleteFeedback("");
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
        <button onClick={handleChooseFile} disabled={operationInProgress}>
          Choose XLSX File
        </button>

        {filePath && <p>{filePath}</p>}
      </section>

      <section className="import-fields">
        <label>
          Month
          <select
            value={month}
            onChange={(event) => setMonth(Number(event.target.value))}
            disabled={operationInProgress}
          >
            {months.map((m) => (
              <option key={m.value} value={m.value}>
                {m.name}
              </option>
            ))}
          </select>
        </label>

        <label>
          Year
          <select
            value={year}
            onChange={(event) => setYear(Number(event.target.value))}
            disabled={operationInProgress}
          >
            {[2026, 2025, 2024].map((y) => (
              <option key={y} value={y}>
                {y}
              </option>
            ))}
          </select>
        </label>
      </section>

      <button
        className="import-action"
        onClick={handleImport}
        disabled={operationInProgress}
      >
        {status === "loading" ? "Importing..." : "Import"}
      </button>

      {status === "error" && (
        <div className="import-feedback error" role="alert">
          {feedback}
        </div>
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

      <section className="delete-month-section">
        <div className="delete-month-heading">
          <div>
            <h2>Imported the wrong data?</h2>
            <p>
              Remove payslip records for a specific month if an import went
              wrong.
            </p>
          </div>
        </div>

        {!showDeletePanel && (
          <button
            className="delete-action"
            onClick={handleOpenDeletePanel}
            disabled={operationInProgress}
          >
            Delete monthly records
          </button>
        )}

        {showDeletePanel && (
          <div className="delete-month-panel">
            <h3>Select records to delete</h3>

            <div className="delete-month-fields">
              <label>
                Month
                <select
                  value={deleteMonth}
                  onChange={(event) =>
                    handleDeleteMonthChange(event.target.value)
                  }
                  disabled={operationInProgress}
                >
                  {months.map((m) => (
                    <option key={m.value} value={m.value}>
                      {m.name}
                    </option>
                  ))}
                </select>
              </label>

              <label>
                Year
                <select
                  value={deleteYear}
                  onChange={(event) =>
                    handleDeleteYearChange(event.target.value)
                  }
                  disabled={operationInProgress}
                >
                  {[2026, 2025, 2024].map((y) => (
                    <option key={y} value={y}>
                      {y}
                    </option>
                  ))}
                </select>
              </label>
            </div>

            <p className="delete-warning">
              This deletes all payslip records for {selectedDeleteMonth}{" "}
              {deleteYear}, across every employee.
            </p>

            {deleteStatus === "success" && (
              <div className="delete-feedback success" role="status">
                ✓ {deleteFeedback}
              </div>
            )}

            {deleteStatus === "error" && (
              <div className="delete-feedback error" role="alert">
                {deleteFeedback}
              </div>
            )}

            <div className="delete-month-actions">
              <button
                className="delete-cancel"
                onClick={handleCloseDeletePanel}
                disabled={operationInProgress}
              >
                Close
              </button>

              <button
                className="delete-action"
                onClick={() => setShowDeleteConfirm(true)}
                disabled={operationInProgress}
              >
                Delete records
              </button>
            </div>
          </div>
        )}
      </section>

      {showDeleteConfirm && (
        <div
          className="delete-modal-backdrop"
          onMouseDown={(event) => {
            if (event.target === event.currentTarget) {
              setShowDeleteConfirm(false);
            }
          }}
        >
          <section
            className="delete-modal"
            role="alertdialog"
            aria-modal="true"
            aria-labelledby="delete-modal-title"
            aria-describedby="delete-modal-description"
          >
            <div className="delete-modal-icon">!</div>

            <h2 id="delete-modal-title">
              Delete {selectedDeleteMonth} {deleteYear} records?
            </h2>

            <p id="delete-modal-description">
              All payslip records for this period will be permanently deleted
              for every employee. This action cannot be undone.
            </p>

            <div className="delete-modal-period">
              <span>Period</span>
              <strong>
                {selectedDeleteMonth} {deleteYear}
              </strong>

              <span>Scope</span>
              <strong>All employees</strong>
            </div>

            <div className="delete-modal-actions">
              <button
                className="delete-cancel"
                onClick={() => setShowDeleteConfirm(false)}
                disabled={deleteStatus === "loading"}
              >
                Cancel
              </button>

              <button
                className="delete-modal-danger"
                onClick={handleDelete}
                disabled={deleteStatus === "loading"}
              >
                {deleteStatus === "loading" ? "Deleting..." : "Delete records"}
              </button>
            </div>
          </section>
        </div>
      )}
    </main>
  );
}

export default ImportPayroll;
