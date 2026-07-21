package importer

import "testing"

func TestApplyRulesRecognizesLeaveAndStatusAliases(t *testing.T) {
	rules := map[string]Rule{
		ruleKey("leave", "Cuti Besar Dipotong"): {
			SourceField: "leave", Code: "Cuti Besar Dipotong",
			Label: "Cuti Besar Dipotong", Rate: 0.025,
		},
		ruleKey("status", "I"): {
			SourceField: "status", Code: "I",
			Label: "Izin Tidak Masuk", Rate: 0.05,
		},
	}
	record := RawRecord{LeaveType: "  cuti   besar dipotong ", AttendanceStatus: "Izin Tidak Masuk"}

	applyRules(&record, rules)

	if record.DeductionRate != 0.075 {
		t.Fatalf("deduction rate = %v, want 0.075", record.DeductionRate)
	}
	if len(record.DeductionComponents) != 2 {
		t.Fatalf("components = %d, want 2", len(record.DeductionComponents))
	}
}

func TestApplyRulesKeepsZeroRateLeaveAsSeparateComponent(t *testing.T) {
	rules := map[string]Rule{
		ruleKey("leave", "Cuti Tahunan"): {
			SourceField: "leave", Code: "Cuti Tahunan",
			Label: "Cuti Tahunan", Rate: 0,
		},
	}
	record := RawRecord{LeaveType: "Cuti Tahunan"}

	applyRules(&record, rules)

	if record.DeductionRate != 0 {
		t.Fatalf("deduction rate = %v, want 0", record.DeductionRate)
	}
	if len(record.DeductionComponents) != 1 || record.DeductionComponents[0].Label != "Cuti Tahunan" {
		t.Fatalf("unexpected components: %#v", record.DeductionComponents)
	}
}
