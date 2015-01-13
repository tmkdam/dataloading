package bloomsource

import (
	"testing"
	"github.com/gocodo/bloomsource/tests"
)

var sourceMapping = &SourceMapping{
	Sources: []Mapping{
		Mapping{
			Name: "source",
			Destinations: []Destination{
				Destination{
					Name: "dest",
					Fields: []MappingField{
						MappingField{
							Source: []interface{}{ "sone" },
							Dest: "id",
						},
						MappingField{
							Source: "stwo",
							Dest: "two",
							Type: "bigint",
						},
					},
				},
			},
		},
	},
}

func TestMappingToCreate(t *testing.T) {
	spec := tests.Spec(t)
	create := MappingToCreate(sourceMapping)
	spec.Expect(create).ToEqual(`CREATE TABLE dest(
id uuid,
two bigint
);
INSERT INTO sources (id, name) VALUES ('c3f70f05-9179-37f5-93b8-1b8f43d291c7', 'source');
`)
}

func TestMappingToDrop(t *testing.T) {
	spec := tests.Spec(t)
	create := MappingToDrop(sourceMapping)
	spec.Expect(create).ToEqual(`DROP TABLE IF EXISTS dest;
DELETE FROM source_versions USING sources WHERE sources.id = source_versions.source_id AND sources.name = 'source';
DELETE FROM sources WHERE sources.name = 'source';
`)
}

func TestMappingToIndex(t *testing.T) {
	spec := tests.Spec(t)
	create := MappingToIndex(sourceMapping)
	spec.Expect(create).ToEqual(`CREATE INDEX ON dest (id);
`)
}