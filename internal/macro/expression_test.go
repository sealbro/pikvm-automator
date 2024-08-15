package macro

import "testing"

func TestExpression_SimpleParse(t *testing.T) {
	expressions := []string{
		"1ms",
		"2s",
		"@left|@right",
		"@left+@0'0+@920'560+@200'200|@right",
		"@0'0|@left|@920'560|@left|@200'200|@right",
		"MetaLeft",
		"MetaLeft|@12'21",
		"@0'0|@12'21|@222'222",
		"MetaLeft+KeyD",
		"MetaLeft+KeyD+1ms|KeyC+KeyV+10ms",
		"MetaLeft+2ms+KeyD+1ms",
		"MetaLeft+KeyS+1ms|100s",
		"[5](MetaLeft+KeyD|10s)",
		"[10](MetaLeft+KeyD+1ms)",
		"[10](MetaLeft+KeyD+1ms|10s)",
		"[10](MetaLeft+KeyD+1ms|10s)|MetaLeft|5s",
		"[5]([20](MetaLeft+KeyD+1ms|12s)|KeyA|15s)",
		"[5]([20](MetaLeft+KeyD+1ms|12s)|KeyA|15s)|MetaLeft|5s",
	}

	for _, exp := range expressions {
		t.Run(exp, func(t *testing.T) {
			macros := New(exp)
			macros.Parse()
			s := macros.String()
			if s != exp {
				t.Errorf("expected %s, got %s", exp, s)
			}
		})
	}
}

func TestExpression_ParseWithExceptions(t *testing.T) {
	expressions := map[string]string{
		"[10](MetaLeft+KeyD+1ms|10s)MetaLeft|5s":               "[10](MetaLeft+KeyD+1ms|10s)|MetaLeft|5s",
		"[5]([20](MetaLeft+KeyD+1ms|12s)|KeyA|15s)MetaLeft|5s": "[5]([20](MetaLeft+KeyD+1ms|12s)|KeyA|15s)|MetaLeft|5s",
	}

	for _, exp := range expressions {
		t.Run(exp, func(t *testing.T) {
			macros := New(exp)
			macros.Parse()
			s := macros.String()
			if s != exp {
				t.Errorf("expected %s, got %s", exp, s)
			}
		})
	}
}
