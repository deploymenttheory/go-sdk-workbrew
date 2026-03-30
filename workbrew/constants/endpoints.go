package constants

// ============================================================================
// Workbrew API Endpoints
// ============================================================================
// All endpoints are relative to the workspace base URL:
// https://console.workbrew.com/workspaces/{workspace_name}
// Workbrew API docs: https://console.workbrew.com/documentation/api

const (
	// Analytics endpoints
	EndpointAnalyticsJSON = "/analytics.json"
	EndpointAnalyticsCSV  = "/analytics.csv"

	// Brew Commands endpoints
	EndpointBrewCommandsJSON          = "/brew_commands.json"
	EndpointBrewCommandsCSV           = "/brew_commands.csv"
	EndpointBrewCommandRunsJSONFormat = "/brew_commands/%s/runs.json" // {brew_command_label}
	EndpointBrewCommandRunsCSVFormat  = "/brew_commands/%s/runs.csv"  // {brew_command_label}

	// Brew Configurations endpoints
	EndpointBrewConfigurationsJSON = "/brew_configurations.json"
	EndpointBrewConfigurationsCSV  = "/brew_configurations.csv"

	// Brewfiles endpoints
	EndpointBrewfilesJSON          = "/brewfiles.json"
	EndpointBrewfilesCSV           = "/brewfiles.csv"
	EndpointBrewfileLabelFormat    = "/brewfiles/%s.json"       // {label}
	EndpointBrewfileRunsJSONFormat = "/brewfiles/%s/runs.json"  // {label}
	EndpointBrewfileRunsCSVFormat  = "/brewfiles/%s/runs.csv"   // {label}

	// Brew Taps endpoints
	EndpointBrewTapsJSON = "/brew_taps.json"
	EndpointBrewTapsCSV  = "/brew_taps.csv"

	// Casks endpoints
	EndpointCasksJSON = "/casks.json"
	EndpointCasksCSV  = "/casks.csv"

	// Device Groups endpoints
	EndpointDeviceGroupsJSON = "/device_groups.json"
	EndpointDeviceGroupsCSV  = "/device_groups.csv"

	// Devices endpoints
	EndpointDevicesJSON = "/devices.json"
	EndpointDevicesCSV  = "/devices.csv"

	// Events endpoints
	EndpointEventsJSON = "/events.json"
	EndpointEventsCSV  = "/events.csv"

	// Formulae endpoints
	EndpointFormulaeJSON = "/formulae.json"
	EndpointFormulaeCSV  = "/formulae.csv"

	// Licenses endpoints
	EndpointLicensesJSON = "/licenses.json"
	EndpointLicensesCSV  = "/licenses.csv"

	// Vulnerabilities endpoints
	EndpointVulnerabilitiesJSON = "/vulnerabilities.json"
	EndpointVulnerabilitiesCSV  = "/vulnerabilities.csv"

	// Vulnerability Changes endpoints
	EndpointVulnerabilityChangesJSON = "/vulnerability_changes.json"
	EndpointVulnerabilityChangesCSV  = "/vulnerability_changes.csv"
)
