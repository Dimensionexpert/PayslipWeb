package resources

import "embed"

//go:embed clusters.xlsx monthly_template.xlsx yearly_template.xlsx
var ResourceFiles embed.FS
