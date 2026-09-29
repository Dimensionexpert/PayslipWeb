import { useEffect, useState } from "react";
import PayslipSelector from "../components/PayslipSelector";
import {
  GetClusters,
  GetSchoolsByCluster,
  GetEmployeesBySchool,
  GenerateMonthlyPayslipsForScope,
} from "../../wailsjs/go/main/App";
import "./ViewPayslips.css";

function ViewPayslips({ mode = "view", onHome }) {
  const [search, setSearch] = useState("");
  const [clusters, setClusters] = useState([]);
  const [schools, setSchools] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [selectedCluster, setSelectedCluster] = useState(null);
  const [selectedSchool, setSelectedSchool] = useState(null);
  const [selectedEmployee, setSelectedEmployee] = useState(null);
  const [employeeAction, setEmployeeAction] = useState(null);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState("");
  const [schoolsLoading, setSchoolsLoading] = useState(false);
  const [schoolsError, setSchoolsError] = useState("");
  const [employeesLoading, setEmployeesLoading] = useState(false);
  const [employeesError, setEmployeesError] = useState("");
  useEffect(() => {
    async function loadClusters() {
      try {
        const result = await GetClusters();
        setClusters(result || []);
      } catch (error) {
        console.error("Failed to load clusters:", error);
        setLoadError(error?.message || String(error));
      } finally {
        setLoading(false);
      }
    }
    loadClusters();
  }, []);
  const filteredClusters = clusters.filter((cluster) =>
    cluster.name.toLowerCase().includes(search.toLowerCase()),
  );
  if (selectedSchool) {
    return (
      <main className="view-payslips">
        <button
          className="payslip-back"
          onClick={() => {
            setSelectedSchool(null);
            setSelectedEmployee(null);
            setEmployeeAction(null);
          }}
        >
          ← Back
        </button>

        <header>
          <h1>{selectedSchool.name}</h1>
          <p>Select an employee to view or generate their payslip.</p>
        </header>

        <section>
          <h2>Employees</h2>

          <div className="cluster-list">
            {employees.map((employee) => {
              const isSelected =
                selectedEmployee?.shalarthId === employee.shalarthId;

              return (
                <div
                  className={`employee-card ${isSelected ? "expanded" : ""}`}
                  key={employee.shalarthId}
                >
                  <button
                    className="employee-card-header"
                    onClick={() => {
                      if (isSelected) {
                        setSelectedEmployee(null);
                        setEmployeeAction(null);
                      } else {
                        setSelectedEmployee(employee);
                        setEmployeeAction(null);
                      }
                    }}
                  >
                    <div>
                      <strong>{employee.name}</strong>
                      <span>{employee.designation}</span>
                    </div>

                    <span className="cluster-arrow">
                      {isSelected ? "↑" : "→"}
                    </span>
                  </button>

                  {isSelected && (
                    <div className="employee-card-content">
                      {!employeeAction ? (
                        <div className="employee-actions">
                          <button onClick={() => setEmployeeAction("view")}>
                            View Payslip
                          </button>

                          <button onClick={() => setEmployeeAction("generate")}>
                            Generate Payslip
                          </button>
                        </div>
                      ) : (
                        <PayslipSelector
                          employee={selectedEmployee}
                          mode={employeeAction}
                        />
                      )}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </section>
      </main>
    );
  }

  if (selectedSchool) {
    if (employeesLoading) {
      return (
        <main className="view-payslips">
          <button className="payslip-back" onClick={onHome}>
            ← Back
          </button>

          <div className="view-state">
            <h2>Loading employees...</h2>
          </div>
        </main>
      );
    }

    if (employeesError) {
      return (
        <main className="view-payslips">
          <button className="payslip-back" onClick={onHome}>
            ← Back
          </button>

          <div className="view-state">
            <h2>Unable to load employees</h2>
            <p>{employeesError}</p>
          </div>
        </main>
      );
    }

    if (employees.length === 0) {
      return (
        <main className="view-payslips">
          <button
            className="payslip-back"
            onClick={() => setSelectedCluster(null)}
          >
            ← Back
          </button>

          <div className="view-state">
            <h2>No employee found</h2>
            <p>This school does not have any employee available.</p>
          </div>
        </main>
      );
    }
    return (
      <main className="view-payslips">
        <button
          className="payslip-back"
          onClick={() => setSelectedSchool(null)}
        >
          ← Back
        </button>

        <header>
          <h1>{selectedSchool.name}</h1>
          <p>Select an employee to view their payslip.</p>
        </header>

        <section>
          <h2>Employees</h2>

          <div className="cluster-list">
            {employees.map((employee) => (
              <button
                className="cluster-card"
                key={employee.shalarthId}
                onClick={() => setSelectedEmployee(employee)}
              >
                <div>
                  <strong>{employee.name}</strong>
                  <span>{employee.designation}</span>
                </div>

                <span className="cluster-arrow">→</span>
              </button>
            ))}
          </div>
        </section>
      </main>
    );
  }

  /* * Cluster selected */ if (selectedCluster) {
    if (schoolsLoading) {
      return (
        <main className="view-payslips">
          <button className="payslip-back" onClick={onHome}>
            ← Back
          </button>

          <div className="view-state">
            <h2>Loading schools...</h2>
          </div>
        </main>
      );
    }

    if (schoolsError) {
      return (
        <main className="view-payslips">
          <button className="payslip-back" onClick={onHome}>
            ← Back
          </button>

          <div className="view-state">
            <h2>Unable to load schools</h2>
            <p>{schoolsError}</p>
          </div>
        </main>
      );
    }

    if (schools.length === 0) {
      return (
        <main className="view-payslips">
          <button
            className="payslip-back"
            onClick={() => setSelectedCluster(null)}
          >
            ← Back
          </button>

          <div className="view-state">
            <h2>No schools found</h2>
            <p>This cluster does not have any schools available.</p>
          </div>
        </main>
      );
    }
    return (
      <main className="view-payslips">
        {" "}
        <button
          className="payslip-back"
          onClick={() => setSelectedCluster(null)}
        >
          {" "}
          ← Back{" "}
        </button>{" "}
        <header>
          {" "}
          <h1>{selectedCluster.name}</h1>{" "}
          <p>Select a school to view its employees.</p>{" "}
        </header>{" "}
        <section>
          {" "}
          <h2>Schools</h2>{" "}
          <div className="cluster-list">
            {" "}
            {schools.map((school) => (
              <button
                className="cluster-card"
                key={school.udiseCode}
                onClick={async () => {
                  try {
                    const result = await GetEmployeesBySchool(school.udiseCode);

                    console.log("Employees:", result);

                    setEmployees(result || []);
                    setSelectedSchool(school);
                  } catch (error) {
                    console.error("Failed to load employees:", error);
                  }
                }}
              >
                {" "}
                <div>
                  {" "}
                  <strong>{school.name}</strong>{" "}
                </div>{" "}
                <span className="cluster-arrow">→</span>{" "}
              </button>
            ))}{" "}
          </div>{" "}
        </section>{" "}
      </main>
    );
  }

  if (loading) {
    return (
      <main className="view-payslips">
        <button className="payslip-back" onClick={onHome}>
          ← Back
        </button>

        <div className="view-state">
          <h2>Loading payslips...</h2>
        </div>
      </main>
    );
  }

  if (loadError) {
    return (
      <main className="view-payslips">
        <button className="payslip-back" onClick={onHome}>
          ← Back
        </button>

        <div className="view-state">
          <h2>Unable to load payslips</h2>
          <p>{loadError}</p>
        </div>
      </main>
    );
  }

  if (clusters.length === 0) {
    return (
      <main className="view-payslips">
        <button className="payslip-back" onClick={onHome}>
          ← Back
        </button>

        <div className="view-state">
          <h2>No payroll data available</h2>
          <p>Import a payroll file to get started.</p>
        </div>
      </main>
    );
  }

  /* * Cluster list */ return (
    <main className="view-payslips">
      {" "}
      <button className="payslip-back" onClick={onHome}>
        ← Home
      </button>
      <header>
        {" "}
        <h1>View Payslips</h1>{" "}
        <p>Find an employee by following the school hierarchy.</p>{" "}
      </header>{" "}
      <div className="search-box">
        {" "}
        <span>⌕</span>{" "}
        <input
          type="text"
          placeholder="Search clusters..."
          value={search}
          onChange={(event) => setSearch(event.target.value)}
        />{" "}
      </div>{" "}
      <section>
        {" "}
        <h2>Clusters</h2>{" "}
        <div className="cluster-list">
          {" "}
          {filteredClusters.map((cluster) => (
            <button
              className="cluster-card"
              key={cluster.id}
              onClick={async () => {
                try {
                  const result = await GetSchoolsByCluster(cluster.name);
                  console.log("Schools:", result);
                  setSchools(result || []);
                  setSelectedCluster(cluster);
                } catch (error) {
                  console.error("Failed to load schools:", error);
                }
              }}
            >
              {" "}
              <div>
                {" "}
                <strong>{cluster.name}</strong>{" "}
              </div>{" "}
              <span className="cluster-arrow">→</span>{" "}
            </button>
          ))}{" "}
        </div>{" "}
      </section>{" "}
    </main>
  );
}
export default ViewPayslips;
