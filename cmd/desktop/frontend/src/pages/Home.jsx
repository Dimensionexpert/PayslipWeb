import "./Home.css";

function Home({ onNavigate }) {
  return (
    <main className="home">
      <header className="home-header">
        <div className="home-brand">
          <h1>Payslip</h1>
          <p>Keep payroll workflows clear, fast, and consistent.</p>
        </div>
      </header>

      <section className="home-menu" aria-label="Main menu">
        <button
          className="home-card"
          data-tone="blue"
          onClick={() => onNavigate("view")}
        >
          <div className="home-card-icon">⌕</div>

          <div>
            <strong>View Payslips</strong>
            <span>Find and review employee payslips.</span>
          </div>

          <span className="home-card-arrow">→</span>
        </button>

        <button
          className="home-card"
          data-tone="teal"
          onClick={() => onNavigate("generate")}
        >
          <div className="home-card-icon">＋</div>

          <div>
            <strong>Generate Payslips</strong>
            <span>Generate monthly or yearly payroll batches.</span>
          </div>

          <span className="home-card-arrow">→</span>
        </button>

        <button
          className="home-card"
          data-tone="violet"
          onClick={() => onNavigate("import")}
        >
          <div className="home-card-icon">⇩</div>

          <div>
            <strong>Import Payroll</strong>
            <span>Import payroll data from Excel files.</span>
          </div>

          <span className="home-card-arrow">→</span>
        </button>

        <button
          className="home-card"
          data-tone="amber"
          onClick={() => onNavigate("settings")}
        >
          <div className="home-card-icon">⚙</div>

          <div>
            <strong>Settings</strong>
            <span>Manage app configuration and paths.</span>
          </div>

          <span className="home-card-arrow">→</span>
        </button>
      </section>
    </main>
  );
}

export default Home;
