package dto

type AppPaths struct {

	// Application resources shipped with the app.
	ResourcesDir string // Example: clusters.xlsx, monthly_template.xlsx, yearly_template.xlsx.
	DataDir      string // Persistent application data owned by the user.
	DBPath       string // SQLite database containing imported payroll data.
	SourceDir    string // Archived payroll XLSX files selected/imported by the user.
	CacheDir     string // Rebuildable generated XLSX cache.
}
