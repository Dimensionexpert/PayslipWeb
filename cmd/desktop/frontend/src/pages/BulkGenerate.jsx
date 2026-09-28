import { useEffect, useState } from "react";
import "./BulkGenerate.css";

import { GetClusters, GetSchoolsByCluster } from "../../wailsjs/go/main/App";

import PayslipSelector from "../components/PayslipSelector";

function BulkGenerate({ onHome }) {
  const [scope, setScope] = useState("all");

  const [clusters, setClusters] = useState([]);
  const [schools, setSchools] = useState([]);

  const [selectedCluster, setSelectedCluster] = useState("");
  const [selectedSchool, setSelectedSchool] = useState("");

  useEffect(() => {
    if (scope === "cluster" || scope === "school") {
      GetClusters().then(setClusters).catch(console.error);
    }
  }, [scope]);

  useEffect(() => {
    if (scope !== "school" || !selectedCluster) {
      setSchools([]);
      setSelectedSchool("");
      return;
    }

    GetSchoolsByCluster(selectedCluster).then(setSchools).catch(console.error);
  }, [scope, selectedCluster]);

  return (
    <main className="bulk-generate">
      <button className="payslip-back" onClick={onHome}>
        ← Back
      </button>

      <header className="bulk-header">
        <h1>Generate Payslips</h1>
        <p>Generate payslips for multiple employees at once.</p>
      </header>

      <section className="bulk-section">
        <h2>Generate For</h2>

        <div className="bulk-options">
          <button
            className={`bulk-option ${scope === "all" ? "selected" : ""}`}
            onClick={() => {
              setScope("all");
              setSelectedCluster("");
              setSelectedSchool("");
            }}
          >
            <strong>All</strong>
            <span>All clusters and schools.</span>
          </button>

          <button
            className={`bulk-option ${scope === "cluster" ? "selected" : ""}`}
            onClick={() => {
              setScope("cluster");
              setSelectedSchool("");
            }}
          >
            <strong>Cluster</strong>
            <span>A complete cluster.</span>
          </button>

          <button
            className={`bulk-option ${scope === "school" ? "selected" : ""}`}
            onClick={() => {
              setScope("school");
            }}
          >
            <strong>School</strong>
            <span>A single school.</span>
          </button>
        </div>
      </section>

      {scope === "cluster" && (
        <section className="bulk-picker">
          <label>Cluster</label>

          <select
            value={selectedCluster}
            onChange={(event) => setSelectedCluster(event.target.value)}
          >
            <option value="">Select a cluster</option>

            {clusters.map((cluster) => (
              <option key={cluster.id} value={cluster.name}>
                {cluster.name}
              </option>
            ))}
          </select>
        </section>
      )}

      {scope === "school" && (
        <section className="bulk-picker">
          <label>Cluster</label>

          <select
            value={selectedCluster}
            onChange={(event) => setSelectedCluster(event.target.value)}
          >
            <option value="">Select a cluster</option>

            {clusters.map((cluster) => (
              <option key={cluster.id} value={cluster.name}>
                {cluster.name}
              </option>
            ))}
          </select>

          {selectedCluster && (
            <>
              <label>School</label>

              <select
                value={selectedSchool}
                onChange={(event) => setSelectedSchool(event.target.value)}
              >
                <option value="">Select a school</option>

                {schools.map((school) => (
                  <option key={school.udiseCode} value={school.udiseCode}>
                    {school.name}
                  </option>
                ))}
              </select>
            </>
          )}
        </section>
      )}

      <section className="bulk-selector">
        <PayslipSelector
          mode="bulk"
          scope={scope}
          cluster={selectedCluster}
          udise={selectedSchool}
          schoolName={
            schools.find((school) => school.udiseCode === selectedSchool)
              ?.name || ""
          }
        />
      </section>
    </main>
  );
}

export default BulkGenerate;
