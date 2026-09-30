import { useState } from "react";
import Home from "./pages/Home";
import ViewPayslips from "./pages/ViewPayslips";
import BulkGenerate from "./pages/BulkGenerate";
import ImportPayroll from "./pages/ImportPayroll";
import Settings from "./pages/Settings";
import "./App.css";

function App() {
  const [page, setPage] = useState("home");

  const renderPage = () => {
    switch (page) {
      case "view":
        return <ViewPayslips mode="view" onHome={() => setPage("home")} />;
      case "generate":
        return <BulkGenerate onHome={() => setPage("home")} />;
      case "import":
        return <ImportPayroll onHome={() => setPage("home")} />;
      case "settings":
        return <Settings onHome={() => setPage("home")} />;
      default:
        return <Home onNavigate={setPage} />;
    }
  };

  return <div className="app">{renderPage()}</div>;
}

export default App;
