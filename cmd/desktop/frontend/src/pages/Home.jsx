import "./Home.css";

function Home({ onNavigate }) {
  return (
    <main className="home">
      <header className="home-header">
        <div className="home-brand">
          <h1>Payslip</h1>
          <p>Payroll made simple.</p>
        </div>
      </header>

      <section className="home-menu">
        <button className="home-card" onClick={() => onNavigate("view")}>
          <div className="home-card-icon">⌕</div>

          <div>
            <strong>View Payslips</strong>
            <span>Find and view employee payslips.</span>
          </div>

          <span className="home-card-arrow">→</span>
        </button>

        <button className="home-card" onClick={() => onNavigate("generate")}>
          <div className="home-card-icon">＋</div>

          <div>
            <strong>Generate Payslips</strong>
            <span>Generate monthly or yearly payslips.</span>
          </div>

          <span className="home-card-arrow">→</span>
        </button>

        <button className="home-card" onClick={() => onNavigate("import")}>
          <div className="home-card-icon">⇩</div>

          <div>
            <strong>Import Payroll</strong>
            <span>Import a monthly payroll Excel file.</span>
          </div>

          <span className="home-card-arrow">→</span>
        </button>

        <button className="home-card" onClick={() => onNavigate("settings")}>
          <div className="home-card-icon">⚙</div>

          <div>
            <strong>Settings</strong>
            <span>Manage application settings.</span>
          </div>

          <span className="home-card-arrow">→</span>
        </button>
      </section>
    </main>
  );
}

export default Home;
