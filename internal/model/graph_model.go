package model

type (
	Declaration struct {
		Package   string
		Name      string
		UniqueId  string
		External  bool
		Content   *string
		File      *string
		StartLine *int
		EndLine   *int
	}
	AstFunction struct {
		Declaration
		Receiver *string
	}

	AstStruct struct {
		Declaration
	}

	AstVariable struct {
		Declaration
	}

	BaseRelation struct {
		SourceElementId string
		TargetElementId string
		RelationType    string
	}

	RelationType string
)

const (
	INVOKE     RelationType = "INVOKE"
	REFERENCE  RelationType = "REFERENCE"
	ASSOCIATE  RelationType = "ASSOCIATE"
	DEPENDENCE RelationType = "DEPENDENCE"
)

func (receiver RelationType) String() string {
	return string(receiver)
}
