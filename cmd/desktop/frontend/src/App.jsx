import { useState } from "react";
import Home from "./pages/Home";
import ViewPayslips from "./pages/ViewPayslips";
import BulkGenerate from "./pages/BulkGenerate";
import ImportPayroll from "./pages/ImportPayroll";
import Settings from "./pages/Settings";
import "./App.css";

function App() {
  const [page, setPage] = useState("home");

  if (page === "view") {
    return <ViewPayslips mode="view" onHome={() => setPage("home")} />;
  }

  if (page === "generate") {
    return <BulkGenerate onHome={() => setPage("home")} />;
  }
  if (page === "import") {
    return <ImportPayroll onHome={() => setPage("home")} />;
  }

  if (page === "settings") {
    return <Settings onHome={() => setPage("home")} />;
  }

  return <Home onNavigate={setPage} />;
}

export default App;
