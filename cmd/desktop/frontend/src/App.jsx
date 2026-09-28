import { useState } from "react";
import Home from "./pages/Home";
import ViewPayslips from "./pages/ViewPayslips";
import BulkGenerate from "./pages/BulkGenerate";
import "./App.css";

function App() {
  const [page, setPage] = useState("home");

  if (page === "view") {
    return <ViewPayslips mode="view" onHome={() => setPage("home")} />;
  }

  if (page === "generate") {
    return <BulkGenerate onHome={() => setPage("home")} />;
  }

  if (page === "settings") {
    return (
      <main className="app-page">
        <button className="payslip-back" onClick={() => setPage("home")}>
          ← Home
        </button>

        <h1>Settings</h1>
        <p>Settings will be added here.</p>
      </main>
    );
  }

  return <Home onNavigate={setPage} />;
}

export default App;
