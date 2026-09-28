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
  useEffect(() => {
    async function loadClusters() {
      try {
        const result = await GetClusters();
        console.log("Clusters:", result);
        setClusters(result);
      } catch (error) {
        console.error("Failed to load clusters:", error);
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

                    setEmployees(result);
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
                  setSchools(result);
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
