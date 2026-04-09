package ruletester

import (
	"encoding/json"
	"errors"

	minderv1 "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v3"
)

type RuleTestSuite struct {
	Version string     `yaml:"version"`
	Tests   []RuleTest `yaml:"cases"`
}

type RuleTest struct {
	Name       string               `yaml:"name"`
	Def        map[string]any       `yaml:"def"`
	Params     map[string]any       `yaml:"params"`
	Entity     EntityVersionWrapper `yaml:"entity"`
	Expect     string               `yaml:"expect"`
	ErrorText  string               `yaml:"error_text"`
	Git        *GitTest             `yaml:"git"`
	HTTP       *HTTPTest            `yaml:"http"`
	MockIngest map[string]any       `yaml:"mock_ingest"`
}

type EntityVersionWrapper struct {
	Type   string                    `yaml:"type"`
	Entity protoreflect.ProtoMessage `yaml:"entity"`
}

func (e *EntityVersionWrapper) UnmarshalYAML(value *yaml.Node) error {
	var entity map[string]any
	if err := value.Decode(&entity); err != nil {
		e.Type = "repository"
		e.Entity = &minderv1.Repository{}
		return nil
	}

	typ, ok := entity["type"]
	if !ok {
		e.Type = "repository"
		e.Entity = &minderv1.Repository{}
		return nil
	}

	e.Type, ok = typ.(string)
	if !ok {
		return errors.New("entity type field must be a string")
	}

	switch e.Type {
	case "repo", "repository":
		e.Entity = &minderv1.Repository{}
	default:
		e.Entity = &minderv1.EntityInstance{}
	}

	if entity["entity"] != nil {
		entityBytes, err := json.Marshal(entity["entity"])
		if err != nil {
			return err
		}
		return protojson.Unmarshal(entityBytes, e.Entity)
	}

	return nil
}

type GitTest struct {
	RepoBase string `yaml:"repo_base"`
}

type HTTPTest struct {
	Status   int               `yaml:"status"`
	Body     string            `yaml:"body"`
	BodyFile string            `yaml:"body_file"`
	Headers  map[string]string `yaml:"headers"`
}
