package quamina

import (
	"testing"
)

func TestPatternErrorHandling(t *testing.T) {
	_, err := patternFromJSON([]byte{})
	if err == nil {
		t.Error("accepted empty pattern")
	}
	_, err = patternFromJSON([]byte("33"))
	if err == nil {
		t.Error("accepted non-object JSON text")
	}
	_, err = patternFromJSON([]byte("{"))
	if err == nil {
		t.Error("accepted stub JSON object")
	}
	_, err = patternFromJSON([]byte("{ ="))
	if err == nil {
		t.Error("accepted malformed JSON object")
	}
	_, err = patternFromJSON([]byte(`{ "foo": `))
	if err == nil {
		t.Error("accepted stub JSON object")
	}
	_, err = patternFromJSON([]byte(`{ "foo": [`))
	if err == nil {
		t.Error("accepted stub JSON array")
	}

	_, err = patternFromJSON([]byte(`{ "foo": [ { "exists" == ] }`))
	if err == nil {
		t.Error("accepted stub JSON array")
	}

	_, err = patternFromJSON([]byte(`{ "foo": [ { "exists": false . ] }`))
	if err == nil {
		t.Error("accepted stub JSON array")
	}
}

func TestPatternFromJSON(t *testing.T) {
	bads := []string{
		`x`,
		`{"foo": ]`,
		`{"foo": 11 }`,
		`{"foo": "x" }`,
		`{"foo": true}`,
		`{"foo": null}`,
		`{"oof": [ ]`,
		`[33,22]`,
		`{"xxx": { }`,
		`{"xxx": [ [ 22 ] }`,
		`{"xxx": [ {"x": 1} ]`,
		`{"xxx": [ { [`,
		`{"xxx": [ { "exists": 23 } ] }`,
		`{"xxx": [ { "exists": true }, 15 ] }`,
		`{"xxx": [ { "exists": true, "a": 3 }] }`,
		`{"xxx": [ { "exists": false, "x": ["a", 3 ] }] }`,
		`{"abc": [ {"shellstyle":15} ] }`,
		`{"abc": [ {"shellstyle":"15"] ] }`,
		`{"abc": [ {"shellstyle":"15", "x", 1} ] }`,
		`{"abc": [ {"shellstyle":"a**b"}, "foo" ] }`,
		`{"abc": [ {"prefix":23}, "foo" ] }`,
		`{"abc": [ {"prefix":["a", "b"]}, "foo" ] }`,
		`{"abc": [ {"prefix": - }, "foo" ] }`,
		`{"abc": [ {"prefix":  - "a" }, "foo" ] }`,
		`{"abc": [ {"prefix":  "a" {, "foo" ] }`,
		`{"abc": [ {"equals-ignore-case":23}, "foo" ] }`,
		`{"abc": [ {"wildcard":"15", "x", 1} ] }`,
		`{"abc": [ {"wildcard":"a**b"}, "foo" ] }`,
		`{"abc": [ {"wildcard":"a\\b"}, "foo" ] }`,                                             // after JSON parsing, code sees `a/b`
		`{"abc": [ {"wildcard":"a\\"}, "foo" ] }`,                                              // after JSON parsing, code sees `a\`
		"{\"a\": [ { \"anything-but\": { \"equals-ignore-case\": [\"1\", \"2\" \"3\"] } } ] }", // missing ,
		"{\"a\": [ { \"anything-but\": { \"equals-ignore-case\": [1, 2, 3] } } ] }",            // no numbers
		"{\"a\": [ { \"anything-but\": { \"equals-ignore-case\": [\"1\", \"2\" } } ] }",        // missing ]
		"{\"a\": [ { \"anything-but\": { \"equals-ignore-case\": [\"1\", \"2\" ] } ] }",        // missing }
		"{\"a\": [ { \"equals-ignore-case\": 5 } ] }",
		"{\"a\": [ { \"equals-ignore-case\": [ \"abc\" ] } ] }",
	}
	for _, b := range bads {
		_, err := patternFromJSON([]byte(b))
		if err == nil {
			t.Error("accepted bad pattern: " + b)
		}
	}

	goods := []string{
		`{"x": [ 2 ]}`,
		`{"x": [ null, true, false, "hopp", 3.072e-11] }`,
		`{"x": { "a": [27, 28], "b": { "m": [ "a", "b" ] } } }`,
		`{"x": [ {"exists": true} ] }`,
		`{"x": { "y": [ {"exists": false} ] } }`,
		`{"abc": [ 3, {"shellstyle":"a*b"} ] }`,
		`{"abc": [ {"shellstyle":"a*b"}, "foo" ] }`,
		`{"abc": [ {"shellstyle":"a*b*c"} ] }`,
		`{"x": [ {"equals-ignore-case":"a*b*c"} ] }`,
		`{"abc": [ 3, {"wildcard":"a*b"} ] }`,
		`{"abc": [ {"wildcard":"a*b"}, "foo" ] }`,
		`{"abc": [ {"wildcard":"a*b*c"} ] }`,
		`{"abc": [ {"wildcard":"a*b\\*c"} ] }`,
	}
	w1 := []*patternField{{path: "x", vals: []typedVal{{vType: numberType, val: "2"}}}}
	w2 := []*patternField{{path: "x", vals: []typedVal{
		{literalType, "null", nil, nil},
		{literalType, "true", nil, nil},
		{literalType, "false", nil, nil},
		{stringType, `"hopp"`, nil, nil},
		{numberType, "3.072e-11", nil, nil},
	}}}
	w3 := []*patternField{
		{path: "x\na", vals: []typedVal{
			{numberType, "27", nil, nil},
			{numberType, "28", nil, nil},
		}},
		{path: "x\nb\nm", vals: []typedVal{
			{stringType, `"a"`, nil, nil},
			{stringType, `"b"`, nil, nil},
		}},
	}
	w4 := []*patternField{
		{
			path: "x", vals: []typedVal{
				{vType: existsTrueType, val: ""},
			},
		},
	}
	w5 := []*patternField{
		{
			path: "x\ny", vals: []typedVal{
				{vType: existsFalseType, val: ""},
			},
		},
	}
	w6 := []*patternField{
		{
			path: "abc", vals: []typedVal{
				{vType: stringType, val: "3"},
				{vType: shellStyleType, val: `"a*b"`},
			},
		},
	}
	w7 := []*patternField{
		{
			path: "abc", vals: []typedVal{
				{vType: shellStyleType, val: `"a*b"`},
				{vType: stringType, val: `"foo"`},
			},
		},
	}
	w8 := []*patternField{
		{
			path: "abc", vals: []typedVal{
				{vType: shellStyleType, val: `"a*b*c"`},
			},
		},
	}
	w9 := []*patternField{
		{
			path: "x", vals: []typedVal{
				{vType: monocaseType, val: `"a*b*c"`},
			},
		},
	}
	w10 := []*patternField{
		{
			path: "abc", vals: []typedVal{
				{vType: stringType, val: "3"},
				{vType: wildcardType, val: `"a*b"`},
			},
		},
	}
	w11 := []*patternField{
		{
			path: "abc", vals: []typedVal{
				{vType: wildcardType, val: `"a*b"`},
				{vType: stringType, val: `"foo"`},
			},
		},
	}
	w12 := []*patternField{
		{
			path: "abc", vals: []typedVal{
				{vType: wildcardType, val: `"a*b*c"`},
			},
		},
	}
	w13 := []*patternField{
		{
			path: "abc", vals: []typedVal{
				{vType: wildcardType, val: `"a*b\*c"`},
			},
		},
	}
	wanted := [][]*patternField{w1, w2, w3, w4, w5, w6, w7, w8, w9, w10, w11, w12, w13}

	for i, good := range goods {
		fields, err := patternFromJSON([]byte(good))
		if err != nil {
			t.Error("pattern:" + good + ": " + err.Error())
		}
		w := wanted[i]
		if len(w) != len(fields) {
			t.Errorf("at %d len(w)=%d, len(fields)=%d", i, len(w), len(fields))
		}
		for j, ww := range w {
			if ww.path != fields[j].path {
				t.Error("pathSegments mismatch: " + ww.path + "/" + fields[j].path)
			}
			for k, www := range ww.vals {
				if www.val != fields[j].vals[k].val {
					t.Errorf("At [%d][%d], val mismatch %s/%s", j, k, www.val, fields[j].vals[k].val)
				}
			}
		}
	}
}
