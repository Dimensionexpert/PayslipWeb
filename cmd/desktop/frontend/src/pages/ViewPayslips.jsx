import { useState } from "react";
import PayslipSelector from "../components/PayslipSelector";

const clusters = [
  {
    name: "INDORI",
    schools: 12,
    employees: 47,
  },
  {
    name: "MAVAL",
    schools: 18,
    employees: 93,
  },
  {
    name: "LONAVALA",
    schools: 9,
    employees: 31,
  },
  {
    name: "KAMSHET",
    schools: 7,
    employees: 24,
  },
];

const schools = {
  INDORI: [
    {
      name: "ZPPS Indori",
      employees: 7,
    },
    {
      name: "ZPPS Aamby",
      employees: 4,
    },
    {
      name: "ZPPS Kamshet",
      employees: 5,
    },
  ],

  MAVAL: [
    {
      name: "ZPPS Maval",
      employees: 8,
    },
    {
      name: "ZPPS Somatane",
      employees: 6,
    },
  ],
};

const employees = {
  "ZPPS Indori": [
    { name: "Rahul Patil", id: "EMP001" },
    { name: "Sneha Jadhav", id: "EMP002" },
    { name: "Amit Shinde", id: "EMP003" },
  ],

  "ZPPS Aamby": [
    { name: "Prakash Pawar", id: "EMP004" },
    { name: "Madhuri More", id: "EMP005" },
  ],
};

function ViewPayslips() {
  const [search, setSearch] = useState("");
  const [selectedCluster, setSelectedCluster] = useState(null);
  const [selectedSchool, setSelectedSchool] = useState(null);
  const [selectedEmployee, setSelectedEmployee] = useState(null);

  const filteredClusters = clusters.filter((cluster) =>
    cluster.name.toLowerCase().includes(search.toLowerCase()),
  );

  if (selectedEmployee) {
    return (
      <main className="view-payslips">
        <button
          className="payslip-back"
          onClick={() => setSelectedEmployee(null)}
        >
          ← Back
        </button>

        <header>
          <h1>Payslip</h1>
          <p>Select the payslip you want to view.</p>
        </header>

        <PayslipSelector employee={selectedEmployee} />
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
          <h1>{selectedSchool}</h1>
          <p>Select an employee to view their payslip.</p>
        </header>

        <section>
          <h2>Employees</h2>

          <div className="cluster-list">
            {employees[selectedSchool].map((employee) => (
              <button
                className="cluster-card"
                key={employee.id}
                onClick={() => setSelectedEmployee(employee)}
              >
                <div>
                  <strong>{employee.name}</strong>

                  <span>{employee.id}</span>
                </div>

                <span className="cluster-arrow">→</span>
              </button>
            ))}
          </div>
        </section>
      </main>
    );
  }

  if (selectedCluster) {
    return (
      <main className="view-payslips">
        <button
          className="payslip-back"
          onClick={() => setSelectedCluster(null)}
        >
          ← Back
        </button>

        <header>
          <h1>{selectedCluster}</h1>
          <p>Select a school to view its employees.</p>
        </header>

        <section>
          <h2>Schools</h2>

          <div className="cluster-list">
            {schools[selectedCluster].map((school) => (
              <button
                className="cluster-card"
                key={school.name}
                onClick={() => setSelectedSchool(school.name)}
              >
                <div>
                  <strong>{school.name}</strong>

                  <span>{school.employees} employees</span>
                </div>

                <span className="cluster-arrow">→</span>
              </button>
            ))}
          </div>
        </section>
      </main>
    );
  }

  return (
    <main className="view-payslips">
      <header>
        <h1>View Payslips</h1>
        <p>Find an employee by following the school hierarchy.</p>
      </header>

      <div className="search-box">
        <span>⌕</span>

        <input
          type="text"
          placeholder="Search clusters..."
          value={search}
          onChange={(event) => setSearch(event.target.value)}
        />
      </div>

      <section>
        <h2>Clusters</h2>

        <div className="cluster-list">
          {filteredClusters.map((cluster) => (
            <button
              className="cluster-card"
              key={cluster.name}
              onClick={() => setSelectedCluster(cluster.name)}
            >
              <div>
                <strong>{cluster.name}</strong>

                <span>
                  {cluster.schools} schools · {cluster.employees} employees
                </span>
              </div>

              <span className="cluster-arrow">→</span>
            </button>
          ))}
        </div>
      </section>
    </main>
  );
}

export default ViewPayslips;
