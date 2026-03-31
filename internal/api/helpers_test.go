package api

import (
	"encoding/json"
	"testing"
)

func TestNormalizeThingID(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		prefix string
		want   string
	}{
		{name: "raw post id", input: "abc123", prefix: "t3", want: "t3_abc123"},
		{name: "prefixed post id", input: "t3_abc123", prefix: "t3", want: "t3_abc123"},
		{name: "different prefix retained", input: "t1_xyz", prefix: "t3", want: "t1_xyz"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeThingID(tc.input, tc.prefix); got != tc.want {
				t.Fatalf("normalizeThingID(%q, %q) = %q, want %q", tc.input, tc.prefix, got, tc.want)
			}
		})
	}
}

func TestFlattenComments(t *testing.T) {
	const payload = `{
		"data": {
			"children": [
				{
					"kind": "t1",
					"data": {
						"id": "c1",
						"name": "t1_c1",
						"author": "alice",
						"body": "top level",
						"score": 10,
						"replies": {
							"data": {
								"children": [
									{
										"kind": "t1",
										"data": {
											"id": "c2",
											"name": "t1_c2",
											"author": "bob",
											"body": "reply",
											"score": 5,
											"replies": ""
										}
									}
								]
							}
						}
					}
				}
			]
		}
	}`

	var listing listingResponse[Comment]
	if err := json.Unmarshal([]byte(payload), &listing); err != nil {
		t.Fatalf("unmarshal listing: %v", err)
	}

	comments := flattenComments(listing.Data.Children, 0)
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
	if comments[0].Depth != 0 || comments[1].Depth != 1 {
		t.Fatalf("unexpected depths: %#v", comments)
	}
}
