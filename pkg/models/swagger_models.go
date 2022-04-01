package models

// swagger:parameters idOfSurvivorCreateEndpoint idOfSurvivorUpdateEndpoint
type _ struct {
	// in:body
	// required:true
	Body Survivor
}

// swagger:parameters idOfSurvivorInfectedEndpoint
type _ struct {
	// in:body
	// required:true
	Body SurvivorInfected
}

// swagger:parameters idOfSurvivorInfectionReport
type _ struct {
	// in:body
	// required:true
	Body InfectionReport
}

// swagger:parameters idOfReportCriteriaEndpoint
type _ struct {
	// in:path
	// criteria
	// required:true
	Criteria string `json:"criteria"`
}
